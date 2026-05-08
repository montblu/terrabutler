package inception

import (
	"fmt"
	"os"
	"os/exec"
	"terrabutler/internal/logger"
	"terrabutler/internal/settings"
	"terrabutler/internal/utils"

	"go.uber.org/zap"
)

func inception_init_check() bool {
	dir := utils.Paths["inception"]
	if _, err := os.Stat(dir + "/.terraform"); !os.IsNotExist(err) {
		if _, err2 := os.Stat(dir + "/.terraform/environment"); !os.IsNotExist(err2) {
			return true
		}
	}
	return false
}

func Init_needed() {

	if !inception_init_check() {
		logger.Zap.Error("The initialization hasn't been made yet. Please execute the following command to initialize it: terrabutler init")
		os.Exit(1)
	}
}

func Init() {
	org := settings.Conf.String("general.organization")
	default_env_name := settings.Conf.String("environments.default.name")
	inception_dir := utils.Paths["inception"]
	backend_dir := utils.Paths["backends"]

	if !inception_init_check() {

		//run Terraform init -backend-config LOCATION
		os.Chdir(inception_dir)

		//Verifies that terraform exist in the current dir
		_, err := exec.LookPath("terraform")
		if err != nil {
			logger.Zap.Error("No Terraform executable found.")
			os.Exit(1)
		}

		command := []string{"terraform", "init", "-backend-config", backend_dir + "/" + org + "-" + default_env_name + "-inception.tfvars"}

		logger.Zap.Debug(fmt.Sprintf("Executing %s command with args: %v", command[0], command[1:]))

		//Runs the inception init command
		cmd := exec.Command(command[0], command[1:]...)
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		//Show error message
		if err != nil {
			logger.Zap.Error("There was an error while doing the initialization", zap.Error(err))
			os.Exit(1)
		}
		//Try opening the new inception dir and create a file and write a file
		f, err := os.OpenFile(inception_dir+"/.terraform/environment", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logger.Zap.Error("The file that manages the environments could not be created.", zap.Error(err))
			os.Exit(1)
		}
		_, err = f.Write([]byte(default_env_name))
		if err != nil {
			logger.Zap.Error("Writing on the file that manages the environments wasn't possible.")
			os.Exit(1)
		}
		f.Close()

		//If all is ok, display successfully inception
		logger.Zap.Info("The initialization was successful!")
	} else {
		logger.Zap.Warn("The initialization was already done.")
	}

}
