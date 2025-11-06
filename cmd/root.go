package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/thomas-stecinski/crm_go_project/internal/contacts"
	"github.com/thomas-stecinski/crm_go_project/internal/db"
)

var (
	dataPath string // JSON file path (backend=json)
	backend  string // "json" (par défaut) ou "gorm"
	dsn      string // sqlite dsn (backend=gorm)
	store    contacts.ContactStore
)

var rootCmd = &cobra.Command{
	Use:   "mini-crm",
	Short: "Mini CRM en ligne de commande (contacts JSON / SQLite)",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if store != nil {
			return nil
		}
		switch backend {
		case "json", "":
			// Prépare le dossier si besoin
			dir := filepath.Dir(dataPath)
			if dir != "." && dir != "" {
				_ = os.MkdirAll(dir, 0o755)
			}
			store = contacts.NewStore(dataPath)
			return nil

		case "gorm":
			if dsn == "" {
				// valeur par défaut
				dsn = "data/contacts.db"
			}
			// Crée le dossier si besoin
			dir := filepath.Dir(dsn)
			if dir != "." && dir != "" {
				_ = os.MkdirAll(dir, 0o755)
			}

			gdb, err := db.OpenSQLite(dsn)
			if err != nil {
				return err
			}
			gormStore, err := contacts.NewGormStore(gdb)
			if err != nil {
				return err
			}
			store = gormStore
			return nil

		default:
			return fmt.Errorf("backend inconnu: %s (attendu: json|gorm)", backend)
		}
	},
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&dataPath, "data", "d", "data/contacts.json", "Chemin du fichier JSON (backend=json)")
	rootCmd.PersistentFlags().StringVar(&backend, "backend", "json", "Type de stockage: json | gorm")
	rootCmd.PersistentFlags().StringVar(&dsn, "dsn", "", "DSN SQLite (ex: data/contacts.db) (backend=gorm)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Erreur:", err)
		os.Exit(1)
	}
}
