package env

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/montblu/terrabutler/internal/logger"
	"github.com/montblu/terrabutler/internal/settings"
	"github.com/montblu/terrabutler/internal/tf"
	"github.com/montblu/terrabutler/internal/utils"
	"github.com/montblu/terrabutler/internal/variables"

	"github.com/spf13/afero"
)

// These functions need to be mocked for tests
var GetAvailableEnvs = getAvailableEnvs
var commandRunner = tf.CommandRunner
var runnerNoVisibleOutput = tf.RunnerNoVisibleOutput
var commandRunnerNoVisibleOutput = tf.CommandRunnerNoVisibleOutput
var initAllSites = tf.InitAllSites
var applyAllSites = tf.ApplyAllSites
var destroyAllSites = tf.DestroyAllSites

var current_env = utils.CurrentEnv

func confirmationMenu(question string, scanner bufio.Scanner) (bool, error) {
	fmt.Printf("%s [y/N]: ", question)
	scanner.Scan()
	input := scanner.Bytes()
	inputParsed := strings.ToLower(strings.TrimSpace(string(input)))
	if inputParsed == "y" || inputParsed == "yes" {
		return true, nil
	}
	return false, nil
}

func isProtectedEnv(env string) bool {
	return slices.Contains(settings.Conf.Strings("environments.permanent"), env)
}

func getAvailableEnvs(fs afero.Fs) ([]string, error) {

	// List of env
	var envs []string

	org := settings.Conf.String("general.organization")
	default_env_name := settings.Conf.String("environments.default.name")

	_, err := commandRunnerNoVisibleOutput("init", "inception", []string{}, []string{}, "backend")
	if err != nil {
		return nil, err
	}

	workspaces, err := runnerNoVisibleOutput([]string{"terraform", "workspace", "list"}, "inception", os.Environ())
	if err != nil {
		return nil, errors.New("There was an error from terraform workspace for " + org + "-" + default_env_name + " environment. Error: " + err.Error())
	}
	// After getting the output string, its needed to
	// Remove the new line
	workspaces_trim := strings.ReplaceAll(string(workspaces), "\n", "")
	// Remove '*' and substitute with " " to still have the same number of " "
	workspaces_trim = strings.ReplaceAll(workspaces_trim, "*", " ")
	// Split the workspaces into a slice flag (Normally between each one exists 2 spaces)
	workspaces_split := strings.Split(workspaces_trim, "  ")
	if len(workspaces_split) > 1 {
		// Get the size of workspace
		size := len(workspaces_split)
		// Remove the first two members of the workspace (" " and default)
		workspaces_split = workspaces_split[2:size]
	}
	envs = append(envs, workspaces_split...)

	return envs, nil
}

func SetCurrentEnv(env string, init bool, fs afero.Fs) error {

	available_envs, err := GetAvailableEnvs(fs)
	if err != nil {
		return errors.New("There was an error while getting the available environments. Error: " + err.Error())
	}

	// Check if env is the current in use
	if env == current_env {
		logger.Zap.Warn("The environment you selected is the current one.")
		logger.Zap.Warn("No changes were made.")
		return nil
	}
	// Check if env exists
	if !slices.Contains(available_envs, env) {
		logger.Zap.Error("The environment '" + env + "' does not exist.")
		return errors.New("you can create this environment with the 'new' command")
	}
	pre_hook := settings.Conf.String("hooks.pre_env_select")

	if pre_hook != "" {
		command := strings.Split(pre_hook, " ")
		envVars := []string{
			"TERRABUTLER_OLD_ENV=" + current_env,
			"TERRABUTLER_NEW_ENV=" + env}
		_, err := runnerNoVisibleOutput(command, "inception", envVars)
		if err != nil {
			return errors.New("The pre_env_select hook has failed: " + err.Error())
		}

	}
	// Get the pre_hook -> "pre_env_select", if exists
	// Run pre_hook [CHANGE ENVIRONMENT: "TERRABUTLER_OLD_ENV": current_env, "TERRABUTLER_NEW_ENV": env]
	// Show error if it occurs

	// Try opening the file in path environments and writing the new env
	// Show error if it fails
	f, err := fs.Create(utils.Paths["environment"])
	if err != nil {
		return errors.New("An error has occurred opening the environment file to update it: " + err.Error())
	}
	l, err := f.Write([]byte(env))
	if l == 0 && err != nil {
		_ = f.Close()
		return errors.New("An error has occurred writing to the environment file to update it: " + err.Error())
	}
	_ = f.Close()

	// If init true, run terraform_init_all_sites
	if init {
		if err := initAllSites(); err != nil {
			return err
		}
	}

	// Get the post_hook -> "post_env_select", if exists
	// Run post_hook [CHANGE ENVIRONMENT: "TERRABUTLER_OLD_ENV": current_env, "TERRABUTLER_NEW_ENV": env]
	// Show error if it occurs
	post_hook := settings.Conf.String("hooks.post_env_select")

	if post_hook != "" {
		command := strings.Split(post_hook, " ")
		envVars := []string{
			"TERRABUTLER_OLD_ENV=" + current_env,
			"TERRABUTLER_NEW_ENV=" + env}
		_, err := runnerNoVisibleOutput(command, "inception", envVars)
		if err != nil {
			return errors.New("The post_env_select hook has failed: " + err.Error())
		}

	}

	// Show successfully message at the end
	logger.Zap.Info("Switched to environment '" + env + "'.")
	return nil
}

func DeleteEnv(env string, confirmation bool, destroy bool, fs afero.Fs) error {

	org := settings.Conf.String("general.organization")

	envs, err := GetAvailableEnvs(fs)
	if err != nil {
		return errors.New("There was an error while getting the available environments. Error: " + err.Error())
	}

	// Check if env does exist
	if !slices.Contains(envs, env) {
		logger.Zap.Error("The environment you are trying to delete does not exist.")
		logger.Zap.Error("No changes were made.")
		return nil
	}
	// Check if the env is the current in use
	if env == current_env {
		logger.Zap.Error("The environment you are trying to delete is your active environment.")
		return errors.New("please switch to another workspace and try again")

	}
	// Check if the env is permanent / use --> is_protected_env
	if isProtectedEnv(env) {
		logger.Zap.Error("The environment you are trying to delete is a permanent environment and can not be deleted.")
		return errors.New("no changes were made")
	}

	// Confirmation menu, if yes
	if !confirmation {
		choice, err := confirmationMenu("Do you really want to delete '"+env+"' environment?", *bufio.NewScanner(os.Stdin))
		if err != nil {
			return err
		}
		if !choice {
			return errors.New("deletion cancelled")
		}
	}

	// If destroy and is not a permanent env ^ already checked above, run tf_destroy_all_sites(env)
	if destroy {
		err := destroyAllSites()
		if err != nil {
			return err
		}
	}
	// For each file in variables, remove/delete it
	files, err := afero.ReadDir(fs, utils.Paths["variables"])
	if err != nil {
		return errors.New("There was an error while reading the directory to delete the environment " + env + ". Error: " + err.Error())
	}
	for _, file := range files {
		if !file.IsDir() {
			if strings.Contains(file.Name(), org+"-"+env) {
				logger.Zap.Debug("Deleted File: " + file.Name())
				err = fs.Remove(utils.Paths["variables"] + "/" + file.Name())
				if err != nil {
					return errors.New("There was an error while removing the file " + file.Name() + " of the environment " + env + ". Error: " + err.Error())
				}
			}
		}
	}

	// Run the terraform workspace delete env [Path Inception] (with output and check)
	err = commandRunner("workspace delete "+env, "inception", []string{}, []string{}, "")
	// Show a error message if process not executed correctly
	if err != nil {
		return errors.New("There was an error while deleting the " + env + " environment: " + err.Error())
	}

	// In the end show a successfully executed message
	logger.Zap.Info("The environment '" + env + "' was deleted!")
	return nil
}

func CreateEnv(env string, confirmation bool, temporary bool, apply bool, fs afero.Fs) error {

	envs, err := GetAvailableEnvs(fs)
	if err != nil {
		return errors.New("There was an error while getting the available environments. Error: " + err.Error())
	}

	// Check if env already exists
	if slices.Contains(envs, env) {
		logger.Zap.Warn("The environment you are trying to create already exists.")
		logger.Zap.Warn("No changes were made.")
		return nil
	}

	// Make a confirmation menu
	// If confirmation is true execute the function
	if !confirmation {
		choice, err := confirmationMenu("Do you really want to create '"+env+"' environment?", *bufio.NewScanner(os.Stdin))
		if err != nil {
			return err
		}
		if !choice {
			return errors.New("creation cancelled")
		}
	}

	// Run Terraform workspace new env [path inception] (with output and check)
	err = commandRunner("workspace new "+env, "inception", []string{}, []string{}, "")
	if err != nil {
		return errors.New("There was an error while creating the new environment: " + err.Error())
	}

	// If temporary is true, generate the var files for the env
	if temporary {
		if err := variables.Generate_var_files(env, fs); err != nil {
			return err
		}
	} else {
		// Else is a permanent environment for the list
		// Get the config file, append the new env to the config file and write the new config file
		envs := settings.Conf.Strings("environments.permanent")
		envs = append(envs, env)
		err := settings.Conf.Set("environments.permanent", envs)
		if err != nil {
			return errors.New("Error adding the new environment to the config file. Error: " + err.Error())
		}
		err = settings.Write_settings(fs, settings.Conf)
		if err != nil {
			return err
		}

	}
	// If temporary and apply are true, terraform_apply_all_sites
	// Should OR apply?
	if temporary && apply {
		err := applyAllSites()
		if err != nil {
			return err
		}
	}

	// In the end show a successfully executed message
	logger.Zap.Info("Created and switched to the environment '" + env + "'!")
	return nil

}
