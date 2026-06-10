package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v3"

	"github.com/montblu/terrabutler/internal/env"
	"github.com/montblu/terrabutler/internal/inception"
	"github.com/montblu/terrabutler/internal/logger"
	"github.com/montblu/terrabutler/internal/requirements"
	"github.com/montblu/terrabutler/internal/settings"
	"github.com/montblu/terrabutler/internal/tf"
	"github.com/montblu/terrabutler/internal/utils"
)

func Run(appName, version, commit, date string, fs afero.Fs) error {

	// Verify the semantic version
	_ = utils.IsSemanticVersion(version)

	// Changing the default help flag
	cli.HelpFlag = &cli.BoolFlag{
		Name:    "help",
		Aliases: []string{"H", "h"},
		Usage:   "Show this message and exit",
	}

	// Using the new HelpPrinter
	cli.HelpPrinter = HelpPrinterNewFunctions

	// Applying the new templates for the helper
	cli.RootCommandHelpTemplate = RootCommandHelpTemplate
	cli.CommandHelpTemplate = CommandHelpTemplate
	cli.SubcommandHelpTemplate = SubcommandHelpTemplate

	defer func() { _ = logger.Zap.Sync() }()

	cmd := &cli.Command{
		Name:      appName,
		Usage:     "The utility that helps keeping your IaC in one piece",
		UsageText: appName + " [OPTIONS] COMMAND [ARGS]...",
		Version:   version,
		// Hides Help Command to "Remove" HelpCommand, you need to hide it for each command
		HideHelpCommand:       true,
		HideVersion:              true,
		EnableShellCompletion:    true,
		Suggest:                  true,
		CommandNotFound:          CommandNotFound,
		OnUsageError:             OnUsageError,
		InvalidFlagAccessHandler: InvalidFlagAccessHandler,
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			if len(os.Args) > 1 && os.Args[1] == "version" {
				return ctx, nil
			}
			if err := settings.ValidateSettings(fs); err != nil {
				return ctx, err
			}
			if err := requirements.CheckRequirement(fs); err != nil {
				return ctx, err
			}
			return ctx, nil
		},
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "Show version and exit",
			Action: func(ctx context.Context, c *cli.Command) error {
				_, _ = fmt.Fprintf(c.Root().Writer, "%s %s (commit: %s, date: %s)\n", appName, version, commit, date)
				return nil
			},
				CommandNotFound:          CommandNotFound,
				OnUsageError:             OnUsageError,
				InvalidFlagAccessHandler: InvalidFlagAccessHandler,
			},
			{
				Name:                     "env",
				Usage:                    "Manage environments",
				UsageText:                appName + " env [OPTIONS] COMMAND [ARGS]...",
				HideHelp:                 true,
				Suggest:                  true,
				EnableShellCompletion:    true,
				CommandNotFound:          CommandNotFound,
				OnUsageError:             OnUsageError,
				InvalidFlagAccessHandler: InvalidFlagAccessHandler,
				Commands: []*cli.Command{
					// Subcommands of Env
					{
						Name:                     "delete",
						Aliases:                  []string{""},
						Usage:                    "Delete an environment",
						UsageText:                appName + " env delete [OPTIONS] NAME",
						ArgsUsage:                "NAME",
						HideHelp:                 true,
						Suggest:                  true,
						EnableShellCompletion:    true,
						CommandNotFound:          CommandNotFound,
						OnUsageError:             OnUsageError,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
						Arguments:                []cli.Argument{&cli.StringArg{Name: "ENV"}},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								// Added destroy as alias to the -d flag
								Name:    "destroy",
								Aliases: []string{"d"},
								Usage:   "Destroy all sites by inverse order.",
							},
							&cli.BoolFlag{
								Name:    "y",
								Aliases: []string{"Y"},
								Usage:   "Delete without asking for confirmation.",
							},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							if c.StringArg("ENV") == "" {
								return errors.New("missing argument 'NAME'")
							}
							return env.DeleteEnv(c.StringArg("ENV"), c.Bool("y"), c.Bool("d"), fs)
						}},
					{
						Name:      "list",
						Aliases:   []string{""},
						Usage:     "List environments",
						UsageText: appName + " env list [OPTIONS]",
						HideHelp:  true,
						Suggest:   true,
						CommandNotFound:          CommandNotFound,
						OnUsageError:             OnUsageError,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
						Action: func(context.Context, *cli.Command) error {
							envs, err := env.GetAvailableEnvs(fs)
							if err != nil {
								return err
							}
							for _, env := range envs {
								if env == utils.GetCurrentEnv() {
									fmt.Println("\u2192", env)
								} else {
									fmt.Println(env)
								}
							}
							return nil
						}},
					{
						Name:      "new",
						Aliases:   []string{""},
						Usage:     "Create a new environment",
						UsageText: appName + " env new [OPTIONS] NAME",
						HideHelp:  true,
						ArgsUsage: "NAME",
						Arguments: []cli.Argument{&cli.StringArg{Name: "ENV"}},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "y",
								Aliases: []string{"Y"},
								Usage:   "Delete without asking for confirmation.",
							},
							&cli.BoolFlag{
								Name:    "t",
								Aliases: []string{"temp"},
								Usage:   "Create a temporary environment.",
							},
							&cli.BoolFlag{
								Name:    "a",
								Aliases: []string{"apply"},
								Usage:   "Apply all terraform sites prior the creation of the environment.",
							},
						},
						CommandNotFound:          CommandNotFound,
						OnUsageError:             OnUsageError,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
						Action: func(ctx context.Context, c *cli.Command) error {
							if c.StringArg("ENV") == "" {
								return errors.New("missing argument 'NAME'")
							}
							return env.CreateEnv(c.StringArg("ENV"), c.Bool("y"), c.Bool("t"), c.Bool("a"), fs)
						}},
					{
						Name:                     "reload",
						Aliases:                  []string{""},
						HideHelp:                 true,
						Usage:                    "Reload the current environment",
						UsageText:                appName + " env reload [OPTIONS]",
						CommandNotFound:          CommandNotFound,
						OnUsageError:             OnUsageError,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
						Action: func(context.Context, *cli.Command) error {
							return tf.InitAllSites()
						}},
					{
						Name:      "select",
						Aliases:   []string{""},
						Usage:     "Select a environment",
						UsageText: appName + " env select [OPTIONS] NAME",
						HideHelp:  true,
						ArgsUsage: "NAME",
						Arguments: []cli.Argument{&cli.StringArg{Name: "ENV"}},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "no-init",
								Usage: "Skip auto init of the sites when selecting an environment.",
							},
						},
						CommandNotFound:          CommandNotFound,
						OnUsageError:             OnUsageError,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
						Action: func(ctx context.Context, c *cli.Command) error {
							if c.StringArg("ENV") == "" {
								return errors.New("missing argument 'NAME'")
							}
							init := !c.Bool("no-init")
							return env.SetCurrentEnv(c.StringArg("ENV"), init, fs)
						}},
					{
						Name:                     "show",
						Aliases:                  []string{""},
						HideHelp:                 true,
						Usage:                    "Show the name of the current environment",
						UsageText:                appName + " env show [OPTIONS]",
						CommandNotFound:          CommandNotFound,
						OnUsageError:             OnUsageError,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
					Action: func(ctx context.Context, c *cli.Command) error {
						_, _ = fmt.Fprintf(c.Root().Writer, "%s\n", utils.GetCurrentEnv())
						return nil
					}},
				},
				Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
					return ctx, inception.InitNeeded(fs)
				},
			},
			{
				Name:                     "init",
				Usage:                    "Initialize the manager",
				UsageText:                appName + " init [OPTIONS]",
				HideHelp:                 true,
				CommandNotFound:          CommandNotFound,
				OnUsageError:             OnUsageError,
				InvalidFlagAccessHandler: InvalidFlagAccessHandler,
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return inception.Init(fs)
				},
			},
			{
				Name:                  "tf",
				Usage:                 "Manage terraform commands",
				UsageText:             appName + " tf [OPTIONS] COMMAND [ARGS]...",
				HideHelp:              true,
				EnableShellCompletion: true,
				Suggest:               true,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "site", Required: true, Usage: "Site where to run terraform.  [required]"}},
				Commands: []*cli.Command{
					{
						Name:      "apply",
						HideHelp:  true,
						Usage:     "Create or update infrastructure.",
						UsageText: appName + " tf apply [OPTIONS]",
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "auto-approve", Usage: "Skip interactive approval of plan before applying."},
							&cli.BoolFlag{Name: "destroy", Usage: "Select the 'destroy' planning mode, which creates a plan to destroy all objects currently managed by this Terraform configuration instead of the usual behavior."},
							// Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-input", Usage: "Don't ask for input for variables if not directly set."},
							// Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-lock", Usage: "Don't hold a state lock during backend migration. This is dangerous if others might concurrently run commands against the same workspace."},
							&cli.StringFlag{Name: "lock-timeout", Usage: "Duration to retry a state lock."},
							&cli.BoolFlag{Name: "no-color", Usage: "If specified, output won't contain any color."},
							&cli.BoolFlag{Name: "refresh-only", Usage: "Select the 'refresh only' planning mode, which checks whether remote objects still match the outcome of the most recent Terraform apply but does not propose any actions to undo any changes made outside of Terraform."},
							// Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-refresh", Usage: "Skip checking for external changes to remote objects while creating the plan. This can potentially make planning faster, but at the expense of possibly planning against a stale record of the remote system state."},
							&cli.StringSliceFlag{Name: "target", Usage: "Limit the planning operation to only the given module, resource, or resource instance and all of its dependencies. You can use this option multiple times to include more than one object. This is for exceptional use only."},
							&cli.StringSliceFlag{Name: "var", Usage: "Set a value for one of the input variables in the root module of the configuration. Use this option more than once to set more than one variable."},
						},
						OnUsageError:             OnUsageErrorSite,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
						Action: func(ctx context.Context, c *cli.Command) error {
							options := []string{}
							if c.Bool("auto-approve") {
								options = append(options, "-auto-approve")
							}
							if c.Bool("destroy") {
								options = append(options, "-destroy")
							}
							if c.Bool("no-input") {
								options = append(options, "-input=false")
							}
							if c.Bool("no-lock") {
								options = append(options, "-lock=false")
							}
							if c.String("lock-timeout") != "" {
								options = append(options, "-lock-timeout="+c.String("lock-timeout"))
							}
							if c.Bool("no-color") {
								options = append(options, "-no-color")
							}
							if c.Bool("refresh-only") {
								options = append(options, "-refresh-only")
							}
							if c.Bool("no-refresh") {
								options = append(options, "-refresh=false")
							}
							for _, target := range c.StringSlice("target") {
								options = append(options, "-target="+target)
							}
							for _, v := range c.StringSlice("var") {
								options = append(options, "-var="+v)
							}
							return tf.CommandRunner("apply", c.String("site"), []string{}, options, "var")
						},
					},
					{
						Name:      "console",
						HideHelp:  true,
						Usage:     "Try Terraform expressions at an interactive command...",
						UsageText: appName + " tf console [OPTIONS]",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "state", Usage: "Legacy option for the local backend only. See the local backend's documentation for more information."},
							&cli.BoolFlag{Name: "plan", Usage: "Create a new plan (as if running \"terraform plan\") and then evaluate expressions against its planned state, instead of evaluating against the current state. You can use this to inspect the effects of configuration changes that haven't been applied yet.."},
							&cli.StringSliceFlag{Name: "var", Usage: "Set a variable in the Terraform configuration. This flag can be set multiple times."},
						},
						OnUsageError:             OnUsageErrorSite,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
						Action: func(ctx context.Context, c *cli.Command) error {
							options := []string{}
							if c.String("state") != "" {
								options = append(options, "-state="+c.String("state"))
							}
							if c.Bool("plan") {
								options = append(options, "-plan")
							}
							for _, v := range c.StringSlice("var") {
								options = append(options, "-var="+v)
							}
							return tf.CommandRunner("console", c.String("site"), []string{}, options, "var")
						},
					},
					{
						Name:      "destroy",
						HideHelp:  true,
						Usage:     "Prepare your working directory for other commands",
						UsageText: appName + " tf destroy [OPTIONS]",
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "auto-approve", Usage: "Skip interactive approval of plan before applying."},
							// Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-input", Usage: "Don't ask for input for variables if not directly set."},
							// Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-lock", Usage: "Don't hold a state lock during backend migration. This is dangerous if others might concurrently run commands against the same workspace."},
							&cli.StringFlag{Name: "lock-timeout", Usage: "Duration to retry a state lock."},
							&cli.BoolFlag{Name: "no-color", Usage: "If specified, output won't contain any color."},
							&cli.BoolFlag{Name: "refresh-only", Usage: "Select the 'refresh only' planning mode, which checks whether remote objects still match the outcome of the most recent Terraform apply but does not propose any actions to undo any changes made outside of Terraform."},
							// Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-refresh", Usage: "Skip checking for external changes to remote objects while creating the plan. This can potentially make planning faster, but at the expense of possibly planning against a stale record of the remote system state."},
							&cli.StringSliceFlag{Name: "target", Usage: "Limit the planning operation to only the given module, resource, or resource instance and all of its dependencies. You can use this option multiple times to include more than one object. This is for exceptional use only."},
							&cli.StringSliceFlag{Name: "var", Usage: "Set a value for one of the input variables in the root module of the configuration. Use this option more than once to set more than one variable."},
						},
						OnUsageError:             OnUsageErrorSite,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
						Action: func(ctx context.Context, c *cli.Command) error {
							options := []string{}
							if c.Bool("auto-approve") {
								options = append(options, "-auto-approve")
							}
							if c.Bool("no-lock") {
								options = append(options, "-lock=false")
							}
							if c.String("lock-timeout") != "" {
								options = append(options, "-lock-timeout="+c.String("lock-timeout"))
							}
							if c.Bool("no-color") {
								options = append(options, "-no-color")
							}
							if c.Bool("refresh-only") {
								options = append(options, "-refresh-only")
							}
							if c.Bool("no-refresh") {
								options = append(options, "-refresh=false")
							}
							for _, target := range c.StringSlice("target") {
								options = append(options, "-target="+target)
							}
							for _, v := range c.StringSlice("var") {
								options = append(options, "-var="+v)
							}
							return tf.CommandRunner("destroy", c.String("site"), []string{}, options, "var")
						},
					},
					{
						Name:      "fmt",
						HideHelp:  true,
						Usage:     "Reformat your configuration in the standardstyle",
						UsageText: appName + " tf fmt [OPTIONS]",
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "diff", Usage: "Display diffs of formatting changes."},
							&cli.BoolFlag{Name: "no-color", Usage: "If specified, output won't contain any color."},
							&cli.BoolFlag{Name: "recursive", Usage: "Also process files in subdirectories. By default, only the given directory (or current directory) is processed."},
						},
						OnUsageError:             OnUsageErrorSite,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
						Action: func(ctx context.Context, c *cli.Command) error {
							options := []string{}
							if c.Bool("diff") {
								options = append(options, "-diff")
							}
							if c.Bool("no-color") {
								options = append(options, "-no-color")
							}
							if c.Bool("recursive") {
								options = append(options, "-recursive")
							}
							return tf.CommandRunner("fmt", c.String("site"), []string{}, options, "")
						},
					},
					{
						Name:      "force-unlock",
						HideHelp:  true,
						Usage:     "Release a stuck lock on the current workspace",
						UsageText: appName + " tf force-unlock [OPTIONS] LOCK_ID",
						ArgsUsage: "LOCK_ID",
						Arguments: []cli.Argument{&cli.StringArg{Name: "LOCK-ID", Value: ""}},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "force",
								Usage: "Don't ask for input for unlock confirmation.",
							},
						},
						OnUsageError:             OnUsageErrorSite,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
						Action: func(ctx context.Context, c *cli.Command) error {
							options := []string{}
							if c.StringArg("LOCK-ID") == "" {
								return errors.New("missing argument 'LOCK_ID'")

							}
							args := append([]string{}, c.StringArg("LOCK-ID"))
							if c.Bool("force") {
								options = append(options, "-force")
							}
							return tf.CommandRunner("force-unlock", c.String("site"), args, options, "")
						},
					},
					{
						Name:      "generate-options",
						HideHelp:  true,
						Usage:     "Generate terraform options",
						UsageText: appName + " tf generate-options [OPTIONS] {init|plan|apply}",
						ArgsUsage: "{init|plan|apply}",
						Arguments: []cli.Argument{
							&cli.StringArg{Name: "Choice"},
						},
						OnUsageError:             OnUsageErrorSite,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
						Action: func(ctx context.Context, c *cli.Command) error {
							if c.StringArg("Choice") != "init" && c.StringArg("Choice") != "plan" && c.StringArg("Choice") != "apply" {
								return errors.New("missing argument '{init|plan|apply}' choose one of the choices: init, plan or apply")
							}
							logger.Zap.Info("Options:\n" + tf.ArgsPrint(c.StringArg("Choice"), c.String("site")))
							return nil
						},
					},
					{
						Name:      "import",
						HideHelp:  true,
						Usage:     "Associate existing infrastructure with a Terraform...",
						UsageText: appName + " tf import [OPTIONS] ADDR ID",
						ArgsUsage: "ADDR ID",
						Arguments: []cli.Argument{
							&cli.StringArg{Name: "ADDR"},
							&cli.StringArg{Name: "ID"},
						},
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "allow-missing-config", Usage: "Allow import when no resource configuration block exists."},
							// Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-input", Usage: "Don't ask for input for variables if not directly set."},
							// Requires BOOLEAN value -- Reversing
							&cli.BoolFlag{Name: "no-lock", Usage: "Don't hold a state lock during the operation. This is dangerous if others might concurrently run commands against the same workspace."},
							&cli.BoolFlag{Name: "no-color", Usage: "If specified, output won't contain any color."},
							&cli.StringSliceFlag{Name: "var", Usage: "Set a variable in the Terraform configuration. This flag can be set multiple times."},
							&cli.BoolFlag{Name: "ignore-remote-version", Usage: "A rare option used for the remote backend only. See the remote backend documentation for more information."},
						},
						OnUsageError:             OnUsageErrorSite,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
						Action: func(ctx context.Context, c *cli.Command) error {
							options := []string{}
							if c.StringArg("ADDR") == "" {
								return errors.New("missing argument 'ADDR'")
							}
							if c.StringArg("ID") == "" {
								return errors.New("missing argument 'ID'")
							}
							args := append([]string{}, c.StringArg("ADDR"), c.StringArg("ID"))
							if c.Bool("allow-missing-config") {
								options = append(options, "-allow-missing-config")
							}
							if c.Bool("no-input") {
								options = append(options, "-input=false")
							}
							if c.Bool("no-lock") {
								options = append(options, "-lock=false")
							}
							if c.Bool("no-color") {
								options = append(options, "-no-color")
							}
							for _, v := range c.StringSlice("var") {
								options = append(options, "-var="+v)
							}
							if c.Bool("ignore-remote-version") {
								options = append(options, "-ignore-remote-version")
							}
							return tf.CommandRunner("import", c.String("site"), args, options, "var")
						},
					},
					{
						Name:      "init",
						HideHelp:  true,
						Usage:     "Prepare your working directory for other commands",
						UsageText: appName + " tf init [OPTIONS]",
						Flags: []cli.Flag{
							// Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-backend", Usage: "Disable backend or Terraform Cloud initialization for this configuration and use what what was previously initialized instead."},
							&cli.BoolFlag{Name: "force-copy", Usage: "Allow import when no resource configuration block exists."},
							// Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-get", Usage: "Disable downloading modules for this configuration."},
							// Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-input", Usage: "Disable interactive prompts. Note that some actions may require interactive prompts and will error if input is disabled."},
							// Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-lock", Usage: "Don't hold a state lock during backend migration. This is dangerous if others might concurrently run commands against the same workspace."},
							&cli.BoolFlag{Name: "no-color", Usage: "If specified, output won't contain any color."},
							&cli.BoolFlag{Name: "reconfigure", Usage: "Reconfigure a backend, ignoring any saved configuration."},
							&cli.BoolFlag{Name: "migrate-state", Usage: "Reconfigure a backend, and attempt to migrate any existing state."},
							&cli.BoolFlag{Name: "upgrade", Usage: "Install the latest module and provider versions allowed within configured constraints, overriding the default behavior of selecting exactly the version recorded in the dependency lockfile."},
							&cli.StringFlag{Name: "lockfile", Usage: "Set a dependency lockfile mode. Currently only 'readonly' is valid."},
							&cli.BoolFlag{Name: "ignore-remote-version", Usage: "A rare option used for Terraform Cloud and the remote backend only. Set this to ignore checking that the local and remote Terraform versions use compatible state representations, making an operation proceed even when there is a potential mismatch. See the documentation on configuring Terraform with Terraform Cloud for more information."},
						},
						OnUsageError:             OnUsageErrorSite,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
						Action: func(ctx context.Context, c *cli.Command) error {
							options := []string{}
							if c.Bool("no-backend") {
								options = append(options, "-backend=false")
							}
							if c.Bool("force-copy") {
								options = append(options, "-force-copy")
							}
							if c.Bool("no-get") {
								options = append(options, "-get=false")
							}
							if c.Bool("no-input") {
								options = append(options, "-input=false")
							}
							if c.Bool("no-lock") {
								options = append(options, "-lock=false")
							}
							if c.Bool("no-color") {
								options = append(options, "-no-color")
							}
							if c.Bool("reconfigure") {
								options = append(options, "-reconfigure")
							}
							if c.Bool("migrate-state") {
								options = append(options, "-migrate-state")
							}
							if c.Bool("upgrade") {
								options = append(options, "-upgrade")
							}
							if c.String("lockfile") != "" {
								options = append(options, "-lockfile="+c.String("lockfile"))
							}
							if c.Bool("ignore-remote-version") {
								options = append(options, "-ignore-remote-version")
							}
							return tf.CommandRunner("init", c.String("site"), []string{}, options, "backend")
						},
					},
					{
						Name:      "output",
						HideHelp:  true,
						Usage:     "Show output values from your root module",
						UsageText: appName + " tf output [OPTIONS]",
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "no-color", Usage: "If specified, output won't contain any color."},
							&cli.BoolFlag{Name: "json", Usage: "If specified, machine readable output will be printed in JSON format."},
							&cli.BoolFlag{Name: "raw", Usage: "For value types that can be automatically converted to a string, will print the raw string directly, rather than a human-oriented representation of the value."},
						},
						OnUsageError:             OnUsageErrorSite,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
						Action: func(ctx context.Context, c *cli.Command) error {
							options := []string{}
							if c.Bool("no-color") {
								options = append(options, "-no-color")
							}
							if c.Bool("json") {
								options = append(options, "-json")
							}
							if c.Bool("raw") {
								options = append(options, "-raw")
							}
							return tf.CommandRunner("output", c.String("site"), []string{}, options, "")
						},
					},
					{
						Name:      "plan",
						HideHelp:  true,
						Usage:     "Show changes required by the current configuration",
						UsageText: appName + " tf plan [OPTIONS]",
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "destroy", Usage: "Select the 'destroy' planning mode, which creates a plan to destroy all objects currently managed by this"},
							// Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-input", Usage: "Don't ask for input for variables if not directly set."},
							// Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-lock", Usage: "Don't hold a state lock during backend migration. This is dangerous if others might concurrently run commands against the same workspace."},
							&cli.StringFlag{Name: "lock-timeout", Usage: "Duration to retry a state lock."},
							&cli.BoolFlag{Name: "no-color", Usage: "If specified, output won't contain any color."},
							&cli.BoolFlag{Name: "refresh-only", Usage: "Select the 'refresh only' planning mode, which checks whether remote objects still match the outcome of the most recent Terraform apply but does not propose any actions to undo any changes made outside of Terraform."},
							// Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-refresh", Usage: "Skip checking for external changes to remote objects while creating the plan. This can potentially make planning faster, but at the expense of possibly planning against a stale record of the remote system state."},
							&cli.StringSliceFlag{Name: "target", Usage: "Limit the planning operation to only the given module, resource, or resource instance and all of its dependencies. You can use this option multiple times to include more than one object. This is for exceptional use only."},
							&cli.StringSliceFlag{Name: "var", Usage: "Set a value for one of the input variables in the root module of the configuration. Use this option more than once to set more than one variable."},
							&cli.StringFlag{Name: "out", Usage: "Write a plan file to the given path. This can be used as input to the \"apply\" command."},
						},
						OnUsageError:             OnUsageErrorSite,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
						Action: func(ctx context.Context, c *cli.Command) error {
							options := []string{}
							if c.Bool("destroy") {
								options = append(options, "-destroy")
							}
							if c.Bool("no-input") {
								options = append(options, "-input=false")
							}
							if c.Bool("no-lock") {
								options = append(options, "-lock=false")
							}
							if c.String("lock-timeout") != "" {
								options = append(options, "-lock-timeout="+c.String("lock-timeout"))
							}
							if c.Bool("no-color") {
								options = append(options, "-no-color")
							}
							if c.Bool("refresh-only") {
								options = append(options, "-refresh-only")
							}
							if c.Bool("no-refresh") {
								options = append(options, "-refresh=false")
							}
							for _, target := range c.StringSlice("target") {
								options = append(options, "-target="+target)
							}
							for _, v := range c.StringSlice("var") {
								options = append(options, "-var="+v)
							}
							if c.String("out") != "" {
								options = append(options, "-out="+c.String("out"))
							}

							return tf.CommandRunner("plan", c.String("site"), []string{}, options, "var")
						},
					},
					{
						Name:      "providers",
						HideHelp:  true,
						Usage:     "Show the providers required for this configuration",
						UsageText: appName + " tf providers [OPTIONS] COMMAND [ARGS]...",
						Commands: []*cli.Command{
							{
								Name:      "lock",
								Usage:     "Write out dependency locks for the configured providers",
								ArgsUsage: "PROVIDERS...",
								UsageText: appName + " tf providers lock [OPTIONS] PROVIDERS...",
								Arguments: []cli.Argument{
									&cli.StringArgs{Max: -1, Name: "Providers"},
								},
								Flags: []cli.Flag{
									&cli.StringFlag{Name: "fs-mirror", Usage: "Consult the given filesystem mirror directory instead of the origin registry for each of the given providers."},
									&cli.StringFlag{Name: "net-mirror", Usage: "Consult the given network mirror (given as a base URL) instead of the origin registry for each of the given providers."},
									&cli.StringFlag{Name: "platform", Usage: "Choose a target platform to request package checksums for."},
								},
								HideHelp:                 true,
								OnUsageError:             OnUsageErrorSite,
								InvalidFlagAccessHandler: InvalidFlagAccessHandler,
								Action: func(ctx context.Context, c *cli.Command) error {
									options := []string{}
									args := []string{}
									if len(c.StringArgs("Providers")) == 0 {
										return errors.New("missing arguments 'PROVIDERS...'")
									}
									args = append(args, c.StringArgs("Providers")...)
									if c.String("fs-mirror") != "" {
										options = append(options, "-fs-mirror="+c.String("fs-mirror"))
									}
									if c.String("net-mirror") != "" {
										options = append(options, "-net-mirror="+c.String("net-mirror"))
									}
									if c.String("platform") != "" {
										options = append(options, "-platform="+c.String("platform"))
									}
									return tf.CommandRunner("providers lock", c.String("site"), args, options, "")
								},
							},
							{
								Name:      "mirror",
								Usage:     "Save local copies of all required provider plugins",
								ArgsUsage: "TARGET_DIR",
								UsageText: appName + " tf providers mirror [OPTIONS] TARGET_DIR",
								Arguments: []cli.Argument{
									&cli.StringArg{Name: "DIR"},
								},
								Flags: []cli.Flag{
									&cli.StringFlag{Name: "platform", Usage: "Choose a target platform to request package checksums for."},
								},
								HideHelp:                 true,
								OnUsageError:             OnUsageErrorSite,
								InvalidFlagAccessHandler: InvalidFlagAccessHandler,
								Action: func(ctx context.Context, c *cli.Command) error {
									options := []string{}
									args := []string{}
									if c.StringArg("DIR") == "" {
										return errors.New("missing argument 'TARGET_DIR'")
									}
									args = append(args, c.StringArg("DIR"))
									if c.String("platform") != "" {
										options = append(options, "-platform="+c.String("platform"))
									}
									return tf.CommandRunner("providers mirror", c.String("site"), args, options, "")
								},
							},
							{
								Name:      "schema",
								Usage:     "Show schemas for the providers used in the configuration",
								UsageText: appName + " tf providers schema [OPTIONS]",
								Flags: []cli.Flag{
									&cli.BoolFlag{Name: "json", Required: true, Usage: "Prints out a json representation of the schemas for all providers used in the current configuration.  [required]"},
								},
								HideHelp:                 true,
								OnUsageError:             OnUsageErrorSite,
								InvalidFlagAccessHandler: InvalidFlagAccessHandler,
								Action: func(ctx context.Context, c *cli.Command) error {
									options := []string{}
									if c.Bool("json") {
										options = append(options, "-json")
									}
									return tf.CommandRunner("providers schema", c.String("site"), []string{}, options, "")
								},
							},
						},
						CommandNotFound:          CommandNotFound,
						OnUsageError:             OnUsageErrorSite,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
					},
					{
						Name:      "refresh",
						HideHelp:  true,
						Usage:     "Update the state to match remote systems",
						UsageText: appName + " tf refresh [OPTIONS]",
						Flags: []cli.Flag{
							// Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-input", Usage: "Don't ask for input for variables if not directly set."},
							// Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-lock", Usage: "Don't hold a state lock during the operation. This is dangerous if others might concurrently run commands against the same workspace."},
							&cli.BoolFlag{Name: "no-color", Usage: "If specified, output won't contain any color."},
							&cli.StringSliceFlag{Name: "target", Usage: "Resource to target. Operation will be limited to this resource and its dependencies. This flag can be used multiple times."},
							&cli.StringSliceFlag{Name: "var", Usage: "Set a variable in the Terraform configuration. This flag can be set multiple times."},
						},
						CommandNotFound:          CommandNotFound,
						OnUsageError:             OnUsageErrorSite,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
						Action: func(ctx context.Context, c *cli.Command) error {
							options := []string{}
							if c.Bool("no-input") {
								options = append(options, "-input=false")
							}
							if c.Bool("no-lock") {
								options = append(options, "-lock=false")
							}
							if c.Bool("no-color") {
								options = append(options, "-no-color")
							}
							for _, target := range c.StringSlice("target") {
								options = append(options, "-target="+target)
							}
							for _, v := range c.StringSlice("var") {
								options = append(options, "-var="+v)
							}
							return tf.CommandRunner("refresh", c.String("site"), []string{}, options, "var")
						},
					},
					{
						Name:      "show",
						HideHelp:  true,
						Usage:     "Show the current state or a saved plan",
						UsageText: appName + " tf show [OPTIONS] [PATH]",
						ArgsUsage: "[PATH]",
						Arguments: []cli.Argument{
							&cli.StringArg{Name: "PATH"},
						},
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "no-color", Usage: "If specified, output won't contain any color."},
							&cli.BoolFlag{Name: "json", Usage: "Output the version information as a JSON object."},
						},
						CommandNotFound:          CommandNotFound,
						OnUsageError:             OnUsageErrorSite,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
						Action: func(ctx context.Context, c *cli.Command) error {
							options := []string{}
							args := []string{}
							if c.StringArg("PATH") != "" {
								args = append(args, c.StringArg("PATH"))
							}
							if c.Bool("json") {
								options = append(options, "-json")
							}
							if c.Bool("no-color") {
								options = append(options, "-no-color")
							}
							return tf.CommandRunner("show", c.String("site"), args, options, "")
						},
					},
					{
						Name:      "state",
						HideHelp:  true,
						Usage:     "Advanced state management",
						UsageText: appName + " tf state [OPTIONS] COMMAND [ARGS]...",
						Commands: []*cli.Command{
							{
								Name:      "list",
								Usage:     "List resources in the state",
								ArgsUsage: "[ADDRESS]",
								Arguments: []cli.Argument{
									&cli.StringArg{Name: "ADDR"},
								},
								Flags: []cli.Flag{
									&cli.StringFlag{Name: "state", Usage: "Path to a Terraform state file to use to look up Terraform-managed resources. By default, Terraform will consult the state of the currently-selected workspace."},
									&cli.StringFlag{Name: "id", Usage: "Filters the results to include only instances whose resource types have an attribute named 'id' whose value equals the given id string."},
								},
								OnUsageError:             OnUsageErrorSite,
								InvalidFlagAccessHandler: InvalidFlagAccessHandler,
								Action: func(ctx context.Context, c *cli.Command) error {
									options := []string{}
									args := []string{}
									if c.StringArg("ADDR") != "" {
										args = append(args, c.StringArg("ADDR"))
									}
									if c.String("state") != "" {
										options = append(options, "-state "+c.String("state"))
									}
									if c.String("id") != "" {
										options = append(options, "-id "+c.String("id"))
									}
									return tf.CommandRunner("state list", c.String("site"), args, options, "")
								},
							},
							{
								Name:      "mv",
								Usage:     "Move an item in the state",
								ArgsUsage: "SOURCE DESTINATION",
								Arguments: []cli.Argument{
									&cli.StringArg{Name: "SOURCE"},
									&cli.StringArg{Name: "DESTINATION"},
								},
								Flags: []cli.Flag{
									&cli.BoolFlag{Name: "dry-run", Usage: "If set, prints out what would've been moved but doesn't actually move anything."},
									&cli.BoolFlag{Name: "no-lock", Usage: "Don't hold a state lock during the operation. This is dangerous if others might concurrently run commands against the same workspace."},
									&cli.StringFlag{Name: "lock-timeout", Usage: "Duration to retry a state lock."},
									&cli.BoolFlag{Name: "ignore-remote-version", Usage: "A rare option used for the remote backend only. See the remote backend documentation for more information."},
								},
								OnUsageError:             OnUsageErrorSite,
								InvalidFlagAccessHandler: InvalidFlagAccessHandler,
								Action: func(ctx context.Context, c *cli.Command) error {
									options := []string{}
									args := []string{}
									if c.StringArg("SOURCE") == "" {
										return errors.New("missing argument 'SOURCE'")
									}
									if c.StringArg("DESTINATION") == "" {
										return errors.New("missing argument 'DESTINATION'")
									}
									args = append(args, c.String("SOURCE"), c.String("DESTINATION"))
									if c.Bool("dry-run") {
										options = append(options, "-dry-run")
									}
									if c.Bool("no-lock") {
										options = append(options, "-lock=false")
									}
									if c.String("lock-timeout") != "" {
										options = append(options, "-lock-timeout="+c.String("lock-timeout"))
									}
									if c.Bool("ignore-remote-version") {
										options = append(options, "-ignore-remote-version")
									}
									return tf.CommandRunner("state mv", c.String("site"), args, options, "")
								},
							},
							{
								Name:         "pull",
								Usage:        "Pull current state and output to stdouts",
								OnUsageError: OnUsageErrorSite,
								Action: func(ctx context.Context, c *cli.Command) error {
									return tf.CommandRunner("state", c.String("site"), []string{"pull"}, []string{}, "")
								},
							},
							{
								Name:      "push",
								Usage:     "Update remote state from a local state file",
								ArgsUsage: "PATH",
								Arguments: []cli.Argument{
									&cli.StringArg{Name: "PATH"},
								},
								Flags: []cli.Flag{
									&cli.BoolFlag{Name: "force", Usage: "Write the state even if lineages don't match or the remote serial is higher."},
									&cli.BoolFlag{Name: "no-lock", Usage: "Don't hold a state lock during the operation. This is dangerous if others might concurrently run commands against the same workspace."},
									&cli.StringFlag{Name: "lock-timeout", Usage: "Duration to retry a state lock."},
								},
								OnUsageError:             OnUsageErrorSite,
								InvalidFlagAccessHandler: InvalidFlagAccessHandler,
								Action: func(ctx context.Context, c *cli.Command) error {
									options := []string{}
									args := []string{}
									if c.StringArg("PATH") == "" {
										return errors.New("missing argument 'PATH'")
									}
									args = append(args, c.StringArg("PATH"))
									if c.Bool("force") {
										options = append(options, "-force")
									}
									if c.Bool("no-lock") {
										options = append(options, "-lock=false")
									}
									if c.String("lock-timeout") != "" {
										options = append(options, "-lock-timeout="+c.String("lock-timeout"))
									}
									return tf.CommandRunner("state push", c.String("site"), args, options, "")
								},
							},
							{
								Name:      "replace-provider",
								Usage:     "Replace provider for resources in the Terraform state.",
								ArgsUsage: "FROM_PROVIDER_FQDN TO_PROVIDER_FQDN",
								Arguments: []cli.Argument{
									&cli.StringArg{Name: "FROM_FQDN"},
									&cli.StringArg{Name: "TO_FQDN"},
								},
								Flags: []cli.Flag{
									&cli.BoolFlag{Name: "auto-approve", Usage: "Skip interactive approval of plan before applying."},
									&cli.BoolFlag{Name: "no-lock", Usage: "Don't hold a state lock during the operation. This is dangerous if others might concurrently run commands against the same workspace."},
									&cli.StringFlag{Name: "lock-timeout", Usage: "Duration to retry a state lock."},
									&cli.BoolFlag{Name: "ignore-remote-version", Usage: "A rare option used for the remote backend only. See the remote backend documentation for more information."},
								},
								OnUsageError:             OnUsageErrorSite,
								InvalidFlagAccessHandler: InvalidFlagAccessHandler,
								Action: func(ctx context.Context, c *cli.Command) error {
									options := []string{}
									args := []string{}
									if c.StringArg("FROM_FQDN") == "" {
										return errors.New("missing argument 'FROM_PROVIDER_FQDN'")
									}
									if c.StringArg("TO_FQDN") == "" {
										return errors.New("missing argument 'TO_PROVIDER_FQDN'")
									}
									args = append(args, c.StringArg("FROM_FQDN"), c.StringArg("TO_FQDN"))
									if c.Bool("auto-approve") {
										options = append(options, "-auto-approve")
									}
									if c.Bool("no-lock") {
										options = append(options, "-lock=false")
									}
									if c.String("lock-timeout") != "" {
										options = append(options, "-lock-timeout="+c.String("lock-timeout"))
									}
									if c.Bool("ignore-remote-version") {
										options = append(options, "-ignore-remote-version")
									}
									return tf.CommandRunner("state replace-provider", c.String("site"), args, options, "")
								},
							},
							{
								Name:      "rm",
								Usage:     "Remove instances from the state",
								ArgsUsage: "ADDRESS...",
								Arguments: []cli.Argument{
									&cli.StringArg{Name: "ADDR"},
								},
								Flags: []cli.Flag{
									&cli.BoolFlag{Name: "dry-run", Usage: "If set, prints out what would've been moved but doesn't actually move anything."},
									&cli.StringFlag{Name: "backup", Usage: "Path where Terraform should write the backup state."},
									&cli.BoolFlag{Name: "no-lock", Usage: "Don't hold a state lock during the operation. This is dangerous if others might concurrently run commands against the same workspace."},
									&cli.StringFlag{Name: "lock-timeout", Usage: "Duration to retry a state lock."},
									&cli.StringFlag{Name: "state", Usage: "Path to the state file to update. Defaults to the current workspace state."},
									&cli.BoolFlag{Name: "ignore-remote-version", Usage: "A rare option used for the remote backend only. See the remote backend documentation for more information."},
								},
								OnUsageError:             OnUsageErrorSite,
								InvalidFlagAccessHandler: InvalidFlagAccessHandler,
								Action: func(ctx context.Context, c *cli.Command) error {
									options := []string{}
									args := []string{}
									if c.StringArg("ADDR") == "" {
										return errors.New("missing argument 'ADDR'")
									}
									args = append(args, c.StringArg("ADDR"))
									if c.Bool("dry-run") {
										options = append(options, "-dry-run")
									}
									if c.String("backup") != "" {
										options = append(options, "-backup")
									}
									if c.Bool("no-lock") {
										options = append(options, "-lock=false")
									}
									if c.String("lock-timeout") != "" {
										options = append(options, "-lock-timeout="+c.String("lock-timeout"))
									}
									if c.String("state") != "" {
										options = append(options, "-state "+c.String("state"))
									}
									if c.Bool("ignore-remote-version") {
										options = append(options, "-ignore-remote-version")
									}
									return tf.CommandRunner("state rm", c.String("site"), args, options, "")
								},
							},
							{
								Name:      "show",
								Usage:     "Show a resource in the state",
								ArgsUsage: "ADDRESS",
								Arguments: []cli.Argument{
									&cli.StringArg{Name: "ADDR"},
								},
								Flags: []cli.Flag{
									&cli.StringFlag{Name: "state", Usage: "Path to the state file to update. Defaults to the current workspace state."},
								},
								OnUsageError:             OnUsageErrorSite,
								InvalidFlagAccessHandler: InvalidFlagAccessHandler,
								Action: func(ctx context.Context, c *cli.Command) error {
									options := []string{}
									args := []string{}
									if c.StringArg("ADDR") == "" {
										return errors.New("missing argument 'ADDR'")
									}
									args = append(args, c.StringArg("ADDR"))
									if c.String("state") != "" {
										options = append(options, "-state "+c.StringArg("state"))
									}
									return tf.CommandRunner("state show", c.String("site"), args, options, "")
								},
							},
						},
						CommandNotFound:          CommandNotFound,
						OnUsageError:             OnUsageErrorSite,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
					},
					{
						Name:      "taint",
						HideHelp:  true,
						Usage:     "Mark a resource instance as not fully functional",
						UsageText: appName + " tf taint [OPTIONS] ADDRESS",
						ArgsUsage: "ADDRESS",
						Arguments: []cli.Argument{
							&cli.StringArg{Name: "ADDR"},
						},
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "allow-missing", Usage: "If specified, the command will succeed (exit code 0) even if the resource is missing."},
							// Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-lock", Usage: "Don't hold a state lock during the operation. This is dangerous if others might concurrently run commands against the same workspace."},
							&cli.StringFlag{Name: "lock-timeout", Usage: "Duration to retry a state lock."},
							&cli.BoolFlag{Name: "ignore-remote-version", Usage: "A rare option used for the remote backend only. See the remote backend documentation for more information."},
						},
						CommandNotFound:          CommandNotFound,
						OnUsageError:             OnUsageErrorSite,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
						Action: func(ctx context.Context, c *cli.Command) error {
							options := []string{}
							if c.StringArg("ADDR") == "" {
								return errors.New("missing argument 'ADDR'")
							}
							args := append([]string{}, c.StringArg("ADDR"))
							if c.Bool("allow-missing") {
								options = append(options, "-allow-missing")
							}
							if c.Bool("no-lock") {
								options = append(options, "-lock=false")
							}
							if c.String("lock-timeout") != "" {
								options = append(options, "-lock-timeout="+c.String("lock-timeout"))
							}
							if c.Bool("ignore-remote-version") {
								options = append(options, "-ignore-remote-version")
							}

							return tf.CommandRunner("taint", c.String("site"), args, options, "")
						},
					},
					{
						Name:      "untaint",
						HideHelp:  true,
						Usage:     "Remove the 'tainted' state from a resource instance",
						UsageText: appName + " tf untaint [OPTIONS] ADDRESS",
						ArgsUsage: "ADDRESS",
						Arguments: []cli.Argument{
							&cli.StringArg{Name: "ADDR"},
						},
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "allow-missing", Usage: "If specified, the command will succeed (exit code 0) even if the resource is missing."},
							// Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-lock", Usage: "Don't hold a state lock during the operation. This is dangerous if others might concurrently run commands against the same workspace."},
							&cli.StringFlag{Name: "lock-timeout", Usage: "Duration to retry a state lock."},
							&cli.BoolFlag{Name: "ignore-remote-version", Usage: "A rare option used for the remote backend only. See the remote backend documentation for more information."},
						},
						CommandNotFound:          CommandNotFound,
						OnUsageError:             OnUsageErrorSite,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
						Action: func(ctx context.Context, c *cli.Command) error {
							options := []string{}
							if c.StringArg("ADDR") == "" {
								return errors.New("missing argument 'ADDR'")
							}
							args := append([]string{}, c.StringArg("ADDR"))
							if c.Bool("allow-missing") {
								options = append(options, "-allow-missing")
							}
							if c.Bool("no-lock") {
								options = append(options, "-lock=false")
							}
							if c.String("lock-timeout") != "" {
								options = append(options, "-lock-timeout="+c.String("lock-timeout"))
							}
							if c.Bool("ignore-remote-version") {
								options = append(options, "-ignore-remote-version")
							}

							return tf.CommandRunner("untaint", c.String("site"), args, options, "")
						},
					},
					{
						Name:      "validate",
						HideHelp:  true,
						Usage:     "Validate the configuration files",
						UsageText: appName + " tf validate [OPTIONS]",
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "no-color", Usage: "If specified, output won't contain any color."},
							&cli.BoolFlag{Name: "json", Usage: "Output the version information as a JSON object."},
						},
						CommandNotFound:          CommandNotFound,
						OnUsageError:             OnUsageErrorSite,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
						Action: func(ctx context.Context, c *cli.Command) error {
							options := []string{}
							if c.Bool("json") {
								options = append(options, "-json")
							}
							if c.Bool("no-color") {
								options = append(options, "-no-color")
							}
							return tf.CommandRunner("validate", c.String("site"), []string{}, options, "")
						},
					},
					{
						Name:      "version",
						HideHelp:  true,
						Usage:     "Show the current Terraform version",
						UsageText: appName + " tf version [OPTIONS]",
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "json", Usage: "Output the version information as a JSON object."},
						},
						CommandNotFound:          CommandNotFound,
						OnUsageError:             OnUsageErrorSite,
						InvalidFlagAccessHandler: InvalidFlagAccessHandler,
						Action: func(ctx context.Context, c *cli.Command) error {
							options := []string{}
							if c.Bool("json") {
								options = append(options, "-json")
							}
							return tf.CommandRunner("version", c.String("site"), []string{}, options, "")
						},
					},
				},
				CommandNotFound:          CommandNotFound,
				OnUsageError:             OnUsageErrorSite,
				InvalidFlagAccessHandler: InvalidFlagAccessHandler,
				Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {

					err := inception.InitNeeded(fs)
					if err != nil {
						return ctx, err
					}

					sites := settings.Conf.Strings("sites.ordered")
					site := c.String("site")
					if strings.Contains(site, "site_") {
						logger.Zap.Warn("There is no need to pass the 'site_' prefix.")
						if err := c.Set("site", strings.TrimPrefix(site, "site_")); err != nil {
							return ctx, err
						}
					}
					// Changing site again fo correct display of the errors
					site = c.String("site")

					if _, err := fs.Stat("site_" + site); os.IsNotExist(err) && slices.Contains(sites, c.String("site")) {
						return ctx, errors.New("the site " + site + " exists but is missing inside the config file")
					} else if site != "" && !slices.Contains(sites, c.String("site")) {
						return ctx, errors.New("the site " + site + " does not exist")
					}

					return ctx, nil
				},
			}},
	}

	err := cmd.Run(context.Background(), os.Args)

	return err
}
