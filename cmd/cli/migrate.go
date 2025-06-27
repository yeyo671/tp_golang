package cli

import (
	"fmt"
	"log"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/spf13/cobra"
	// Driver SQLite pour GORM
)

// MigrateCmd représente la commande 'migrate'
var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Exécute les migrations de la base de données pour créer ou mettre à jour les tables.",
	Long: `Cette commande se connecte à la base de données configurée (SQLite)
et exécute les migrations automatiques de GORM pour créer les tables 'links' et 'clicks'
basées sur les modèles Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO : Charger la configuration chargée globalement via cmd.cfg

		// TODO 2: Initialiser la connexion à la base de données SQLite avec GORM.

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
		}
		// TODO Assurez-vous que la connexion est fermée après la migration.

		// TODO 3: Exécuter les migrations automatiques de GORM.
		// Utilisez db.AutoMigrate() et passez-lui les pointeurs vers tous vos modèles.

		// Pas touche au log
		fmt.Println("Migrations de la base de données exécutées avec succès.")
	},
}

func init() {
	// Ajouter la commande à RootCmd
	cmd2.RootCmd.AddCommand(MigrateCmd)
}
