package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

// Function to Create offsets to the flags.
func offsetFlags(flag []cli.Flag, fixed int) int {
	max := 0
	for _, f := range flag {
		s := strings.Join(f.Names(), "--,    ")
		if len(s) > max {
			max = len(s)
		}
	}
	return max + fixed
}

// Function to get the type of Value a flag Requires
func typeFlag(flag cli.Flag) string {
	// Other Types aren't required in terraButler, so only Getting returning Strings Types
	var typeFlag string
	// String -> TEXT
	typeFlag = fmt.Sprintf("%T", flag.Get())
	if typeFlag == "string" {
		typeFlag = " TEXT"
		// Special Case for the flag -site of TerraButler
		if strings.Compare(flag.Names()[0], "site") == 0 {
			typeFlag = " SITE"
		}
		//Other types are ignored
	} else {
		typeFlag = ""
	}

	return typeFlag
}

// Function to add "-" to the names available to a flag
func addIndentFlag(names []string) []string {
	for i, n := range names {
		if n != "" {
			if len(names) == 1 {
				names[i] = "-" + n
			} else {
				names[i] = "--" + n
			}
		}
	}
	return names
}

func main() {

	//Changing the default help flag
	cli.HelpFlag = &cli.BoolFlag{
		Name:    "help",
		Aliases: []string{"H", "h"},
		Usage:   "Show this message and exit",
	}

	// New HelpPrinter function with support of:
	//
	// Defining max Length for the text be wrapped
	// Calculating the offset need fr flags...
	cli.HelpPrinter = func(w io.Writer, templ string, data interface{}) {
		funcMap := map[string]interface{}{
			"wrapAt": func() int {
				return 80
			},
			"offsetFlags":   offsetFlags,
			"addIndentFlag": addIndentFlag,
			"typeFlag":      typeFlag,
		}

		cli.HelpPrinterCustom(w, templ, data, funcMap)
	}

	//In the future allocate this templates to another file...
	//This modifications to the template only affect "Visible" flags

	//New template with the flags begin able to support Indentation and Wrapper
	cli.RootCommandHelpTemplate = `NAME:
   {{template "helpNameTemplate" .}}

USAGE:
   {{if .UsageText}}{{wrap .UsageText 3}}{{else}}{{.FullName}} {{if .VisibleFlags}}[global options]{{end}}{{if .VisibleCommands}} [command [command options]]{{end}}{{if .ArgsUsage}} {{.ArgsUsage}}{{else}}{{if .Arguments}} [arguments...]{{end}}{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

VERSION:
   {{.Version}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
   {{template "descriptionTemplate" .}}{{end}}
{{- if len .Authors}}

AUTHOR{{template "authorsTemplate" .}}{{end}}{{if .VisibleCommands}}

COMMANDS:{{template "visibleCommandCategoryTemplate" .}}{{end}}{{if .VisibleFlagCategories}}

GLOBAL OPTIONS:{{template "visibleFlagCategoryTemplate" .}}{{else if .VisibleFlags}}

GLOBAL OPTIONS:{{ $cv := offsetFlags .VisibleFlags 0}}{{range $i, $e := .VisibleFlags}}
	{{$f := addIndentFlag $e.Names}}{{$s := join $f ", "}}{{$t := typeFlag $e}}{{$s}}{{$t}}{{ $sp := subtract $cv (offset $s 3) }}{{ $sp = subtract $sp (offset $t -4)}}{{ indent $sp ""}}{{wrap $e.Usage 6}}{{end}}{{end}}{{if .Copyright}}

COPYRIGHT:
   {{template "copyrightTemplate" .}}{{end}}
`
	//New template for the Sub-SubCommands with the flags begin able to support Indentation and Wrapper
	cli.CommandHelpTemplate = `NAME:
   {{template "helpNameTemplate" .}}

USAGE:
   {{template "usageTemplate" .}}{{if .Category}}

CATEGORY:
   {{.Category}}{{end}}{{if .Description}}

DESCRIPTION:
   {{template "descriptionTemplate" .}}{{end}}{{if .VisibleFlagCategories}}

OPTIONS:{{template "visibleFlagCategoryTemplate" .}}{{else if .VisibleFlags}}

OPTIONS:
	{{ $cv := offsetFlags .VisibleFlags 8}}{{range $i, $e := .VisibleFlags}}
	{{$f := addIndentFlag $e.Names}}{{$s := join $f ", "}}{{$t := typeFlag $e}}{{$s}}{{$t}}{{ $sp := subtract $cv (offset $s 3) }}{{ $sp = subtract $sp (offset $t -4)}}{{ indent $sp ""}}{{wrap $e.Usage 6}}{{end}}{{end}}{{if .VisiblePersistentFlags}}

GLOBAL OPTIONS:{{template "visiblePersistentFlagTemplate" .}}{{end}}
`
	//New template for the SubCommands with the flags begin able to support Indentation and Wrapper
	cli.SubcommandHelpTemplate = `NAME:
   {{template "helpNameTemplate" .}}

USAGE:
   {{if .UsageText}}{{wrap .UsageText 3}}{{else}}{{.FullName}}{{if .VisibleCommands}} [command [command options]]{{end}}{{if .ArgsUsage}} {{.ArgsUsage}}{{else}}{{if .Arguments}} [arguments...]{{end}}{{end}}{{end}}{{if .Category}}

CATEGORY:
   {{.Category}}{{end}}{{if .Description}}

DESCRIPTION:
   {{template "descriptionTemplate" .}}{{end}}{{if .VisibleCommands}}

COMMANDS:{{template "visibleCommandTemplate" .}}{{end}}{{if .VisibleFlagCategories}}

OPTIONS:{{template "visibleFlagCategoryTemplate" .}}{{else if .VisibleFlags}}

OPTIONS:{{ $cv := offsetFlags .VisibleFlags 8}}{{range $i, $e := .VisibleFlags}}
	{{$f := addIndentFlag $e.Names}}{{$s := join $f ", "}}{{$t := typeFlag $e}}{{$s}}{{$t}}{{ $sp := subtract $cv (offset $s 3) }}{{ $sp = subtract $sp (offset $t -4)}}{{ indent $sp ""}}{{wrap $e.Usage 6}}{{end}}{{end}}{{if .VisiblePersistentFlags}}

GLOBAL OPTIONS:{{template "visiblePersistentFlagTemplate" .}}{{end}}
`

	//CLI
	//
	// TODO:
	// Error Handling with wrong flags --> With Logger
	// Version with Semantic Versioning --> Not Urgent
	// Logs (Using Prints for Debugging) --> Next Step
	//
	// After CLI, start Configuration File (settings.py)

	cmd := &cli.Command{
		Name:      "terrabutler",
		Usage:     "The utility that helps keeping your IaC in one piece",
		UsageText: "terrabutler [OPTIONS] COMMAND [ARGS]...",
		Version:   "v3.0.0",
		//Hides Help Command to "Remove" HelpCommand, you need to hide it for each command
		HideHelpCommand:       true,
		HideVersion:           true,
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "Show version and exit",
				Action: func(ctx context.Context, c *cli.Command) error {
					fmt.Fprintf(c.Root().Writer, "%s: %s\n", c.Root().Name, c.Root().Version)
					return nil
				},
			},
			// env Command
			//
			// What is Done:
			// Added all SubCommands
			// All Flags and Arguments of the SubCommands
			//
			// TODO:
			// Finished for now...
			{
				Name:      "env",
				Usage:     "Manage environments",
				UsageText: "terrabutler env [OPTIONS] COMMAND [ARGS]...",
				HideHelp:  true,
				Commands: []*cli.Command{
					//Subcommands of Env
					{
						Name:      "delete",
						Aliases:   []string{""},
						Usage:     "Delete an environment",
						UsageText: "terrabutler env delete [OPTIONS] NAME",
						ArgsUsage: "NAME",
						HideHelp:  true,
						Arguments: []cli.Argument{&cli.StringArg{Name: "ENV"}},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								//Added destroy as alias to the -d flag
								Name:    "destroy",
								Aliases: []string{"d"},
								Usage:   "Destroy all sites by inverse order.",
							},
							&cli.BoolFlag{
								Name:    "y",
								Aliases: []string{""},
								Usage:   "Delete without asking for confirmation.",
							},
							&cli.BoolFlag{
								//Flags with more than 1 letter are shown with double dash, but it is still accepted with -
								Name:    "s3",
								Aliases: []string{"S3"},
								Usage:   "Access S3 instead of parsing terraform output.",
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							//Test Ouput
							fmt.Println("Deleted Environment", cmd.StringArg("ENV"), "\nActive Flags: \n-d", cmd.Bool("destroy"), "\n-y", cmd.Bool("y"), "\n-s3", cmd.Bool("s3"))
							return nil
						}},

					{
						Name:     "list",
						Aliases:  []string{""},
						Usage:    "List environments",
						HideHelp: true,
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "s3",
								Aliases: []string{"S3"},
								Usage:   "Access S3 instead of parsing terraform output.",
							}},
						Action: func(context.Context, *cli.Command) error {
							//Test Ouput
							fmt.Println("Listed Environments")
							return nil
						}},
					{
						Name:      "new",
						Aliases:   []string{""},
						Usage:     "Create a new environment",
						UsageText: "terrabutler env new [OPTIONS] NAME",
						HideHelp:  true,
						ArgsUsage: "NAME",
						Arguments: []cli.Argument{&cli.StringArg{Name: "ENV"}},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "y",
								Aliases: []string{""},
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
							&cli.BoolFlag{
								Name:    "s3",
								Aliases: []string{"S3"},
								Usage:   "Access S3 instead of parsing terraform output.",
							},
						},
						Action: func(context.Context, *cli.Command) error {
							//Test Ouput
							fmt.Println("Created Environment")
							return nil
						}},
					{
						Name:      "reload",
						Aliases:   []string{""},
						HideHelp:  true,
						Usage:     "Reload the current environment",
						UsageText: "terrabutler env reload [OPTIONS]",
						Action: func(context.Context, *cli.Command) error {
							//Test Ouput
							fmt.Println("Reloaded Environment")
							return nil
						}},
					{
						Name:      "select",
						Aliases:   []string{""},
						Usage:     "Select a environment",
						UsageText: "terrabutler env select [OPTIONS] NAME",
						HideHelp:  true,
						ArgsUsage: "NAME",
						Arguments: []cli.Argument{&cli.StringArg{Name: "ENV"}},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "init",
								Aliases: []string{},
								Usage:   "Disable auto init of the sites.",
							},
							&cli.BoolFlag{
								Name:    "s3",
								Aliases: []string{"S3"},
								Usage:   "Access S3 instead of parsing terraform output.",
							},
						},
						Action: func(context.Context, *cli.Command) error {
							//Test Ouput
							fmt.Println("Selected Environment")
							return nil
						}},
					{
						Name:      "show",
						Aliases:   []string{""},
						HideHelp:  true,
						Usage:     "Show the name of the current environment",
						UsageText: "terrabutler env show [OPTIONS]",
						Action: func(context.Context, *cli.Command) error {
							//Test Ouput
							fmt.Println("Current Environment is ...")
							return nil
						}},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return nil
				},
			},
			// init Command
			//
			// TODO:
			// Concluded for now
			{
				Name:      "init",
				Usage:     "Initialize the manager",
				UsageText: "terrabutler init [OPTIONS]",
				HideHelp:  true,
				Action: func(ctx context.Context, cmd *cli.Command) error {
					//Test Ouput
					fmt.Println("The initialization was successfull!")
					return nil
				},
			},
			// tf Command
			//
			// What is done:
			// The required flag -site
			// All subcommands, flags and arguments..
			//
			//
			// TODO:
			// Flag site Warnings
			//
			//
			{
				Name:      "tf",
				Usage:     "Manage terraform commands",
				UsageText: "terrabutler tf [OPTIONS] COMMAND [ARGS]...",
				HideHelp:  true,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "site", Required: true, Usage: "Site where to run terraform.  [required]"}},
				Commands: []*cli.Command{
					{
						Name:      "apply",
						HideHelp:  true,
						Usage:     "Create or update infrastructure.",
						UsageText: "",
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "auto-approve", Usage: "Skip interactive approval of plan before applying."},
							&cli.BoolFlag{Name: "destroy", Usage: "Select the 'destroy' planning mode, which creates a plan to destroy all objects currently managed by this Terraform configuration instead of the usual behavior."},
							// Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-input", Usage: "Don't ask for input for variables if not directly set."},
							// Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-lock", Usage: `Don't hold a state lock during backend migration. This is dangerous if others might concurrently run commands against the same workspace.`},
							&cli.StringFlag{Name: "lock-timeout", Usage: "Duration to retry a state lock."},
							&cli.BoolFlag{Name: "no-color", Usage: "If specified, output won't contain any color."},
							&cli.BoolFlag{Name: "refresh-only", Usage: "Select the 'refresh only' planning mode, which checks whether remote objects still match the outcome of the most recent Terraform apply but does not propose any actions to undo any changes made outside of Terraform."},
							//Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-refresh", Usage: "Skip checking for external changes to remote objects while creating the plan. This can potentially make planning faster, but at the expense of possibly planning against a stale record of the remote system state."},
							&cli.StringSliceFlag{Name: "target", Usage: "Limit the planning operation to only the given module, resource, or resource instance and all of its dependencies. You can use this option multiple times to include more than one object. This is for exceptional use only."},
							&cli.BoolFlag{Name: "var", Usage: "Set a value for one of the input variables in the root module of the configuration. Use this option more than once to set more than one variable."},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							// To verify if its possible multiple "Targets" and if i get its values.
							fmt.Println(c.StringSlice("target"))
							return nil
						},
					},
					{
						Name:      "console",
						HideHelp:  true,
						Usage:     "Try Terraform expressions at an interactive command...",
						UsageText: "",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "state", Usage: "Legacy option for the local backend only. See the local backend's documentation for more information."},
							&cli.BoolFlag{Name: "plan", Usage: "Create a new plan (as if running \"terraform plan\") and then evaluate expressions against its planned state, instead of evaluating against the current state. You can use this to inspect the effects of configuration changes that haven't been applied yet.."},
							&cli.StringFlag{Name: "var", Usage: "Set a variable in the Terraform configuration. This flag can be set multiple times."},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "destroy",
						HideHelp:  true,
						Usage:     "Prepare your working directory for other commands",
						UsageText: "",
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "auto-approve", Usage: "Skip interactive approval of plan before applying."},
							// Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-input", Usage: "Don't ask for input for variables if not directly set."},
							// Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-lock", Usage: "Don't hold a state lock during backend migration. This is dangerous if others might concurrently run commands against the same workspace."},
							&cli.StringFlag{Name: "lock-timeout", Usage: "Duration to retry a state lock."},
							&cli.BoolFlag{Name: "no-color", Usage: "If specified, output won't contain any color."},
							&cli.BoolFlag{Name: "refresh-only", Usage: "Select the 'refresh only' planning mode, which checks whether remote objects still match the outcome of the most recent Terraform apply but does not propose any actions to undo any changes made outside of Terraform."},
							//Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-refresh", Usage: " Skip checking for external changes to remote objects while creating the plan. This can potentially make planning faster, but at the expense of possibly planning against a stale record of the remote system state."},
							&cli.StringFlag{Name: "target", Usage: "Limit the planning operation to only the given module, resource, or resource instance and all of its dependencies. You can use this option multiple times to include more than one object. This is for exceptional use only."},
							&cli.StringFlag{Name: "var", Usage: "Set a value for one of the input variables in the root module of the configuration. Use this option more than once to set more than one variable."},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "fmt",
						HideHelp:  true,
						Usage:     "Reformat your configuration in the standardstyle",
						UsageText: "",
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "diff", Usage: "Display diffs of formatting changes."},
							&cli.BoolFlag{Name: "no-color", Usage: "If specified, output won't contain any color."},
							&cli.BoolFlag{Name: "recursive", Usage: "Also process files in subdirectories. By default, only the given directory (or current directory) is processed."},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "force-unlock",
						HideHelp:  true,
						Usage:     "Release a stuck lock on the current workspace",
						UsageText: "",
						ArgsUsage: "LOCK-ID",
						Arguments: []cli.Argument{&cli.Int16Arg{Name: "LOCK-ID"}},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "force",
								Usage: "Don't ask for input for unlock confirmation.",
							},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "generate-options",
						HideHelp:  true,
						Usage:     "Generate terraform options",
						UsageText: "",
						ArgsUsage: "{init|plan|apply}",
						Arguments: []cli.Argument{
							&cli.StringArg{Name: "Choice"},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "import",
						HideHelp:  true,
						Usage:     "Associate existing infrastructure with a Terraform...",
						UsageText: "",
						ArgsUsage: "ADDR ID",
						Arguments: []cli.Argument{
							&cli.StringArgs{Min: 1, Max: 1, Name: "ADDR"},
							&cli.Int16Args{Min: 1, Max: 1, Name: "ID"},
						},
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "allow-missing-config", Usage: "Allow import when no resource configuration block exists."},
							//Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-input", Usage: "Don't ask for input for variables if not directly set."},
							//Requires BOOLEAN value -- Reversing
							&cli.BoolFlag{Name: "no-lock", Usage: "Don't hold a state lock during the operation. This is dangerous if others might concurrently run commands against the same workspace."},
							&cli.BoolFlag{Name: "no-color", Usage: "If specified, output won't contain any color."},
							&cli.StringSliceFlag{Name: "var", Usage: "Set a variable in the Terraform configuration. This flag can be set multiple times."},
							&cli.StringFlag{Name: "ignore-remote-version", Usage: "A rare option used for the remote backend only. See the remote backend documentation for more information."},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "init",
						HideHelp:  true,
						Usage:     "Prepare your working directory for other commands",
						UsageText: "",
						Flags: []cli.Flag{
							//Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-backend", Usage: "Disable backend or Terraform Cloud initialization for this configuration and use what what was previously initialized instead."},
							&cli.BoolFlag{Name: "force-copy", Usage: "Allow import when no resource configuration block exists."},
							//Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-get", Usage: "Disable downloading modules for this configuration."},
							//Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-input", Usage: "Disable interactive prompts. Note that some actions may require interactive prompts and will error if input is disabled."},
							//Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-lock", Usage: "Don't hold a state lock during backend migration. This is dangerous if others might concurrently run commands against the same workspace."},
							&cli.BoolFlag{Name: "no-color", Usage: "If specified, output won't contain any color."},
							&cli.BoolFlag{Name: "reconfigure", Usage: "Reconfigure a backend, ignoring any saved configuration."},
							&cli.BoolFlag{Name: "migrate-state", Usage: "Reconfigure a backend, and attempt to migrate any existing state."},
							&cli.BoolFlag{Name: "upgrade", Usage: "Install the latest module and provider versions allowed within configured constraints, overriding the default behavior of selecting exactly the version recorded in the dependency lockfile."},
							&cli.StringFlag{Name: "lockfile", Usage: "Set a dependency lockfile mode. Currently only 'readonly' is valid."},
							&cli.BoolFlag{Name: "ignore-remote-version", Usage: "A rare option used for Terraform Cloud and the remote backend only. Set this to ignore checking that the local and remote Terraform versions use compatible state representations, making an operation proceed even when there is a potential mismatch. See the documentation on configuring Terraform with Terraform Cloud for more information."},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "output",
						HideHelp:  true,
						Usage:     "Show output values from your root module",
						UsageText: "",
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "no-color", Usage: "If specified, output won't contain any color."},
							&cli.BoolFlag{Name: "json", Usage: "If specified, machine readable output will be printed in JSON format."},
							&cli.BoolFlag{Name: "raw", Usage: "For value types that can be automatically converted to a string, will print the raw string directly, rather than a human-oriented representation of the value."},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "plan",
						HideHelp:  true,
						Usage:     "Show changes required by the current configuration",
						UsageText: "",
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "destroy", Usage: "Select the 'destroy' planning mode, which creates a plan to destroy all objects currently managed by this"},
							//Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-input", Usage: "Don't ask for input for variables if not directly set."},
							//Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-lock", Usage: "Don't hold a state lock during backend migration. This is dangerous if others might concurrently run commands against the same workspace."},
							&cli.StringFlag{Name: "lock-timeout", Usage: "Duration to retry a state lock."},
							&cli.BoolFlag{Name: "no-color", Usage: "If specified, output won't contain any color."},
							&cli.BoolFlag{Name: "refresh-only", Usage: "Select the 'refresh only' planning mode, which checks whether remote objects still match the outcome of the most recent Terraform apply but does not propose any actions to undo any changes made outside of Terraform."},
							//Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-refresh", Usage: "Skip checking for external changes to remote objects while creating the plan. This can potentially make planning faster, but at the expense of possibly planning against a stale record of the remote system state."},
							&cli.StringSliceFlag{Name: "target", Usage: "Limit the planning operation to only the given module, resource, or resource instance and all of its dependencies. You can use this option multiple times to include more than one object. This is for exceptional use only."},
							&cli.StringFlag{Name: "var", Usage: "Set a value for one of the input variables in the root module of the configuration. Use this option more than once to set more than one variable."},
							&cli.StringFlag{Name: "out", Usage: "Write a plan file to the given path. This can be used as input to the \"apply\" command."},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "providers",
						HideHelp:  true,
						Usage:     "Show the providers required for this configuration",
						UsageText: "",
						Commands: []*cli.Command{
							//Makeup Providers...
							{
								Name:      "lock",
								Usage:     "Write out dependency locks for the configured providers",
								ArgsUsage: "PROVIDERS...",
								Arguments: []cli.Argument{
									&cli.StringArgs{Name: "Providers"},
								},
								Flags: []cli.Flag{
									&cli.StringFlag{Name: "fs-mirror", Usage: "Consult the given filesystem mirror directory instead of the origin registry for each of the given providers."},
									&cli.StringFlag{Name: "net-mirror", Usage: "Consult the given network mirror (given as a base URL) instead of the origin registry for each of the given providers."},
									&cli.StringFlag{Name: "platform", Usage: "Choose a target platform to request package checksums for."},
								},
								HideHelp: true,
							}, //Makeup DIRS..
							{
								Name:      "mirror",
								Usage:     "Save local copies of all required provider plugins",
								ArgsUsage: "TARGET_DIR",
								Arguments: []cli.Argument{
									&cli.StringArgs{Name: "DIR"},
								},
								Flags: []cli.Flag{
									&cli.StringFlag{Name: "platform", Usage: "Choose a target platform to request package checksums for."},
								},
								HideHelp: true,
							},
							{
								Name:  "schema",
								Usage: "Show schemas for the providers used in the configuration",
								Flags: []cli.Flag{
									&cli.BoolFlag{Name: "json", Required: true, Usage: "Prints out a json representation of the schemas for all providers used in the current configuration.  [required]"},
								},
								HideHelp: true,
							},
						},
					},
					{
						Name:      "refresh",
						HideHelp:  true,
						Usage:     "Update the state to match remote systems",
						UsageText: "",
						Flags: []cli.Flag{
							//Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-input", Usage: "Don't ask for input for variables if not directly set."},
							//Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "lock", Usage: "Don't hold a state lock during the operation. This is dangerous if others might concurrently run commands against the same workspace."},
							&cli.BoolFlag{Name: "no-color", Usage: "If specified, output won't contain any color."},
							&cli.StringFlag{Name: "target", Usage: "Resource to target. Operation will be limited to this resource and its dependencies. This flag can be used multiple times."},
							&cli.StringSliceFlag{Name: "var", Usage: "Set a variable in the Terraform configuration. This flag can be set multiple times."},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "show",
						HideHelp:  true,
						Usage:     "Show the current state or a saved plan",
						UsageText: "",
						ArgsUsage: "[PATH]",
						Arguments: []cli.Argument{
							&cli.StringArg{Name: "PATH"},
						},
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "no-color", Usage: "If specified, output won't contain any color."},
							&cli.BoolFlag{Name: "json", Usage: " Output the version information as a JSON object."},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "state",
						HideHelp:  true,
						Usage:     "Advanced state management",
						UsageText: "",
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
								Action: func(ctx context.Context, c *cli.Command) error {
									return nil
								},
							},
							{
								Name:      "mv",
								Usage:     "Move an item in the state",
								ArgsUsage: "SOURCE DESTINATION",
								Arguments: []cli.Argument{
									&cli.StringArgs{Name: "SOURCE"},
									&cli.StringArg{Name: "DESTINATION"},
								},
								Flags: []cli.Flag{
									&cli.BoolFlag{Name: "dry-run", Usage: "If set, prints out what would've been moved but doesn't actually move anything."},
									&cli.BoolFlag{Name: "lock", Usage: "Don't hold a state lock during the operation. This is dangerous if others might concurrently run commands against the same workspace."},
									&cli.StringFlag{Name: "lock-timeout", Usage: "Duration to retry a state lock."},
									&cli.BoolFlag{Name: "ignore-remote-version", Usage: "A rare option used for the remote backend only. See the remote backend documentation for more information."},
								},
							},
							{
								Name:  "pull",
								Usage: "Pull current state and output to stdouts",
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
									&cli.BoolFlag{Name: "lock", Usage: "Don't hold a state lock during the operation. This is dangerous if others might concurrently run commands against the same workspace."},
									&cli.StringFlag{Name: "lock-timeout", Usage: "Duration to retry a state lock."},
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
									&cli.BoolFlag{Name: "lock", Usage: "Don't hold a state lock during the operation. This is dangerous if others might concurrently run commands against the same workspace."},
									&cli.StringFlag{Name: "lock-timeout", Usage: "Duration to retry a state lock."},
									&cli.BoolFlag{Name: "ignore-remote-version", Usage: "A rare option used for the remote backend only. See the remote backend documentation for more information."},
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
									&cli.BoolFlag{Name: "lock", Usage: "Don't hold a state lock during the operation. This is dangerous if others might concurrently run commands against the same workspace."},
									&cli.StringFlag{Name: "lock-timeout", Usage: "Duration to retry a state lock."},
									&cli.StringFlag{Name: "state", Usage: "Path to the state file to update. Defaults to the current workspace state."},
									&cli.BoolFlag{Name: "ignore-remote-version", Usage: "A rare option used for the remote backend only. See the remote backend documentation for more information."},
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
							},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "taint",
						HideHelp:  true,
						Usage:     "Mark a resource instance as not fully functional",
						UsageText: "",
						ArgsUsage: "ADDRESS",
						Arguments: []cli.Argument{
							&cli.StringArg{Name: "ADDR"},
						},
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "allow-missing", Usage: " If specified, the command will succeed (exit code 0) even if the resource is missing."},
							//Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-lock", Usage: " Don't hold a state lock during the operation. This is dangerous if others might concurrently run commands against the same workspace."},
							&cli.StringFlag{Name: "lock-timeout", Usage: "Duration to retry a state lock."},
							&cli.BoolFlag{Name: "ignore-remote-version", Usage: "A rare option used for the remote backend only. See the remote backend documentation for more information."},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "untaint",
						HideHelp:  true,
						Usage:     "Remove the 'tainted' state from a resource instance",
						UsageText: "",
						ArgsUsage: "ADDRESS",
						Arguments: []cli.Argument{
							&cli.StringArg{Name: "ADDR"},
						},
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "allow-missing", Usage: " If specified, the command will succeed (exit code 0) even if the resource is missing."},
							//Requires BOOLEAN value --> Reversing
							&cli.BoolFlag{Name: "no-lock", Usage: " Don't hold a state lock during the operation. This is dangerous if others might concurrently run commands against the same workspace."},
							&cli.StringFlag{Name: "lock-timeout", Usage: "Duration to retry a state lock."},
							&cli.BoolFlag{Name: "ignore-remote-version", Usage: "A rare option used for the remote backend only. See the remote backend documentation for more information."},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "validate",
						HideHelp:  true,
						Usage:     "Validate the configuration files",
						UsageText: "",
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "no-color", Usage: "If specified, output won't contain any color."},
							&cli.BoolFlag{Name: "json", Usage: " Output the version information as a JSON object."},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "version",
						HideHelp:  true,
						Usage:     "Show the current Terraform version",
						UsageText: "",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "json", Usage: "Output the version information as a JSON object."},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					fmt.Println("Missing Command.")
					return nil
				},
			}},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
