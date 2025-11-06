// cmd/interactive.go
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thomas-stecinski/crm_go_project/internal/cli"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Lancer le menu interactif (console)",
	RunE: func(cmd *cobra.Command, args []string) error {
		app := &cli.App{Store: store}
		app.RunInteractive()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}