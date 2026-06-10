package tf

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"
	"sync"

	"github.com/montblu/terrabutler/internal/logger"
	"github.com/montblu/terrabutler/internal/settings"
	"github.com/montblu/terrabutler/internal/utils"

	"golang.org/x/term"
)

// Used for generate-options, prints arguments
func ArgsPrint(command string, site string) string {
	var needed_options string
	switch command {
	case "init":
		needed_options = "backend"
	case "plan", "apply":
		needed_options = "var"
	default:
		needed_options = ""
	}

	options := NeededOptionsBuilder(needed_options, site)
	return strings.Join(options, " ")
}

// Create array of needed options for backend or var files
func NeededOptionsBuilder(needed_options string, site string) []string {
	org := settings.Conf.String("general.organization")
	default_env := settings.Conf.String("environments.default.name")
	current_env := utils.GetCurrentEnv()

	switch needed_options {
	case "backend":
		backend_dir := utils.Paths["backends"]

		if site == "inception" { // Init inception with default ENV
			return []string{"-backend-config", backend_dir + "/" + org + "-" + default_env + "-inception.tfvars"}
		} else {
			return []string{"-backend-config", backend_dir + "/" + org + "-" + current_env + "-" + site + ".tfvars"}
		}
	case "var":
		variables_dir := utils.Paths["variables"]

		return []string{"-var-file", variables_dir + "/global.tfvars",
			"-var-file", variables_dir + "/" + org + "-" + current_env + ".tfvars",
			"-var-file", variables_dir + "/" + org + "-" + current_env + "-" + site + ".tfvars"}

	default:
		return []string{}
	}
}

// Command builder
func CommandBuilder(command string, site string, args []string, options []string, needed_options string) []string {

	base_command := []string{"terraform"}
	base_command = append(base_command, strings.Split(command, " ")...)

	if needed_options == "backend" || needed_options == "var" {
		aux := NeededOptionsBuilder(needed_options, site)
		base_command = append(base_command, aux...)
	}

	base_command = append(base_command, options...)
	base_command = append(base_command, args...)

	return base_command
}

// Main runner function, which forms a terraform command and executes it
func CommandRunner(command string, site string, args []string, options []string, needed_options string) error {

	// Verifies if terraform exists
	_, err := exec.LookPath("terraform")
	if err != nil {
		return errors.New("no Terraform executable found. Please install Terraform")
	}

	// Builds the terraform command
	runner_command := CommandBuilder(command, site, args, options, needed_options)

	// Executes the command
	return Runner(runner_command, site)

}

// Executes a command with its output on the console
func Runner(command []string, site string) error {

	// Runs the terraform command
	//nolint:gosec // the command is built from internal constants, not user input
	cmd := exec.Command(command[0], command[1:]...)
	// Changes the current directory
	cmd.Dir = utils.Paths["root"] + "/site_" + site
	// Uses the console input
	cmd.Stdin = os.Stdin
	// Prints the output to the console
	cmd.Stdout = os.Stdout
	// Prints the errors to the console
	cmd.Stderr = os.Stderr
	// Runs the command
	err := cmd.Run()
	if err != nil {
		return errors.New("There was an error during execution of terraform " + command[0] + " in the site " + site + " in the environment " + utils.GetCurrentEnv() + ", Error: " + err.Error())
	}
	return nil
}

// Runner function form a terraform commands that require no output visible
func CommandRunnerNoVisibleOutput(command string, site string, args []string, options []string, needed_options string) ([]byte, error) {

	// Verifies if terraform exists
	_, err := exec.LookPath("terraform")
	if err != nil {
		return nil, errors.New("no Terraform executable found. Please install Terraform")
	}

	// Builds the terraform command
	runner_command := CommandBuilder(command, site, args, options, needed_options)

	// Executes the command
	return RunnerNoVisibleOutput(runner_command, site, os.Environ())

}

// Execute a command with a defined environment variables and no visible output
func RunnerNoVisibleOutput(command []string, site string, envVars []string) ([]byte, error) {

	//nolint:gosec // the command is built from internal constants, not user input
	cmd := exec.Command(command[0], command[1:]...)
	// Changes the current directory
	cmd.Dir = utils.Paths["root"] + "/site_" + site
	// Defining Environment Variables
	cmd.Env = envVars
	// Enabling error output
	cmd.Stderr = os.Stderr
	// Runs the command
	output, err := cmd.Output()
	if err != nil {
		return nil, errors.New("There was an error during execution of " + strings.Join(command, " ") + " in the site " + site + " in the environment " + utils.GetCurrentEnv() + ", Error: " + err.Error())
	}
	return output, nil
}

// New commands to be used in all sites
func DestroyAllSites() error {
	sites := settings.Conf.Strings("sites.ordered")
	slices.Reverse(sites)
	for _, site := range sites {
		err := CommandRunner("destroy", site, []string{}, []string{"-auto-approve"}, "var")
		if err != nil {
			return errors.New("Error destroying all sites, during site " + site + ", Error: " + err.Error())
		}

	}
	return nil
}

func ApplyAllSites() error {
	sites := settings.Conf.Strings("sites.ordered")
	for _, site := range sites {
		if site != "inception" {
			err := CommandRunner("init", site, []string{}, []string{"-reconfigure"}, "backend")
			if err != nil {
				return errors.New("Error initializing site during apply-all, site " + site + ", Error: " + err.Error())
			}
		}
		err := CommandRunner("apply", site, []string{}, []string{"-auto-approve"}, "var")
		if err != nil {
			return errors.New("Error applying all sites, during site " + site + ", Error: " + err.Error())
		}
	}
	return nil
}

func InitAllSites() error {
	sites := settings.Conf.Strings("sites.ordered")
	// Remove "inception" from the list of sites to be initialized.
	if index := slices.Index(sites, "inception"); index != -1 {
		sites = slices.Delete(sites, index, index+1)
	}
	for _, site := range sites {

		logger.Zap.Warn("Initializing " + site + " site")
		err := CommandRunner("init", site, []string{}, []string{"-reconfigure"}, "backend")
		if err != nil {
			return errors.New("Error initializing all sites, during site " + site + ", Error: " + err.Error())
		}
	}
	return nil
}
