package utils

import (
	"errors"

	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/afero"
	"golang.org/x/mod/semver"

	"github.com/montblu/terrabutler/internal/logger"
)

// - Gets the Current Environment
// - Get Root Env
// - Define Paths
// - Add semantic Versioning

var Paths, _ = initPaths()
var testEnv string

// GetCurrentEnv reads and returns the active environment from the filesystem.
func GetCurrentEnv() string {
	env, err := currentEnv(afero.NewOsFs())
	if err != nil {
		return ""
	}
	return env
}

// SetCurrentEnvForTest sets a mock environment override for tests.
// Call with an empty string to reset.
func SetCurrentEnvForTest(env string) {
	testEnv = env
}

// Function that initializes the Path Map
func initPaths() (map[string]string, error) {

	var paths = make(map[string]string)

	// Loading koanf instance
	var k = koanf.New(".")

	// Getting the environment variables
	err := k.Load(env.Provider(".", env.Opt{Prefix: "TERRABUTLER_"}), nil)
	if err != nil {
		return paths, errors.New("An error occurred while loading the environment variables: " + err.Error())
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
func SettingsPath() string {
	return Paths["settings"]
}

// Check if the version corresponds to the semantic versioning.
func IsSemanticVersion(version string) error {
	if !semver.IsValid(version) {
		return errors.New("the version of terrabutler is not valid")
	}
	return nil
}

// currentEnv reads the environment file from the given filesystem.
func currentEnv(fs afero.Fs) (string, error) {
	if testEnv != "" {
		return testEnv, nil
	}

	// Open site environment file
	env, err := afero.ReadFile(fs, Paths["environment"])
	logger.Zap.Debug("Path of Environment is: " + Paths["environment"])
	if err != nil {
		return "", errors.New("An error has occurred while reading the current environment: " + err.Error())
	}

	return string(env), nil
}
