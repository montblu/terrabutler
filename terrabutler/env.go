package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"go.uber.org/zap"
)

func confirmation_menu(question string) bool {
	fmt.Printf("%s [y/N]: ", question)
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		logger.Error("Fail to read input.", zap.Error(err))
	}
	input = strings.ToLower(strings.TrimSpace(input))
	if input == "y" || input == "yes" {
		return true
	}
	return false
}

func get_current_env() string {

	//Open site environment file
	env, err := os.ReadFile(paths["environment"])
	if err != nil {
		logger.Error("An error has occurred:", zap.Error(err))
		os.Exit(1)
	}
	logger.Debug("Current environment is " + string(env))
	return string(env)
}

func is_protected_env(env string) bool {

	if slices.Contains(settings.Strings("environments.permanent"), env) {
		return true
	}
	return false
}

func get_available_envs() []string {

	//List of env
	var envs []string

	directory := paths["inception"]
	org := settings.String("general.organization")
	default_env_name := settings.String("environments.default.name")

	os.Chdir(directory)

	cmd := exec.Command("terraform", "init", "-backend-config", paths["backend"]+"/"+org+"-"+default_env_name+"-inception.tfvars")
	err := cmd.Run()
	if err != nil {
		logger.Error("There was an error from terraform init for "+org+"-"+default_env_name+" environment.", zap.Error(err))
		os.Exit(1)
	}

	cmd = exec.Command("terraform", "workspace", "list")
	workspaces, err := cmd.Output()
	if err != nil {
		logger.Error("There was an error from terraform workspace for "+org+"-"+default_env_name+" environment.", zap.Error(err))
		os.Exit(1)
	}
	//After getting the output string, its needed to
	//Remove the new line
	workspaces_trim := strings.ReplaceAll(string(workspaces), "\n", "")
	//Remove '*' and substitute with " " to still have the same number of " "
	workspaces_trim = strings.ReplaceAll(workspaces_trim, "*", " ")
	//Split the workspaces into a slice flag (Normally between each one exists 2 spaces)
	workspaces_split := strings.Split(workspaces_trim, "  ")
	if len(workspaces_split) > 1 {
		//Get the size of workspace
		size := len(workspaces_split)
		//Remove the first two members of the workspace (" " and default)
		workspaces_split = workspaces_split[2:size]
	}
	envs = append(envs, workspaces_split...)

	return envs
}

func set_current_env(env string, init bool) {

	current_env := get_current_env()
	available_envs := get_available_envs()

	//Check if env is the current in use
	if env == current_env {
		logger.Warn("The environment you selected is the current one.")
		logger.Warn("No changes were made.")
		os.Exit(0)
	}
	//Check if env exists
	if !slices.Contains(available_envs, env) {
		logger.Error("The environment '" + env + "' does not exist.")
		logger.Error("You can create this environment with the 'new' command.")
		os.Exit(1)
	}
	pre_hook := settings.String("hooks.pre_env_select")

	if pre_hook != "" {
		command := strings.Split(pre_hook, " ")
		cmd := exec.Command(command[0], command[1:len(command)]...)
		cmd.Env = []string{
			"TERRABUTLER_OLD_ENV=" + current_env,
			"TERRABUTLER_NEW_ENV=" + env}
		err := cmd.Run()
		if err != nil {
			logger.Error("pre_env_select hook failed:", zap.Error(err))
			os.Exit(1)
		}

	}
	//Get the pre_hook -> "pre_env_select", if exists
	//Run pre_hook [CHANGE ENVIRONMENT: "TERRABUTLER_OLD_ENV": current_env, "TERRABUTLER_NEW_ENV": env]
	//Show error if it occurs

	//Try opening the file in path environments and writing the new env
	//Show error if it fails
	f, err := os.Create(paths["environment"])
	if err != nil {
		logger.Error(fmt.Sprint("An error has occurred opening the file: ", err))
		os.Exit(1)
	}
	l, err := f.Write([]byte(env))
	if l == 0 && err != nil {
		logger.Error(fmt.Sprint("An error has occurred writing to the file: ", err))
		f.Close()
		os.Exit(1)
	}
	err = f.Close()

	//If init true, run terraform_init_all_sites
	if init {
		tf_init_all_sites()
	}

	//Get the post_hook -> "post_env_select", if exists
	//Run post_hook [CHANGE ENVIRONMENT: "TERRABUTLER_OLD_ENV": current_env, "TERRABUTLER_NEW_ENV": env]
	//Show error if it occurs
	post_hook := settings.String("hooks.post_env_select")

	if post_hook != "" {
		command := strings.Split(post_hook, " ")
		cmd := exec.Command(command[0], command[1:len(command)]...)
		cmd.Env = []string{
			"TERRABUTLER_OLD_ENV=" + current_env,
			"TERRABUTLER_NEW_ENV=" + env}
		err := cmd.Run()
		if err != nil {
			logger.Error("post_env_select hook failed:", zap.Error(err))
			os.Exit(1)
		}

	}

	//Show successfully message at the end
	logger.Info("Switched to environment '" + env + "'.")

}

func delete_env(env string, confirmation bool, destroy bool) {

	org := settings.String("general.organization")

	//Check if env does exist
	if !slices.Contains(get_available_envs(), env) {
		logger.Error("The environment you are trying to delete does not exist.")
		logger.Error("No changes were made.")
		os.Exit(1)
	}
	//Check if the env is the current in use
	if env == get_current_env() {
		logger.Error("The environment you are trying to delete is your active environment.")
		logger.Error("Please switch to another workspace and try again.")
		os.Exit(1)
	}
	//Check if the env is permanent / use --> is_protected_env
	if is_protected_env(env) {
		logger.Error("The environment you are trying to delete is a permanent environment and can not be deleted.")
		logger.Error("No changes were made.")
		os.Exit(1)
	}

	//Confirmation menu, if yes
	if confirmation || confirmation_menu("Do you really want to delete '"+env+"' environment?") {

		//If destroy and is not a permanent env ^ already checked above, run tf_destroy_all_sites(env)
		if destroy {
			tf_destroy_all_sites()
		}
		//For each file in variables, remove/delete it
		files, _ := os.ReadDir(paths["variables"])
		for _, file := range files {
			if !file.IsDir() {
				fmt.Println(file.Name())
				if strings.Contains(file.Name(), org+"-"+env) {
					os.Remove(paths["variables"] + "/" + file.Name())
				}
			}
		}

		//Run the terraform workspace delete env [Path Inception] (with output and check)
		os.Chdir(paths["inception"])
		cmd := exec.Command("terraform", "workspace", "delete", env)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		//Show a error message if process not executed correctly
		if err != nil {
			logger.Error("here was an error while deleting the new environment:", zap.Error(err))
			os.Exit(1)
		}

		//In the end show a successfully executed message
		logger.Info("The environment '" + env + "' was deleted!")
	} else {
		logger.Error("Deletion cancelled.")
		os.Exit(1)
	}

}

// TODO: variables.go for the temporary environments can use the templates
func create_env(env string, confirmation bool, temporary bool, apply bool) {

	//Check if env already exists
	if slices.Contains(get_available_envs(), env) {
		logger.Warn("The environment you are trying to create already exists.")
		logger.Warn("No changes were made.")
		os.Exit(1)
	}

	//Make a confirmation menu
	//If confirmation is true execute the function
	if confirmation || confirmation_menu("Do you really want to create '"+env+"' environment?") {
		//Run Terraform workspace new env [path inception] (with output and check)
		os.Chdir(paths["inception"])
		cmd := exec.Command("terraform", "workspace", "new", env)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			logger.Error("here was an error while creating the new environment:", zap.Error(err))
			os.Exit(1)
		}

		//If temporary is true, generate the var files for the env
		if temporary {
			//generate_var_files(env)
		} else {
			//Else is a permanent environment for the list
			// Get the config file, append the new env to the config file and write the new config file
			envs := settings.Strings("environments.permanent")
			envs = append(envs, env)
			err := settings.Set("environments.permanent", envs)
			if err != nil {
				logger.Error("Error adding the new environment to the config file.", zap.Error(err))
			}
			write_settings(settings)

		}
		//If temporary and apply are true, terraform_apply_all_sites
		//Should OR apply?
		if temporary && apply {
			tf_apply_all_sites()
		}

		//In the end show a successfully executed message
		logger.Info("Created and switched to the environment '" + env + "'!")
	} else {
		logger.Error("Creation cancelled.")
		os.Exit(1)
	}

}
