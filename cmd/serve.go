package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Define command
var CmdServe = &cobra.Command{
	Use:   "serve [flags]",
	Short: "...",
	Args:  cobra.NoArgs,
	Run:   runServe,
}

// Initialize command options
func init() {
	// Server
	CmdServe.Flags().String("host", ":80", "host the server should listen on")
}

// runServe is called when the "list" command is used.
func runServe(cmd *cobra.Command, args []string) {
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
