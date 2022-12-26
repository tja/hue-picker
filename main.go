package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Define root command
var CmdRoot = &cobra.Command{
	Use:               "hue-picker",
	Long:              "Philips Hue Color Picker",
	Args:              cobra.NoArgs,
	Version:           "0.0.1",
	PersistentPreRunE: setup,
	Run:               runRoot,
}

// Initialize command options
func init() {
	// Logging
	CmdRoot.PersistentFlags().String("log-level", "info", "verbosity of logging output")
	CmdRoot.Flags().Bool("json", false, "change logging format to JSON")

	// API
	CmdRoot.Flags().String("host", ":80", "host API server should listen on")
}

// runRoot is called when the root command is used.
func runRoot(cmd *cobra.Command, args []string) {
	// Start HTTP server
	srv := &http.Server{
		Addr: viper.GetString("host"),
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Fatal("Failed to start server")
		}
	}()

	logrus.WithField("host", srv.Addr).Info("Server is listening")

	// Wait for user termination
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done

	// Stop server
	logrus.Info("Server shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.WithError(err).Fatal("Failed to gracefully shut down server")
	}
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
	viper.AddConfigPath("$HOME/.config/hue-picker")
	viper.AddConfigPath(".")

	viper.ReadInConfig() //nolint:errcheck

	// Logging
	log.SetOutput(io.Discard)

	lvl, err := logrus.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		return fmt.Errorf("parse log level: %w", err)
	}

	logrus.SetLevel(lvl)

	if viper.GetBool("json") {
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
