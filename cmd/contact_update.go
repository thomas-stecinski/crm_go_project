// cmd/contact_update.go
package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thomas-stecinski/crm_go_project/internal/contacts"
)

var (
	updID    int
	updName  string
	updEmail string
)

var contactUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Mettre à jour un contact (nom/email) par ID",
	Long:  "La commande update permet de modifier un utilisateur existant dans la base de données \n Exemple : go run . contact update --id 1 --name ''Alice'' --email ''alice@mail.com''",
	RunE: func(cmd *cobra.Command, args []string) error {
		if updID <= 0 {
			return errors.New("--id est requis et doit être > 0")
		}
		existing, err := store.Get(updID)
		if err != nil {
			if errors.Is(err, contacts.ErrNotFound) {
				return fmt.Errorf("contact %d introuvable", updID)
			}
			return fmt.Errorf("récupération contact: %w", err)
		}

		if updName != "" {
			existing.Nom = updName
		}
		if updEmail != "" {
			existing.Email = updEmail
		}
		if !existing.IsValid() {
			return errors.New("nom et email ne peuvent pas être vides")
		}
		if err := store.Update(updID, existing); err != nil {
			return fmt.Errorf("mise à jour impossible: %w", err)
		}
		fmt.Printf("Contact %d mis à jour.\n", updID)
		return nil
	},
}

func init() {
	contactCmd.AddCommand(contactUpdateCmd)
	contactUpdateCmd.Flags().IntVar(&updID, "id", 0, "ID (requis)")
	contactUpdateCmd.Flags().StringVar(&updName, "name", "", "Nouveau nom (optionnel)")
	contactUpdateCmd.Flags().StringVar(&updEmail, "email", "", "Nouvel email (optionnel)")
	_ = contactUpdateCmd.MarkFlagRequired("id")
}