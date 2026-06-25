// This file contains some functions that are used in the cli and the TextWrapper customized for the cli

package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/montblu/terrabutler/internal/logger"

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
		// Other types are ignored
	} else {
		typeFlag = ""
	}

	return typeFlag
}

// Function to add "-" to the names available to a flag
func addIndentFlag(names []string) []string {
	for i, n := range names {
		if n != "" {
			if len(n) == 1 {
				names[i] = "-" + n
			} else {
				names[i] = "--" + n
			}
		}
	}
	return names
}

// Default Functions for the Logger, of commands and flags not found.
func CommandNotFound(ctx context.Context, c *cli.Command, s string) {
	fmt.Println("Usage: " + c.UsageText)
	fmt.Println("Try '" + c.FullName() + " -h' for help.")
	logger.Zap.Error("No such command '" + s + "'.")
}

func OnUsageError(ctx context.Context, cmd *cli.Command, err error, isSubcommand bool) error {
	return nil
}

func InvalidFlagAccessHandler(ctx context.Context, c *cli.Command, s string) {
	fmt.Println("Usage: " + c.UsageText)
	fmt.Println("Try '" + c.FullName() + " -h' for help.")
	logger.Zap.Error("No such option: '" + s + "'.")
}

// Function for the Subcommands of tf, to show the required use of the flag -site
func OnUsageErrorSite(ctx context.Context, cmd *cli.Command, err error, isSubcommand bool) error {

	switch {
	case err.Error() == "flag needs an argument: -site":
		fmt.Println("Usage: " + cmd.UsageText)
		fmt.Println("Try '" + cmd.FullName() + " -h' for help.")
		return errors.New("option '-site' requires an argument")
	case err.Error() == "Required flag \"site\" not set":
		fmt.Println("Usage: " + cmd.UsageText)
		fmt.Println("Try '" + cmd.FullName() + " -h' for help.")
		return errors.New("missing option '-site'")
	case strings.Contains(err.Error(), "flag provided but not defined:"):
		return nil
	}
	return err

}

// The New HelpPrinter function with support of:
//
// Defining max Length for the text be wrapped
// Calculating the offset need fr flags...
func HelpPrinterNewFunctions(w io.Writer, templ string, data interface{}) {
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

// New template with the flags begin able to support Indentation and Wrapper
var RootCommandHelpTemplate = `NAME:
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

// New template for the Sub-SubCommands with the flags begin able to support Indentation and Wrapper
var CommandHelpTemplate = `NAME:
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

// New template for the SubCommands with the flags begin able to support Indentation and Wrapper
var SubcommandHelpTemplate = `NAME:
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
