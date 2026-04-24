package main

import (
	"fmt"
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

// - Uses utils.go for the paths
// - Makes the Schema for settings
//		- Be able to add schemas for json/yaml/toml/huml/hcl
// - Validates the settings
// - Writes the settings

var settingsPath = init_utils()

// Global koanf instance, it has the Settings Configuration
var settings = koanf.New(".")

// Loads the settings from the settings.yaml
func get_settings() {

	// Load default values using the confmap provider.
	settings.Load(confmap.Provider(map[string]any{
		"general.organization":                                nil,
		"general.secrets_key_id":                              nil,
		"sites.ordered":                                       []string{},
		"environments.default.domain":                         nil,
		"environments.default.name":                           nil,
		"environments.default.profile_name":                   nil,
		"environments.default.region":                         nil,
		"environments.permanent":                              []string{},
		"environments.temporary.secrets.firebase_credentials": nil,
		"environments.temporary.secrets.mail_password":        nil,
		//Optional
		"hooks.pre_env_select":  nil,
		"hooks.post_env_select": nil,
	}, "."), nil)

	// Load YAML config on top of the default values.
	if err := settings.Load(file.Provider(settingsPath), yaml.Parser()); err != nil {
		logger.Error("Error occurred loading the settings: ", zap.Error(err))
		os.Exit(1)
	}

	//logger.Debug("Settings File Loaded", zap.String("Settings", fmt.Sprint(settings.All())))
}

// Validates the settings files
func validate_settings() {

	//Gets the settings
	get_settings()

	//Verify if all the required camps are filled
	for key, value := range settings.All() {

		value = fmt.Sprint(value)
		if key != "hooks.post_env_select" && key != "hooks.pre_env_select" && (value == "<nil>" || value == "[]") {
			logger.Error(fmt.Sprintf("Invalid settings file. Invalid value for %s", key))
			os.Exit(1)
		}

	}

}

// Writes settings file
func write_settings(newSettings *koanf.Koanf) {

	//Marshal the new Settings file, to the specific type
	b, _ := newSettings.Marshal(yaml.Parser())

	f, err := os.Create(settingsPath)
	if err != nil {
		logger.Error(fmt.Sprint("An error has occurred opening the file: ", err))
		os.Exit(1)
	}
	l, err := f.Write(b)
	if l == 0 && err != nil {
		logger.Error(fmt.Sprint("An error has occurred writing to the file: ", err))
		f.Close()
		os.Exit(1)
	}
	err = f.Close()

	// Writes out to the location of the old settings
	logger.Info("Settings written successfully.")

}
