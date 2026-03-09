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
		Aliases: []string{"V", "v"},
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
			// Rearrange the UsageText of the SubCommands
			{
				Name:            "env",
				Usage:           "Manage environments",
				UsageText:       "terrabutler env [OPTIONS] COMMAND [ARGS]...",
				HideHelpCommand: true,
				Commands: []*cli.Command{
					//Subcommands of Env
					{
						Name:            "delete",
						Aliases:         []string{""},
						Usage:           "Delete an environment",
						HideHelpCommand: true,
						ArgsUsage:       "NAME",
						Arguments:       []cli.Argument{&cli.StringArg{Name: "ENV"}},
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
						Name:    "list",
						Aliases: []string{""},
						Usage:   "List environments",
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
						Name:            "new",
						Aliases:         []string{""},
						Usage:           "Create a new environment",
						HideHelpCommand: true,
						ArgsUsage:       "NAME",
						Arguments:       []cli.Argument{&cli.StringArg{Name: "ENV"}},
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
						Name:    "reload",
						Aliases: []string{""},
						Usage:   "Reload the current environment",
						Action: func(context.Context, *cli.Command) error {
							//Test Ouput
							fmt.Println("Reloaded Environment")
							return nil
						}},
					{
						Name:            "select",
						Aliases:         []string{""},
						Usage:           "Select a environment",
						HideHelpCommand: true,
						ArgsUsage:       "NAME",
						Arguments:       []cli.Argument{&cli.StringArg{Name: "ENV"}},
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
							fmt.Println("Selected Enviroment")
							return nil
						}},
					{
						Name:    "show",
						Aliases: []string{""},
						Usage:   "Show the name of the current environment",
						Action: func(context.Context, *cli.Command) error {
							//Test Ouput
							fmt.Println("Current Enviroment is ...")
							return nil
						}},
				},
			},
			// init Command
			//
			// Concluded for now
			{
				Name:            "init",
				Usage:           "Initialize the manager",
				UsageText:       "terrabutler init [OPTIONS]",
				HideHelpCommand: true,
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
			// (With String flags it is possible to add "help" sites, without hiding Help Command)
			//
			// TODO:
			// Add all subcommands
			//
			{
				Name:            "tf",
				Usage:           "Initialize the manager",
				UsageText:       "terrabutler tf [OPTIONS] COMMAND [ARGS]...",
				HideHelpCommand: true,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "site", Required: true, Usage: "Site where to runterraform.  [required]"}},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					//Test Ouput
					fmt.Println("TerraForm Start")
					return nil
				},
			}},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
