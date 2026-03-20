package main

import (
	"os"

	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/v2"
	"golang.org/x/mod/semver"
)

// - Uses requirements.go to validate the requirements
// - Get Root Env
// - Define Paths
// - Add semantic Versioning

var paths = make(map[string]string)

// Function that initializes the Path Map
func utils() {

	check_requirement()

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

}

// Initializes utils function and returns the settings
func init_utils() string {
	utils()
	return paths["settings"]
}

// Check if the version corresponds to the semantic versioning.
func is_semantic_version(version string) {
	if !semver.IsValid(version) {
		logger.Error("The version of terrabutler is not valid.")
		os.Exit(1)
	}
}
