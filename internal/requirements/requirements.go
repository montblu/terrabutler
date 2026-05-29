package requirements

import (
	"errors"
	"os"

	"github.com/l58193/terrabutler/internal/logger"

	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/afero"
)

// Checks the requirements before running the application.
func Check_requirement(fs afero.Fs) error {

	// Sync logger
	defer logger.Zap.Sync()

	// Loading koanf instance
	var k = koanf.New(".")

	//Getting the environment variables
	err := k.Load(env.Provider(".", env.Opt{Prefix: "TERRABUTLER_"}), nil)
	if err != nil {
		return errors.New("An error occured while loading the environment variables: " + error.Error(err))
	}
	root := k.String("TERRABUTLER_ROOT")
	isEnabled := k.Bool("TERRABUTLER_ENABLE")
	settingsFile := root + "/configs/settings.yml"

	if !isEnabled {
		return errors.New("Terrabutler is not currently enabled on this folder. Please set 'TERRABUTLER_ENABLE' in your environment to true to enable it.")
	}
	if root == "" {
		return errors.New("Terrabutler can't determine the root folder of your project or it doesn't exist. Please set 'TERRABUTLER_ROOT' in your environment pointing to the root folder of your project.")
	}
	if _, err := fs.Stat(settingsFile); os.IsNotExist(err) {
		return errors.New("Terrabutler can't find you settings file. Please create a 'settings.yml' file inside the 'configs' folder.")
	}

	return nil
}
