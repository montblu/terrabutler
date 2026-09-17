package inception

import (
	"testing"

	"github.com/montblu/terrabutler/internal/settings"
	"github.com/montblu/terrabutler/internal/utils"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestInitNeeded(t *testing.T) {

	// Define Path "inception"
	utils.Paths["inception"] = "inception"

	// Use the in-memory filesystem
	fs := afero.NewMemMapFs()

	// Creating environment file initialized with env
	_ = afero.WriteFile(fs, utils.Paths["inception"]+"/.terraform/environment", []byte("env"), 0644)

	assert.NoError(t, InitNeeded(fs), "Failed, the environment file exists.")

}

func TestInit(t *testing.T) {

	// Mockable tf function (the output isn't important here...)
	commandRunnerNoVisibleOutput = func(command, site string, args, options []string, needed_options string, fs afero.Fs) ([]byte, error) {
		return nil, nil
	}

	envDefault := "env"

	settings.Conf.Set("environments.default.name", envDefault) //nolint:errcheck
	utils.Paths["inception"] = "inception"
	// Creating the sites which will be tested the file creation
	settings.Conf.Set("sites.ordered", []string{"site-a", "site-b", "site-c"})

	// Use the in-memory filesystem
	fs := afero.NewMemMapFs()
	_ = fs.MkdirAll(utils.Paths["inception"]+"/.terraform", 0644)
	_ = fs.MkdirAll(utils.Paths["backends"], 0644)

	// Create the environment file
	assert.NoError(t, Init(fs), "Failed, the environment file couldn't be created")

	// Validate the environment file created
	env, err := afero.ReadFile(fs, utils.Paths["inception"]+"/.terraform/environment")
	assert.NoError(t, err, "Failed, the environment file couldn't be read.")

	assert.Equal(t, settings.Conf.String("environments.default.name"), string(env), "Failed, the default environment name isn't correct in the environment file.")

	// Test for each site if the environment file was created with the default environment name
	for _, site := range settings.Conf.Strings("sites.ordered") {
		siteEnv, err := afero.ReadFile(fs, utils.Paths["root"]+"/site_"+site+"/.terraform/.terrabutler_env")
		assert.NoError(t, err, "Failed, the environment file for site "+site+" couldn't be read.")
		assert.Equal(t, envDefault, string(siteEnv), "Failed, the default environment name isn't correct in the environment file for site "+site+".")
	}
}
