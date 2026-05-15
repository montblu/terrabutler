package utils

import (
	"errors"

	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/afero"
	"golang.org/x/mod/semver"

	"github.com/l58193/terrabutler/internal/logger"
)

// - Gets the Current Environment
// - Get Root Env
// - Define Paths
// - Add semantic Versioning

var Paths, _ = init_paths()

var CurrentEnv, _ = getCurrentEnv(afero.NewOsFs())

// Function that initializes the Path Map
func init_paths() (map[string]string, error) {

	var paths = make(map[string]string)

	// Loading koanf instance
	var k = koanf.New(".")

	//Getting the environment variables
	err := k.Load(env.Provider(".", env.Opt{Prefix: "TERRABUTLER_"}), nil)
	if err != nil {
		return paths, errors.New("An error occured while loading the environment variables: " + error.Error(err))
	}
	root := k.String("TERRABUTLER_ROOT")

	paths["backends"] = root + "/configs/backends"
	paths["environment"] = root + "/site_inception/.terraform/environment"
	paths["inception"] = root + "/site_inception"
	paths["root"] = root
	paths["settings"] = root + "/configs/settings.yml"
	paths["templates"] = root + "/configs/templates"
	paths["variables"] = root + "/configs/variables"

	return paths, nil
}

// Returns the settings path
func Settings_path() string {
	return Paths["settings"]
}

// Check if the version corresponds to the semantic versioning.
func Is_semantic_version(version string) error {
	if !semver.IsValid(version) {
		return errors.New("The version of terrabutler is not valid.")
	}
	return nil
}

// Get current_environment
func getCurrentEnv(fs afero.Fs) (string, error) {

	//Open site environment file
	env, err := afero.ReadFile(fs, Paths["environment"])
	logger.Zap.Debug("Path of Environment is: " + Paths["environment"])
	if err != nil {
		return "", errors.New("An error has occurred while reading the current environment: " + err.Error())
	}

	return string(env), nil
}
