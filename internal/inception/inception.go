package inception

import (
	"errors"
	"os"

	"github.com/montblu/terrabutler/internal/logger"
	"github.com/montblu/terrabutler/internal/settings"
	"github.com/montblu/terrabutler/internal/tf"
	"github.com/montblu/terrabutler/internal/utils"

	"github.com/spf13/afero"
)

// Creating var for mockable function
var commandRunnerNoVisibleOutput = tf.CommandRunnerNoVisibleOutput

func inceptionInitCheck(fs afero.Fs) bool {
	dir := utils.Paths["inception"]
	if _, err := fs.Stat(dir + "/.terraform"); !os.IsNotExist(err) {
		if _, err2 := fs.Stat(dir + "/.terraform/environment"); !os.IsNotExist(err2) {
			return true
		}
	}
	return false
}

func InitNeeded(fs afero.Fs) error {

	if !inceptionInitCheck(fs) {
		return errors.New("the initialization hasn't been made yet. Please execute the following command to initialize it: terrabutler init")
	}
	return nil
}

func Init(fs afero.Fs) error {

	default_env_name := settings.Conf.String("environments.default.name")
	inception_dir := utils.Paths["inception"]

	if !inceptionInitCheck(fs) {

		_, err := commandRunnerNoVisibleOutput("init", "inception", []string{}, []string{}, "backend")
		if err != nil {
			return errors.New("There was an error while doing the initialization, Error info: " + err.Error())
		}
		// Try opening the new inception dir and create a file and write a file
		f, err := fs.OpenFile(inception_dir+"/.terraform/environment", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return errors.New("The file that manages the environments could not be created, Error info: " + err.Error())
		}
		_, err = f.Write([]byte(default_env_name))
		if err != nil {
			return errors.New("Writing on the file that manages the environments wasn't possible, Error info: " + err.Error())
		}
		if err := f.Close(); err != nil {
			return errors.New("failed to close the environment file: " + err.Error())
		}

		// If all is ok, display successfully inception
		logger.Zap.Info("The initialization was successful!")
	} else {
		logger.Zap.Warn("The initialization was already done.")
	}

	return nil

}
