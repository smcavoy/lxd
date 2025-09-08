package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/canonical/lxd/shared/version"
)

type cmdGlobal struct {
	flagVersion bool
	flagHelp    bool

	flagLogVerbose bool
	flagLogDebug   bool

	flagExcludeInterfaces string
}

func main() {
	// agent command (main)
	agentCmd := cmdAgent{}
	app := agentCmd.Command()
	app.SilenceUsage = true
	app.CompletionOptions = cobra.CompletionOptions{DisableDefaultCmd: true}

	// Workaround for main command
	app.Args = cobra.ArbitraryArgs

	// Global flags
	globalCmd := cmdGlobal{}
	agentCmd.global = &globalCmd
	app.PersistentFlags().BoolVar(&globalCmd.flagVersion, "version", false, "Print version number")
	app.PersistentFlags().BoolVarP(&globalCmd.flagHelp, "help", "h", false, "Print help")
	app.PersistentFlags().BoolVarP(&globalCmd.flagLogVerbose, "verbose", "v", false, "Show all information messages")
	app.PersistentFlags().BoolVarP(&globalCmd.flagLogDebug, "debug", "d", false, "Show all debug messages")
	app.PersistentFlags().StringVar(&globalCmd.flagExcludeInterfaces, "exclude-interfaces", "", "Comma-separated list of network interfaces to exclude from stats/state queries")

	// Version handling
	app.SetVersionTemplate("{{.Version}}\n")
	app.Version = version.Version
	if version.IsLTSVersion {
		app.Version = version.Version + " LTS"
	}

	// Run the main command and handle errors
	err := app.Execute()
	if err != nil {
		// Ensure we exit with a non-zero exit code.
		os.Exit(1)
	}
}
