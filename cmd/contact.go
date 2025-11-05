// cmd/contact.go
package cmd

import "github.com/spf13/cobra"

var contactCmd = &cobra.Command{
	Use:   "contact",
	Short: "GÃ©rer les contacts (add, list, delete, update)",
}

func init() {
	rootCmd.AddCommand(contactCmd)
}