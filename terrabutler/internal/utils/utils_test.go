package utils

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestInitPaths(t *testing.T) {

	rootDir := "ROOT"

	os.Setenv("TERRABUTLER_ROOT", rootDir)

	Paths, err := init_paths()

	assert.NoError(t, err, "An error occured.")

	assert.Equal(t, Paths["backends"], rootDir+"/configs/backends", "The path for backends is incorrect.")
	assert.Equal(t, Paths["environment"], rootDir+"/site_inception/.terraform/environment", "The path for environment is incorrect.")
	assert.Equal(t, Paths["inception"], rootDir+"/site_inception", "The path for inception is incorrect.")
	assert.Equal(t, Paths["root"], rootDir, "The path for root is incorrect.")
	assert.Equal(t, Paths["settings"], rootDir+"/configs/settings.yml", "The path for settings is incorrect.")
	assert.Equal(t, Paths["templates"], rootDir+"/configs/templates", "The path for templates is incorrect.")
	assert.Equal(t, Paths["variables"], rootDir+"/configs/variables", "The path for variables is incorrect.")

}

func TestValidSemanticVersion(t *testing.T) {
	assert.NoError(t, Is_semantic_version("v3.3.3"), "Failed, a valid version was given.")
}

func TestValidCurrentEnv(t *testing.T) {

	//Defining Paths
	Paths["environments"] = "ROOT/site_inception/.terraform/environment"
	envName := "env"

	// Use the in-memory filesystem
	fs := afero.NewMemMapFs()

	//Creating the environment file
	afero.WriteFile(fs, Paths["environment"], []byte(envName), 0644)

	env, err := getCurrentEnv(fs)

	assert.Equal(t, env, envName, "Failed, returned name of the current environment is incorrect.")
	assert.NoError(t, err, "Failed, an error occurred opening/reading the file.")
}

func TestInvalidSemanticVersion(t *testing.T) {
	assert.Error(t, Is_semantic_version("v1.Invalid.1"), "Failed, it was accepted a invalid semantic version.")
}

func TestInvalidCurrentEnv(t *testing.T) {

	envName := "env"

	// Use the in-memory filesystem
	fs := afero.NewMemMapFs()

	env, err := getCurrentEnv(fs)

	assert.NotEqual(t, env, envName, "Failed, returned name of the current environment is the same.")
	assert.Error(t, err, "Failed, there was no file to be read.")
}
