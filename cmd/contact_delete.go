// cmd/contact_delete.go
package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thomas-stecinski/crm_go_project/internal/contacts"
)

var delID int

var contactDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Supprimer un contact par ID",
	Long:  "La commande delete permet de supprimer un utilisateur de la base de données \n Exemple : go run . contact delete --id 1",
	RunE: func(cmd *cobra.Command, args []string) error {
		if delID <= 0 {
			return errors.New("--id est requis et doit être > 0")
		}
		if err := store.Delete(delID); err != nil {
			if errors.Is(err, contacts.ErrNotFound) {
				return fmt.Errorf("contact %d introuvable", delID)
			}
			return fmt.Errorf("suppression impossible: %w", err)
		}
		fmt.Printf("🗑️  Contact %d supprimé.\n", delID)
		return nil
	},
}

func init() {
	contactCmd.AddCommand(contactDeleteCmd)
	contactDeleteCmd.Flags().IntVar(&delID, "id", 0, "ID du contact (requis)")
	_ = contactDeleteCmd.MarkFlagRequired("id")
}