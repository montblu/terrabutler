package tf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"slices"
	"strings"
	"sync"
	"syscall"

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

// trapTerminationSignals prevents Go's default terminate-on-signal behavior so the
// parent can wait for the terraform child to exit cleanly, and forwards the signal
// to the child process. Do not read sigChan for any other purpose.
func trapTerminationSignals(cmd *exec.Cmd) (stop func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range sigChan {
			if cmd.Process != nil {
				_ = cmd.Process.Signal(sig)
			}
		}
	}()
	return func() {
		signal.Stop(sigChan)
		close(sigChan)
	}
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

	// Starts the command first so cmd.Process is fully assigned before the
	// signal-trap goroutine reads it, avoiding a data race with exec.Cmd.Start().
	if err := cmd.Start(); err != nil {
		return errors.New("There was an error during execution of terraform " + command[0] + " in the site " + site + " in the environment " + utils.GetCurrentEnv() + ", Error: " + err.Error())
	}

	// Trap ctrl+C and just wait for terraform
	stop := trapTerminationSignals(cmd)
	defer stop()

	// Waits for the command to finish
	err := cmd.Wait()
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
	// Captures stdout manually since cmd.Output() would call Start() internally,
	// which must happen before trapTerminationSignals is set up (see below).
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	// Starts the command first so cmd.Process is fully assigned before the
	// signal-trap goroutine reads it, avoiding a data race with exec.Cmd.Start().
	if err := cmd.Start(); err != nil {
		return nil, errors.New("There was an error during execution of " + strings.Join(command, " ") + " in the site " + site + " in the environment " + utils.GetCurrentEnv() + ", Error: " + err.Error())
	}

	// Trap ctrl+C and just wait for terraform
	stop := trapTerminationSignals(cmd)
	defer stop()

	// Waits for the command to finish
	err := cmd.Wait()
	if err != nil {
		return nil, errors.New("There was an error during execution of " + strings.Join(command, " ") + " in the site " + site + " in the environment " + utils.GetCurrentEnv() + ", Error: " + err.Error())
	}
	return stdout.Bytes(), nil
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

// Creating var for mockable function in tests
var commandRunnerNoVisibleOutputVar = CommandRunnerNoVisibleOutput

func InitAllSites() error {
	sites := settings.Conf.Strings("sites.ordered")
	// Remove "inception" from the list of sites to be initialized.
	if index := slices.Index(sites, "inception"); index != -1 {
		sites = slices.Delete(sites, index, index+1)
	}

	if len(sites) == 0 {
		return nil
	}

	total := len(sites)
	logger.Zap.Info(fmt.Sprintf("Initializing %d sites in parallel...", total))

	type result struct {
		site string
		err  error
	}

	results := make(chan result, total)
	var wg sync.WaitGroup

	for _, site := range sites {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			_, err := commandRunnerNoVisibleOutputVar("init", s, []string{}, []string{"-reconfigure"}, "backend")
			results <- result{site: s, err: err}
		}(site)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	isTerm := isTerminal(os.Stderr)
	completed := 0
	failed := 0
	var errs []error

	for r := range results {
		completed++
		if r.err != nil {
			failed++
			errs = append(errs, fmt.Errorf("site %s: %w", r.site, r.err))
		}

		if isTerm {
			drawProgressBar(completed, failed, total)
		} else {
			status := "✔"
			if r.err != nil {
				status = "✗"
			}
			fmt.Fprintf(os.Stderr, "  %s site %s\n", status, r.site)
		}
	}

	if isTerm {
		fmt.Fprintln(os.Stderr)
	}

	for _, e := range errs {
		logger.Zap.Error(e.Error())
	}

	if len(errs) > 0 {
		return fmt.Errorf("%d/%d sites failed to initialize", failed, total)
	}

	logger.Zap.Info("All sites initialized successfully")
	return nil
}

func drawProgressBar(done, failed, total int) {
	const width = 30
	filled := (done * width) / total
	bar := strings.Repeat("=", filled) + strings.Repeat(" ", width-filled)

	status := fmt.Sprintf("[%s] %d/%d sites", bar, done, total)
	if failed > 0 {
		status += fmt.Sprintf(" (%d failed)", failed)
	}
	fmt.Fprintf(os.Stderr, "\r\033[K%s", status)
}

func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(int(f.Fd()))
}
