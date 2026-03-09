package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {

	//Criação da Version Flag
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{},
		Usage:   "Show the version and exit",
	}
	//Custom Version Flag
	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Fprintf(cmd.Root().Writer, "%s: %s\n", cmd.Root().Name, cmd.Root().Version)
	}
	//Alteração da Help flag defaulf
	cli.HelpFlag = &cli.BoolFlag{
		Name:    "help",
		Aliases: []string{"H", "h"},
		Usage:   "Show this message and exit",
	}

	//Comando Base
	//
	// TODO:
	// Error Handling with wrong flags
	// Version with Semantic Versioning
	// Logs (Using Prints for Debugging)
	// Fix being possible to use -version flag in all subcommands
	//
	cmd := &cli.Command{
		Name:      "terrabutler",
		Usage:     "The utility that helps keeping your IaC in one piece",
		UsageText: "terrabutler [OPTIONS] COMMAND [ARGS]...",
		Version:   "v1.1.2",
		//Hides Help Command to "Remove" HelpCommand, you need to hide it for each command
		HideHelpCommand:       true,
		EnableShellCompletion: true,
		Suggest:               true,
		Commands: []*cli.Command{
			// env Command
			//
			// What is Done:
			// Added all SubCommands
			// All Flags and Arguments of the SubCommands
			//
			// TODO:
			// Finished for now...
			{
				Name:            "env",
				Usage:           "Manage environments",
				UsageText:       "terrabutler env [OPTIONS] COMMAND [ARGS]...",
				HideHelpCommand: true,
				Commands: []*cli.Command{
					//Subcommands of Env
					{
						Name:      "delete",
						Aliases:   []string{""},
						Usage:     "Delete an environment",
						UsageText: "terrabutler env delete [OPTIONS] NAME",
						HideHelp:  true,
						ArgsUsage: "NAME",
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
						Name:        "select",
						Aliases:     []string{""},
						Usage:       "Select a environment",
						UsageText:   "terrabutler env select [OPTIONS] NAME",
						HideHelp:    true,
						HideVersion: true,
						ArgsUsage:   "NAME",
						Arguments:   []cli.Argument{&cli.StringArg{Name: "ENV"}},
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
			// All subcommands
			//
			//
			// TODO:
			// Add the missing arguments and flags of the subcommands
			//
			{
				Name:      "tf",
				Usage:     "Initialize the manager",
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
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "console",
						HideHelp:  true,
						Usage:     "Try Terraform expressions at an interactive command...",
						UsageText: "",
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "destroy",
						HideHelp:  true,
						Usage:     "Prepare your working directory for other commands",
						UsageText: "",
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "fmt",
						HideHelp:  true,
						Usage:     "Reformat your configuration in the standardstyle",
						UsageText: "",
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "force-unlock",
						HideHelp:  true,
						Usage:     "Release a stuck lock on the current workspace",
						UsageText: "",
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "generate-options",
						HideHelp:  true,
						Usage:     "Generate terraform options",
						UsageText: "",
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "import",
						HideHelp:  true,
						Usage:     "Associate existing infrastructure with a Terraform...",
						UsageText: "",
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "init",
						HideHelp:  true,
						Usage:     "Prepare your working directory for other commands",
						UsageText: "",
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "output",
						HideHelp:  true,
						Usage:     "Show output values from your root module",
						UsageText: "",
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "plan",
						HideHelp:  true,
						Usage:     "Show changes required by the current configuration",
						UsageText: "",
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "providers",
						HideHelp:  true,
						Usage:     "Show the providers required for this configuration",
						UsageText: "",
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "refresh",
						HideHelp:  true,
						Usage:     "Update the state to match remote systems",
						UsageText: "",
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "show",
						HideHelp:  true,
						Usage:     "Show the current state or a saved plan",
						UsageText: "",
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "state",
						HideHelp:  true,
						Usage:     "Advanced state management",
						UsageText: "",
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "taint",
						HideHelp:  true,
						Usage:     "Mark a resource instance as not fully functional",
						UsageText: "",
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "untaint",
						HideHelp:  true,
						Usage:     "Remove the 'tainted' state from a resource instance",
						UsageText: "",
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "validate",
						HideHelp:  true,
						Usage:     "Validate the configuration files",
						UsageText: "",
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
					{
						Name:      "version",
						HideHelp:  true,
						Usage:     "Show the current Terraform version",
						UsageText: "",
						Action: func(ctx context.Context, c *cli.Command) error {
							return nil
						},
					},
				},
			}},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
