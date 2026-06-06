package env

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/montblu/terrabutler/internal/settings"
	"github.com/montblu/terrabutler/internal/utils"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestConfirmationMenu(t *testing.T) {

	// String with the question
	questionString := "Question"

	// Buffer to sent the user input
	yesResponse := bytes.NewBuffer([]byte("yes"))
	scanner := bufio.NewScanner(yesResponse)

	// Calling the function
	choiceYes, err := confirmationMenu(questionString, *scanner)

	// Valitaing the result
	assert.NoError(t, err, "Failed, an error occurred while giving a Yes input.")
	assert.Equal(t, true, choiceYes, "Failed, the input was yes, but returned false")

	// Same for the invalid response
	noResponse := bytes.NewBuffer([]byte("\n invalid input"))
	scanner = bufio.NewScanner(noResponse)
	choiceNo, err := confirmationMenu(questionString, *scanner)
	assert.Equal(t, false, choiceNo, "Failed, the input was no/invalid, but returned true")
	assert.NoError(t, err, "Failed, an error occurred while giving a invalid/no input")

}

func TestProtectedEnv(t *testing.T) {

	// Defining the settings values
	_ = settings.Conf.Set("environments.permanent", []string{"OtherEnv1", "protectedEnv", "OtherEnv2"})

	assert.Equal(t, true, isProtectedEnv("protectedEnv"), "Failed, the environment is in the list of the protected environments.")
	assert.Equal(t, false, isProtectedEnv("randomEnv"), "Failed, the environment is not in the list of the protected environments.")

}

// Mocking execution of GetAvailableEnvs, because it need the output from terraform
func mockGetAvailableEnvs(fs afero.Fs) ([]string, error) {
	return []string{"env0", "env1", "env2", "env3"}, nil
}

// Mocks for all the tf functions used..
func mockCommandRunner(command string, site string, args []string, flags []string, oa string) error {
	return nil
}

func mockRunner(command []string, site string, envVars []string) ([]byte, error) {
	return nil, nil
}

func mockFunc() error {
	return nil
}

// Verify the change to the file
func TestSetCurrentEnv(t *testing.T) {

	// Setting up the mock functions used
	GetAvailableEnvs = mockGetAvailableEnvs
	runnerNoVisibleOutput = mockRunner
	initAllSites = mockFunc

	// Defining the environment file path
	utils.Paths["environment"] = "PATH/environment"

	// Use the in-memory filesystem
	fs := afero.NewMemMapFs()

	// Choose a existing env from the mockGetAvailableEnvs
	env := "env1"

	// Run the tests for a valid environment
	assert.NoError(t, SetCurrentEnv(env, false, fs), "Failed, the environment is in the list of available environments.")
	newEnv, err := afero.ReadFile(fs, utils.Paths["environment"])
	assert.NoError(t, err, "Failed, the environment file couldn't be readen.")
	assert.Equal(t, env, string(newEnv), "Failed, the environment file wasn't been updated.")

	// Choosing a invalid environment
	env = "invalid"

	// Run the tests for a invalid environment
	assert.Error(t, SetCurrentEnv(env, false, fs), "Failed, the environment is not in the list of available environments.")
	newEnv, err = afero.ReadFile(fs, utils.Paths["environment"])
	assert.NoError(t, err, "Failed, the environment file couldn't be readen.")
	assert.NotEqual(t, env, string(newEnv), "Failed, the environment file was been updated to an invalid environment.")

}

// Verify files where deleted
func TestDeleteEnv(t *testing.T) {

	// Mocking GetAvailableEnvironments
	GetAvailableEnvs = mockGetAvailableEnvs
	// Mocking tf commands
	commandRunner = mockCommandRunner

	// Defining the variables used
	org := "org"
	_ = settings.Conf.Set("general.organization", org)
	// The protected envs should exist in the mock GetAvailableEnvironments
	protectedEnvs := []string{"env0", "env3"}
	_ = settings.Conf.Set("environments.permanent", protectedEnvs)
	// Defining the current environment
	current_env = "env0"
	// The environment to be deleted
	env := "env1"

	// Defining the environment file path
	utils.Paths["variables"] = "PATH"

	// Use the in-memory filesystem
	fs := afero.NewMemMapFs()

	_ = afero.WriteFile(fs, utils.Paths["variables"]+"/"+org+"-"+env+"-site.tfvars", nil, 0644)
	_ = afero.WriteFile(fs, utils.Paths["variables"]+"/"+org+"-"+env+".tfvars", nil, 0644)

	assert.NoError(t, DeleteEnv(env, true, false, fs), "Failed, the environment was a valid one.")

	// Verify if all the file were deleted
	files, err := afero.ReadDir(fs, utils.Paths["variables"])
	assert.NoError(t, err, "Failed, an error has occurred while reading the variables folder.")
	for _, file := range files {
		if !file.IsDir() {
			if strings.Contains(file.Name(), org+"-"+env) {
				t.Error("The file " + file.Name() + " hasn't been deleted.")
			}
		}
	}

	// Trying to delete a invalid environment
	env = "invalid"
	assert.NoError(t, DeleteEnv(env, true, false, fs), "Failed, the environment doesn't exist, should give a warning.")
	// Trying to delete the current environment
	assert.Error(t, DeleteEnv(current_env, true, false, fs), "Failed, the environment to be deleted is the current one.")
	// Trying to delete a protected environment
	env = "env3"
	assert.Error(t, DeleteEnv(env, true, false, fs), "Failed, the environment to be deleted is a permanent one.")
}

// Verify if the current environment was changed and the settings file updated.
func TestCreateEnv(t *testing.T) {
	// Create the settings file
	settings.Path = "PATH/settings.yml"

	// Use the in-memory filesystem
	fs := afero.NewMemMapFs()

	// Creating the settings file
	_ = afero.WriteFile(fs, settings.Path, []byte(`
environments:
    default:
        domain: domain
        name: staging
        profile_name: domain-staging
        region: DUMMY
    permanent:
        - env0
        - env3
    temporary:
        secrets:
            firebase_credentials: DUMMY
            mail_password: DUMMY
general:
    organization:	org
    secrets_key_id: alias/org-staging-secrets
hooks:
    post_env_select: null
    pre_env_select: null
sites:
    ordered:
        - inception`), 0644)

	// The new environment to be created
	newEnv := "newEnv"

	// Mocking GetAvailableEnvironments
	GetAvailableEnvs = mockGetAvailableEnvs
	// Mocking tf commands
	commandRunner = mockCommandRunner

	// Update the settings
	assert.NoError(t, settings.Validate_settings(fs), "Failed, the created settings file should be a valid one.")

	assert.NoError(t, CreateEnv(newEnv, true, false, true, fs), "Failed, the environment to be created is a valid one.")
	// Update the settings
	assert.NoError(t, settings.Validate_settings(fs), "Failed, the updated settings file should be a valid one.")
	// Get the new permanent environments
	newPermanentEnvs := settings.Conf.Strings("environments.permanent")
	assert.Contains(t, newPermanentEnvs, newEnv, "Failed, the settings file wasn't been updated correctly.")

	// Trying to create a invalid environment
	// Is a environment that already exists
	newEnv = "env2"
	assert.NoError(t, CreateEnv(newEnv, true, false, true, fs), "Failed, the environment already exists, should give a warning.")

}
