package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"syscall"

	"go.uber.org/zap"
)

// Used for generate-options, prints arguments
func terraform_args_print(command string, site string) string {
	var needed_options string
	if command == "init" {
		needed_options = "backend"
	} else if command == "plan" || command == "apply" {
		needed_options = "var"
	} else {
		needed_options = ""
	}

	options := terraform_needed_options_builder(needed_options, site)
	return strings.Join(options, " ")
}

// Create array of needed options for backend or var files
func terraform_needed_options_builder(needed_options string, site string) []string {
	env := get_current_env()
	org := settings.String("general.organization")
	default_env := settings.String("environments.default.name")

	if needed_options == "backend" {
		backend_dir := paths["backends"]

		if site == "inception" { //Init inception with default ENV
			return []string{"-backend-config", backend_dir + "/" + org + "-" + default_env + "-inception.tfvars"}
		} else {
			return []string{"-backend-config", backend_dir + "/" + org + "-" + env + "-" + site + ".tfvars"}
		}
	} else if needed_options == "var" {
		variables_dir := paths["variables"]

		return []string{"-var-file", variables_dir + "/global.tfvars",
			"-var-file", variables_dir + "/" + org + "-" + env + ".tfvars",
			"-var-file", variables_dir + "/" + org + "-" + env + "-" + site + ".tfvars"}

	} else { // If needed_options is empty, return empty slice
		return []string{}
	}
}

// Command builder
func terraform_command_builder(command string, site string, args []string, options []string, needed_options string) []string {

	base_command := []string{"terraform"}
	base_command = append(base_command, strings.Split(command, " ")...)

	if needed_options == "backend" || needed_options == "var" {
		aux := terraform_needed_options_builder(needed_options, site)
		base_command = append(base_command, aux...)
	}

	base_command = append(base_command, options...)
	base_command = append(base_command, args...)

	return base_command
}

// Main runner function
func terraform_command_runner(command string, site string, args []string, options []string, needed_options string) {

	env := get_current_env()

	tf_binary, err := exec.LookPath("terraform")
	if err != nil {
		logger.Error("No Terraform executable found.")
		os.Exit(1)
	}

	//Changes the current working dir to the site chosen
	err = os.Chdir(paths["root"] + "/site_" + site)
	//In theory this error shouldn't occur, since is begin parsed before the execution of this command
	if err != nil {
		logger.Error("Error in finding the path for the site " + site)
		os.Exit(1)
	}

	runner_command := terraform_command_builder(command, site, args, options, needed_options)

	logger.Debug(fmt.Sprintf("Executing command: %v", runner_command))

	runner_env := os.Environ()

	logger.Debug("Tf Binary loc: " + tf_binary)

	execErr := syscall.Exec(tf_binary, runner_command, runner_env)
	if execErr != nil {
		logger.Error("There was an error during execution of terraform "+command+" in the site "+site+" in the environment "+env, zap.Error(execErr))
		os.Exit(1)
	}

}

// New commands to be used in all sites
func tf_destroy_all_sites() {
	sites := settings.Strings("sites.ordered")
	slices.Reverse(sites)
	for _, site := range sites {
		command := terraform_command_builder("destroy", site, []string{}, []string{"-auto-approve"}, "var")

		//Changes the current working dir to the site chosen
		err := os.Chdir(paths["root"] + "/site_" + site)
		if err != nil {
			logger.Error("Error in finding the path for the site " + site)
			os.Exit(1)
		}
		cmd := exec.Command(command[0], command[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			logger.Error("There was an error during execution of terraform "+command[0]+"in the site "+site+" in the environment "+get_current_env(), zap.Error(err))
			os.Exit(1)
		}
	}

}

func tf_apply_all_sites() {
	sites := settings.Strings("sites.ordered")
	for _, site := range sites {
		if site != "inception" {
			command := terraform_command_builder("init", site, []string{}, []string{"-reconfigure"}, "backend")
			cmd := exec.Command(command[0], command[1:]...)
			cmd.Run()
		}
		command := terraform_command_builder("apply", site, []string{}, []string{"-auto-approve"}, "var")

		//Changes the current working dir to the site chosen
		err := os.Chdir(paths["root"] + "/site_" + site)
		if err != nil {
			logger.Error("Error in finding the path for the site " + site)
			os.Exit(1)
		}
		cmd := exec.Command(command[0], command[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			logger.Error("There was an error during execution of terraform "+command[0]+"in the site "+site+" in the environment "+get_current_env(), zap.Error(err))
			os.Exit(1)
		}
	}
}

func tf_init_all_sites() {
	sites := settings.Strings("sites.ordered")
	if index := slices.Index(sites, "inception"); index != -1 {
		sites = sites[index+1:]
	}
	for _, site := range sites {

		logger.Warn("Initializing " + site + " site")
		command := terraform_command_builder("init", site, []string{}, []string{"-reconfigure"}, "backend")

		//Changes the current working dir to the site chosen
		err := os.Chdir(paths["root"] + "/site_" + site)
		if err != nil {
			logger.Error("Error in finding the path for the site " + site)
			os.Exit(1)
		}
		cmd := exec.Command(command[0], command[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			logger.Error("There was an error during execution of terraform "+command[0]+"in the site "+site+" in the environment "+get_current_env(), zap.Error(err))
			os.Exit(1)
		}
	}
}
