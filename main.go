package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tja/hue-picker/cmd"
)

// CmdRoot defines the CLI root command.
var CmdRoot = &cobra.Command{
	Use:               "hue-picker",
	Long:              "Philips Hue Color Picker",
	Args:              cobra.NoArgs,
	Version:           "0.1.0",
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	PersistentPreRunE: setup,
}

// Initialize CLI options.
func init() {
	// Logging
	CmdRoot.PersistentFlags().String("log-level", "info", "verbosity of logging output")
	CmdRoot.PersistentFlags().Bool("log-json", false, "change logging format to JSON")

	// Register sub-commands
	CmdRoot.AddCommand(cmd.CmdRegister)
	CmdRoot.AddCommand(cmd.CmdList)
	CmdRoot.AddCommand(cmd.CmdServe)
}

// setup will set up configuration management and logging.
//
// Configuration options can be set via the command line, via a configuration file (in the current folder, at
// "/etc/hue-picker/config.yaml" or at "~/.config/hue-picker/config.yaml"), and via environment variables (all
// uppercase and prefixed with "HUE_PICKER_").
func setup(cmd *cobra.Command, args []string) error {
	// Connect all options to Viper
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return fmt.Errorf("bind command line flags: %w", err)
	}

	// Environment variables
	viper.SetEnvPrefix("HUE_PICKER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	// Configuration file
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/hue-picker")
	viper.AddConfigPath(os.Getenv("HOME") + "/.config/hue-picker")
	viper.AddConfigPath(".")

	// Configuration file
	if err := viper.ReadInConfig(); err != nil {
		// Don't fail if config not found
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			logrus.WithError(err).Fatal("Unable to read config file")
		}
	}

	// Logging
	log.SetOutput(io.Discard)

	lvl, err := logrus.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		return fmt.Errorf("parse log level: %w", err)
	}

	logrus.SetLevel(lvl)

	if viper.GetBool("log-json") {
		// Use JSON formatter
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	return nil
}

// main is the main entry point of the command.
func main() {
	if err := CmdRoot.Execute(); err != nil {
		logrus.WithError(err).Fatal("Unable to execute command")
	}
}
