package settings

import (
	"errors"
	"fmt"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/afero"

	"github.com/l58193/terrabutler/internal/logger"
	"github.com/l58193/terrabutler/internal/utils"
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
func get_settings(fs afero.Fs) error {

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

	b, err := afero.ReadFile(fs, Path)
	if err != nil {
		return errors.New("Error occurred while trying to read the settings: " + err.Error())
	}

	// Load YAML config on top of the default values.
	if err := Conf.Load(rawbytes.Provider(b), yaml.Parser()); err != nil {
		return errors.New("Error occurred loading the settings: " + err.Error())
	}
	return nil
}

// Validates the settings files
func Validate_settings(fs afero.Fs) error {

	//Gets the settings
	err := get_settings(fs)
	if err != nil {
		return err
	}

	//Verify if all the required camps are filled
	for key, value := range Conf.All() {

		value = fmt.Sprint(value)
		if key != "hooks.post_env_select" && key != "hooks.pre_env_select" && (value == "<nil>" || value == "[]") {
			return errors.New("Invalid settings file. Invalid value for " + key)
		}

	}

	return nil

}

// Writes settings file
func Write_settings(fs afero.Fs, newSettings *koanf.Koanf) error {

	//Marshal the new Settings file as a yaml file
	b, _ := newSettings.Marshal(yaml.Parser())

	f, err := fs.Create(Path)
	if err != nil {
		f.Close()
		return errors.New("An error has occurred opening the settings file: " + err.Error())
	}
	l, err := f.Write(b)
	if l == 0 && err != nil {
		f.Close()
		return errors.New("An error has occurred writing to the settings file: " + err.Error())
	}
	err = f.Close()

	// Writes out to the location of the old settings
	logger.Zap.Info("Settings written successfully.")
	return nil

}
