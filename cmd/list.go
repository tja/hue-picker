package cmd

import (
	"fmt"
	"strconv"

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

	// Get all lights
	lights, err := bridge.GetLights()
	if err != nil {
		logrus.WithError(err).Fatal("Unable to get lights")
	}

	// Get all groups
	groups, err := bridge.GetGroups()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get groups")
	}

	// Dump lights for each room
	for _, g := range groups {
		// Skip if group is no room
		if g.Type != "Room" {
			continue
		}

		// Skip if no lights in room
		if len(g.Lights) == 0 {
			continue
		}

		// Dump room lights
		fmt.Printf("🏡 %s\n", g.Name)

		for _, lid := range g.Lights {
			for _, l := range lights {
				// Bail on wrong IDs
				if strconv.Itoa(l.ID) != lid {
					continue
				}

				// Dump light
				fmt.Printf("   💡 [%s] %s (%s)\n", l.UniqueID, l.Name, l.ProductName)
			}
		}

		fmt.Println()
	}
}