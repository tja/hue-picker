package cmd

import (
	"github.com/amimof/huego"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Define command
var CmdList = &cobra.Command{
	Use:   "list [flags]",
	Short: "...",
	Args:  cobra.NoArgs,
	Run:   runList,
}

// Initialize command options
func init() {
	// Hue Bridge
	CmdList.Flags().String("bridge", "", "ID of the Hue bridge")
	CmdList.Flags().String("user", "", "ID of user registered to the Hue bridge")
}

// runList is called when the "list" command is used.
func runList(cmd *cobra.Command, args []string) {
	// Discover all bridges in the network
	bridges, err := huego.DiscoverAll()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to discover Hue bridges")
	}

	// Look for the proper bridge
	var bridge *huego.Bridge

	for _, b := range bridges {
		if b.ID == viper.GetString("bridge") {
			bridge = &b
			break
		}
	}

	if bridge == nil {
		logrus.WithField("bridge", viper.GetString("bridge")).Fatal("Unable to find requested Hue bridge")
	}

	bridge = bridge.Login(viper.GetString("user"))

	// Get all lights from the bridge
	lights, err := bridge.GetLights()
	if err != nil {
		logrus.WithError(err).Fatal("Unable to get lights")
	}

	// Dump all lights known
	for _, l := range lights {
		logrus.
			WithField("name", l.Name).
			WithField("product", l.ProductName).
			WithField("id", l.UniqueID).
			Info("💡")
	}
}
