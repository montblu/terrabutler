package utils

import (
	"os"

	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
	"golang.org/x/mod/semver"

	"terrabutler/internal/logger"
	"terrabutler/internal/requirements"
)

// - Uses requirements.go to validate the requirements
// - Get Root Env
// - Define Paths
// - Add semantic Versioning

var Paths = utils()

var CurrentEnv = getCurrentEnv()

// Function that initializes the Path Map
func utils() map[string]string {

	var paths = make(map[string]string)

	requirements.Check_requirement()

	// Loading koanf instance
	var k = koanf.New(".")

	//Getting the environment variables
	k.Load(env.Provider(".", env.Opt{Prefix: "TERRABUTLER_"}), nil)
	root := k.String("TERRABUTLER_ROOT")

	paths["backends"] = root + "/configs/backends"
	paths["environment"] = root + "/site_inception/.terraform/environment"
	paths["inception"] = root + "/site_inception"
	paths["root"] = root
	paths["settings"] = root + "/configs/settings.yml"
	paths["templates"] = root + "/configs/templates"
	paths["variables"] = root + "/configs/variables"

	return paths
}

// Returns the settings path
func Settings_path() string {
	return Paths["settings"]
}

// Check if the version corresponds to the semantic versioning.
func Is_semantic_version(version string) {
	if !semver.IsValid(version) {
		logger.Zap.Error("The version of terrabutler is not valid.")
		os.Exit(1)
	}
}

// Get current_environment
func getCurrentEnv() string {

	//Open site environment file
	env, err := os.ReadFile(Paths["environment"])
	logger.Zap.Debug("Path of Environment is: " + Paths["environment"])
	if err != nil {
		logger.Zap.Error("An error has occurred while reading the current environment: ", zap.Error(err))
		os.Exit(1)
	}

	return string(env)
}
