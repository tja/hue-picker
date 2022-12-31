package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amimof/huego"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Define command
var CmdRegister = &cobra.Command{
	Use:   "register [flags]",
	Short: "...",
	Args:  cobra.NoArgs,
	Run:   runRegister,
}

// Initialize command options
func init() {
	// Disovery
	CmdRegister.Flags().Int("index", -1, "Index of Hue bridge to register for")
}

// runRegister is called when the "register" command is used.
func runRegister(cmd *cobra.Command, args []string) {
	// Discover all bridges in the network
	bridges, err := huego.DiscoverAll()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to discover Hue bridges")
	}

	// Check if the bridge can be identified properly
	if len(bridges) == 0 {
		logrus.Fatal("No Hue bridge discovered")
	}

	if len(bridges) > 1 {
		index := viper.GetInt("index")
		if index < 0 {
			logrus.WithField("bridges", bridges).Fatal("Multiple bridges discovered, but no valid --index given")
		}

		bridges = bridges[index:index]
	}

	// Ask user to press button on Hue bridge
	logrus.WithField("id", bridges[0].ID).Info("Hue bridge discovered")

	fmt.Println(`    ╭──────────────╮`)
	fmt.Println(`    │              │`)
	fmt.Println(`    │      ▁▁      │`)
	fmt.Println(`    │     ╱  ╲     │    Press button`)
	fmt.Println(`    │     ╲▁▁╱     │    on Hue bridge`)
	fmt.Println(`    │              │`)
	fmt.Println(`    │              │`)
	fmt.Println(`    ╰──────────────╯`)

	// Event loop
	ticker := time.NewTicker(500 * time.Millisecond)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

loop:
	for {
		select {
		case <-done:
			// User interruption
			logrus.Info("Interrupted")
			break loop

		case <-ticker.C:
			// Try creating a user account
			user, err := bridges[0].CreateUser("hue-picker")
			if err != nil {
				// Check if API error 101 ("link button not pressed")
				var apiError *huego.APIError
				if errors.As(err, &apiError) && (apiError.Type == 101) {
					continue loop
				}

				// Bail on other error
				logrus.WithError(err).Fatal("Failed to create user")
			}

			// Successfully created user
			fmt.Printf("Successfully registered:\n")
			fmt.Printf("   Host:   %s\n", bridges[0].Host)
			fmt.Printf("   Bridge: %s\n", bridges[0].ID)
			fmt.Printf("   User:   %s\n", user)

			break loop
		}
	}
}
