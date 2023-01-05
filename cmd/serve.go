package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amimof/huego"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tja/hue-picker/internal/api"
)

// Define command
var CmdServe = &cobra.Command{
	Use:   "serve [flags]",
	Short: "Serve web application via built-in web server",
	Args:  cobra.NoArgs,
	Run:   runServe,
}

// Initialize command options
func init() {
	// Hue
	CmdServe.Flags().String("host", "", "Host address of the Hue bridge (or empty)")
	CmdServe.Flags().String("bridge", "", "ID of the Hue bridge (or empty)")
	CmdServe.Flags().String("user", "", "ID of user registered to the Hue bridge")
	CmdServe.Flags().String("light", "", "ID of light registered to the Hue bridge")

	// Server
	CmdServe.Flags().String("listen", ":80", "host the server should listen on")
	CmdServe.Flags().String("www", "./etc/www/", "folder to static web assets")
}

// runServe is called when the "list" command is used.
func runServe(cmd *cobra.Command, args []string) {
	// Set up Hue bridge
	var bridge *huego.Bridge

	if host := viper.GetString("host"); host != "" {
		// Create Hue bridge directly
		bridge = &huego.Bridge{
			Host: host,
			User: viper.GetString("user"),
			ID:   viper.GetString("bridge"),
		}
	} else {
		// Otherwise, discover all bridges in the network
		bridges, err := huego.DiscoverAll()
		if err != nil {
			logrus.WithError(err).Fatal("Failed to discover Hue bridges")
		}

		// Look for the proper bridge
		for _, b := range bridges {
			if b.ID == viper.GetString("bridge") {
				bridge = &b
				break
			}
		}

		if bridge == nil {
			logrus.WithField("bridge", viper.GetString("bridge")).Fatal("Unable to find requested Hue bridge")
		}

		// Log in
		bridge = bridge.Login(viper.GetString("user"))
	}

	// Get all lights from the bridge
	lights, err := bridge.GetLights()
	if err != nil {
		logrus.WithError(err).Fatal("Unable to get lights")
	}

	// Look for the proper light
	var light *huego.Light

	for _, l := range lights {
		if l.UniqueID == viper.GetString("light") {
			light = &l
			break
		}
	}

	if light == nil {
		logrus.WithField("light", viper.GetString("light")).Fatal("Unable to find requested light")
	}

	// Create API server
	s, err := api.NewServer(bridge, light)
	if err != nil {
		logrus.WithError(err).Fatal("Unable to create API server")
	}

	// Start HTTP server
	m := chi.NewMux()

	m.Mount("/api", s.M)
	m.Handle("/*", http.FileServer(http.Dir(viper.GetString("www"))))

	srv := &http.Server{
		Addr:    viper.GetString("listen"),
		Handler: m,
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
