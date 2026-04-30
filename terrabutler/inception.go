package main

import (
	"os"
	"os/exec"

	"go.uber.org/zap"
)

func inception_init_check() bool {
	dir := paths["inception"]
	if _, err := os.Stat(dir + "/.terraform"); !os.IsNotExist(err) {
		if _, err2 := os.Stat(dir + "/.terraform/environment"); !os.IsNotExist(err2) {
			return true
		}
	}
	return false
}

func inception_init_needed() {

	if !inception_init_check() {
		logger.Error("The initialization hasn't been made yet. Please execute the following command to initialize it: terrabutler init")
		os.Exit(1)
	}
}

func inception_init() {
	org := settings.String("general.organization")
	default_env_name := settings.String("environments.default.name")
	inception_dir := paths["inception"]
	backend_dir := paths["backends"]

	if !inception_init_check() {

		//run Terraform init -backend-config LOCATION
		os.Chdir(inception_dir)

		//Verifies that terraform exist in the current dir
		_, err := exec.LookPath("terraform")
		if err != nil {
			logger.Error("No Terraform executable found, please run mise script first.")
			os.Exit(1)
		}

		//Runs the inception init command
		cmd := exec.Command("terraform", "init", "-backend-config", backend_dir+"/"+org+"-"+default_env_name+"-inception.tfvars")
		err = cmd.Run()
		//Show error message
		if err != nil {
			logger.Error("There was an error while doing the initialization", zap.Error(err))
			os.Exit(1)
		}
		//Try opening the new inception dir and create a file and write a file
		f, err := os.OpenFile(inception_dir+"/.terraform/environment", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logger.Error("The file that manages the environments could not be created.", zap.Error(err))
			os.Exit(1)
		}
		_, err = f.Write([]byte(default_env_name))
		if err != nil {
			logger.Error("Writing on the file that manages the environments wasn't possible.")
			os.Exit(1)
		}
		f.Close()

		//If all is ok, display successfully inception
		logger.Info("The initialization was successful!")
	} else {
		logger.Warn("The initialization was already done.")
	}

}
