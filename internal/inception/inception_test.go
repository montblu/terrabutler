package inception

import (
	"terrabutler/internal/settings"
	"terrabutler/internal/utils"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestInitNeeded(t *testing.T) {

	// Define Path "inception"
	utils.Paths["inception"] = "inception"

	// Use the in-memory filesystem
	fs := afero.NewMemMapFs()

	//Creating environment file initialized with env
	afero.WriteFile(fs, utils.Paths["inception"]+"/.terraform/environment", []byte("env"), 0644)

	assert.NoError(t, Init_needed(fs), "Failed, the environment file exists.")

}

func TestInit(t *testing.T) {

	// Mockable tf function (the output isn't important here...)
	commandRunnerNoVisibleOutput = func(command, site string, args, options []string, needed_options string) ([]byte, error) {
		return nil, nil
	}

	settings.Conf.Set("environments.default.name", "env")
	utils.Paths["inception"] = "inception"

	// Use the in-memory filesystem
	fs := afero.NewMemMapFs()
	fs.MkdirAll(utils.Paths["inception"]+"/.terraform", 0644)
	fs.MkdirAll(utils.Paths["backends"], 0644)

	//Create the environment file
	assert.NoError(t, Init(fs), "Failed, the environment file couldn't be created")

	// Validate the environment file created
	env, err := afero.ReadFile(fs, utils.Paths["inception"]+"/.terraform/environment")
	assert.NoError(t, err, "Failed, the environment file couldn't be redden.")

	// See if the content is correct
	assert.Equal(t, settings.Conf.String("environments.default.name"), string(env), "Failed, the default environment name isn't correct in the environment file.")

}
