// cmd/contact_add.go
package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/thomas-stecinski/crm_go_project/internal/contacts"
)

var (
	addID    int
	addName  string
	addEmail string
)

var contactAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Ajouter un contact",
	Long:  "La commande add permet de creer un utilisateur et de l'enregistrer dans la base de données \n Exemple : go run . contact add --name ''Alice'' --email ''alice@mail.com''",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := contacts.Contact{
			ID:    addID,
			Nom:   addName,
			Email: addEmail,
		}
		if !c.IsValid() {
			return errors.New("nom et email sont obligatoires (--name, --email)")
		}
		_, err := store.Add(c)
		if err != nil {
			return fmt.Errorf("ajout impossible: %w", err)
		}
		fmt.Printf("✅ Contact ajouté (ID=%s)\n", idToStr(c.ID))
		return nil
	},
}

func idToStr(id int) string {
	if id == 0 {
		return "auto"
	}
	return strconv.Itoa(id)
}

func init() {
	contactCmd.AddCommand(contactAddCmd)
	contactAddCmd.Flags().IntVar(&addID, "id", 0, "ID explicite (optionnel, auto si 0)")
	contactAddCmd.Flags().StringVar(&addName, "name", "", "Nom (requis)")
	contactAddCmd.Flags().StringVar(&addEmail, "email", "", "Email (requis)")
	_ = contactAddCmd.MarkFlagRequired("name")
	_ = contactAddCmd.MarkFlagRequired("email")
}
