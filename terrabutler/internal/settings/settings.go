package settings

import (
	"fmt"
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"

	"terrabutler/internal/logger"
	"terrabutler/internal/utils"
)

// - Uses utils.go for the paths
// - Makes the Schema for settings
//		- Be able to add schemas for json/yaml/toml/huml/hcl
// - Validates the settings
// - Writes the settings

var Path = utils.Settings_path()

// Global koanf instance, it has the Settings Configuration
var Conf = koanf.New(".")

// Loads the settings from the settings.yaml
func get_settings() {

	// Load default values using the confmap provider.
	Conf.Load(confmap.Provider(map[string]any{
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
	if err := Conf.Load(file.Provider(Path), yaml.Parser()); err != nil {
		logger.Zap.Error("Error occurred loading the settings: ", zap.Error(err))
		os.Exit(1)
	}
}

// Validates the settings files
func Validate_settings() {

	//Gets the settings
	get_settings()

	//Verify if all the required camps are filled
	for key, value := range Conf.All() {

		value = fmt.Sprint(value)
		if key != "hooks.post_env_select" && key != "hooks.pre_env_select" && (value == "<nil>" || value == "[]") {
			logger.Zap.Error(fmt.Sprintf("Invalid settings file. Invalid value for %s", key))
			os.Exit(1)
		}

	}

}

// Writes settings file
func Write_settings(newSettings *koanf.Koanf) {

	//Marshal the new Settings file as a yaml file
	b, _ := newSettings.Marshal(yaml.Parser())

	f, err := os.Create(Path)
	if err != nil {
		logger.Zap.Error(fmt.Sprint("An error has occurred opening the settings file: ", err))
		os.Exit(1)
	}
	l, err := f.Write(b)
	if l == 0 && err != nil {
		logger.Zap.Error(fmt.Sprint("An error has occurred writing to the settings file: ", err))
		f.Close()
		os.Exit(1)
	}
	err = f.Close()

	// Writes out to the location of the old settings
	logger.Zap.Info("Settings written successfully.")

}
