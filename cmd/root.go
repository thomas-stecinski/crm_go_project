// cmd/root.go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/thomas-stecinski/crm_go_project/internal/contacts"
)

var (
	dataPath string
	store    contacts.ContactStore
)

var rootCmd = &cobra.Command{
	Use:   "mini-crm",
	Short: "Mini CRM en ligne de commande (contacts JSON)",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialisation du store une seule fois pour toutes les sous-commandes
		if store == nil {
			// Crée le dossier si besoin
			dir := filepath.Dir(dataPath)
			if dir != "." && dir != "" {
				_ = os.MkdirAll(dir, 0o755)
			}
			store = contacts.NewStore(dataPath)
		}
		return nil
	},
}

func Execute() {
	// Flag global pour toutes les sous-commandes
	rootCmd.PersistentFlags().StringVarP(&dataPath, "data", "d", "data/contacts.json", "Chemin du fichier JSON de stockage")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Erreur:", err)
		os.Exit(1)
	}
}
