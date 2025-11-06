package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thomas-stecinski/crm_go_project/internal/contacts"
	"github.com/thomas-stecinski/crm_go_project/internal/db"
)

var (
	dataPath string // JSON path
	backend  string // "json" or "gorm"
	dsn      string // sqlite DSN
	store    contacts.ContactStore
	cfgFile  string // optional custom config path
)

var rootCmd = &cobra.Command{
	Use:   "mini-crm",
	Short: "Mini CRM CLI (contacts JSON / SQLite)",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Charger la configuration AVANT d'initialiser le store
		if err := initConfig(); err != nil {
			return fmt.Errorf("config: %w", err)
		}

		// Appliquer les valeurs de config (si non définies par flag)
		if backend == "" {
			backend = viper.GetString("backend")
		}
		if dataPath == "" {
			dataPath = viper.GetString("data")
		}
		if dsn == "" {
			dsn = viper.GetString("dsn")
		}

		// === Initialisation du store selon le backend ===
		switch backend {
		case "json", "":
			dir := filepath.Dir(dataPath)
			if dir != "." && dir != "" {
				_ = os.MkdirAll(dir, 0o755)
			}
			store = contacts.NewStore(dataPath)
			fmt.Println("Backend JSON chargé:", dataPath)
			return nil

		case "gorm":
			if dsn == "" {
				dsn = "data/contacts.db"
			}
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
			fmt.Println("Backend GORM chargé:", dsn)
			return nil

		default:
			return fmt.Errorf("backend inconnu: %s (attendu: json|gorm)", backend)
		}
	},
}

func Execute() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Fichier de configuration (ex: config.yaml)")
	rootCmd.PersistentFlags().StringVar(&backend, "backend", "", "Type de stockage (json|gorm)")
	rootCmd.PersistentFlags().StringVar(&dataPath, "data", "", "Chemin JSON (backend=json)")
	rootCmd.PersistentFlags().StringVar(&dsn, "dsn", "", "Chemin SQLite (backend=gorm)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Erreur:", err)
		os.Exit(1)
	}
}

// 🔧 Gestion automatique du fichier de config avec Viper
func initConfig() error {
	if cfgFile != "" {
		// fichier passé via --config
		viper.SetConfigFile(cfgFile)
	} else {
		// config.yaml par défaut dans le répertoire courant
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// variables d’environnement compatibles
	viper.SetEnvPrefix("CRM")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Aucun fichier config trouvé (config.yaml), utilisation des flags/env.")
			return nil
		}
		return err
	}
	fmt.Println("Configuration chargée depuis:", viper.ConfigFileUsed())
	return nil
}
