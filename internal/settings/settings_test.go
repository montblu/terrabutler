package settings

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestWriteSettings(t *testing.T) {

	//Path where the settings are
	Path = "PATH/settings.yml"
	oldPath := "PATH/settingsOld.yml"

	settingsFile := []byte(`
environments:
    default:
        domain: domain
        name: staging
        profile_name: domain-staging
        region: DUMMY
    permanent:
        - staging
        - production
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
        - inception`)

	// Use the in-memory filesystem
	fs := afero.NewMemMapFs()

	//Creating both settings files
	afero.WriteFile(fs, Path, settingsFile, 0644)
	afero.WriteFile(fs, oldPath, settingsFile, 0644)

	Validate_settings(fs)

	//Modify current settings file, with the same use case in the program
	envs := Conf.Strings("environments.permanent")
	envs = append(envs, "NewEnvironment")
	err := Conf.Set("environments.permanent", envs)

	assert.NoError(t, err, "Failed, an error accessing the test files has occurred.")

	assert.NoError(t, Write_settings(fs, Conf), "Failed, an error has occurred.")

	oldSettings, err := afero.ReadFile(fs, oldPath)
	assert.NoError(t, err, "Failed, An error occurred while getting the data from the old settings file.")

	newSettings, err := afero.ReadFile(fs, Path)
	assert.NoError(t, err, "Failed, An error occurred while getting the data from the new settings file.")

	assert.NotEqual(t, string(oldSettings), string(newSettings), "Failed, the settings file was not updated.")

}

func TestGetSettings(t *testing.T) {

	//Path where the settings are
	Path = "PATH/settings.yml"

	// Use the in-memory filesystem
	fs := afero.NewMemMapFs()

	//Creating the settings file
	//It accepts empty files but they will fail in the validateSettings function
	afero.WriteFile(fs, Path, []byte(``), 0644)

	assert.NoError(t, get_settings(fs), "Failed, the file exists with an valid output.")

}

func TestValidValidateSettings(t *testing.T) {

	//Path where the settings are
	Path = "PATH/settings.yml"

	// Use the in-memory filesystem
	fs := afero.NewMemMapFs()

	//Creating the settings file
	afero.WriteFile(fs, Path, []byte(`
environments:
    default:
        domain: domain
        name: staging
        profile_name: domain-staging
        region: DUMMY
    permanent:
        - staging
        - production
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

	assert.NoError(t, Validate_settings(fs), "Failed, it was used a valid settings file.")

}

func TestInvalidGetSettings(t *testing.T) {

	//Path where the settings are
	Path = "PATH/settings.yml"

	// Use the in-memory filesystem
	fs := afero.NewMemMapFs()

	assert.Error(t, get_settings(fs), "Failed, the file doesn't exist.")

	//If the file exists but have a invalid input
	afero.WriteFile(fs, Path, []byte(`Invalid Input`), 0644)
	assert.Error(t, get_settings(fs), "Failed, the file exists but the input is invalid.")

}

func TestInvalidValidateSettings(t *testing.T) {

	//Path where the settings are
	Path = "PATH/settings.yml"

	// Use the in-memory filesystem
	fs := afero.NewMemMapFs()

	//Creating the settings file empty
	afero.WriteFile(fs, Path, []byte(""), 0644)

	assert.Error(t, Validate_settings(fs), "Failed, the settings file was empty.")

	assert.Error(t, Validate_settings(fs), "Failed, the settings was the correct structure but its empty.")

	//Writing the file with the correct structure but with invalid types
	afero.WriteFile(fs, Path, []byte(`
environments:
    default:
        domain: null
        name: null
        profile_name: null
        region: null
    permanent: []
    temporary:
        secrets:
            firebase_credentials: null
            mail_password: null
general:
    organization:	null
    secrets_key_id: null
hooks:
    post_env_select: null
    pre_env_select: null
sites:
    ordered: []`), 0644)

	assert.Error(t, Validate_settings(fs), "Failed, the settings was the correct structure but has nil values.")
}
