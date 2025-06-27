package cli

import (
	"fmt"
	"log"

	// Pour valider le format de l'URL
	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/spf13/cobra"
	// Driver SQLite pour GORM
)

// longURLFlag stockera la valeur du flag --url
var longURLFlag string

// CreateCmd représente la commande 'create'
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Crée une URL courte à partir d'une URL longue.",
	Long: `Cette commande raccourcit une URL longue fournie et affiche le code court généré.

Exemple:
  url-shortener create --url="https://www.google.com/search?q=go+lang"`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO 1: Valider que le flag --url a été fourni.

		// TODO Validation basique du format de l'URL avec le package url et la fonction ParseRequestURI
		// si erreur, os.Exit(1)

		// TODO : Charger la configuration chargée globalement via cmd.cfg

		// TODO : Initialiser la connexion à la base de données SQLite.

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
		}

		// TODO S'assurer que la connexion est fermée à la fin de l'exécution de la commande

		// TODO : Initialiser les repositories et services nécessaires NewLinkRepository & NewLinkService

		// TODO : Appeler le LinkService et la fonction CreateLink pour créer le lien court.
		// os.Exit(1) si erreur

		fullShortURL := fmt.Sprintf("%s/%s", cfg.Server.BaseURL, link.ShortCode)
		fmt.Printf("URL courte créée avec succès:\n")
		fmt.Printf("Code: %s\n", link.ShortCode)
		fmt.Printf("URL complète: %s\n", fullShortURL)
	},
}

// init() s'exécute automatiquement lors de l'importation du package.
// Il est utilisé pour définir les flags que cette commande accepte.
func init() {
	// Définir le flag --url pour la commande create
	CreateCmd.Flags().StringVarP(&longURLFlag, "url", "u", "", "URL longue à raccourcir (requis)")

	// Marquer le flag comme requis
	CreateCmd.MarkFlagRequired("url")

	// Ajouter la commande à RootCmd
	cmd2.RootCmd.AddCommand(CreateCmd)
}
