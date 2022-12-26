package cmd

import (
	"github.com/spf13/cobra"
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
}

// runRegister is called when the "register" command is used.
func runRegister(cmd *cobra.Command, args []string) {
}
