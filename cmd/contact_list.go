// cmd/contact_list.go
package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/thomas-stecinski/crm_go_project/internal/contacts"
)

var contactListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lister tous les contacts",
	Long:  "La commande list permet d'afficher tous les utilisateurs de la base de données \n Exemple : go run . contact list",
	RunE: func(cmd *cobra.Command, args []string) error {
		all := store.All()
		if len(all) == 0 {
			fmt.Println("Aucun contact.")
			return nil
		}
		// tri par ID pour affichage déterministe
		sort.Slice(all, func(i, j int) bool { return all[i].ID < all[j].ID })
		for _, c := range all {
			fmt.Printf("ID: %d | Nom: %s | Email: %s\n", c.ID, c.Nom, c.Email)
		}
		return nil
	},
}

func init() {
	contactCmd.AddCommand(contactListCmd)
	// pas de flags nécessaires
	_ = contacts.ErrNotFound // éviter l'import inutile si pas utilisé ici
}
