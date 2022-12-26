package cmd

import (
	"github.com/spf13/cobra"
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
}

// runList is called when the "list" command is used.
func runList(cmd *cobra.Command, args []string) {
}
