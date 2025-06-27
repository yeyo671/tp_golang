package config

import (
	"fmt"
	"log" // Pour logger les informations ou erreurs de chargement de config

	"github.com/spf13/viper" // La bibliothèque pour la gestion de configuration
)

// TODO Créer Config qui est la structure principale qui mappe l'intégralité de la configuration de l'application.
// Les tags `mapstructure` sont utilisés par Viper pour mapper les clés du fichier de config
// (ou des variables d'environnement) aux champs de la structure Go.
type Config struct {
	Server struct {
		Port     int    `mapstructure:"port"`      // Port du serveur HTTP
		BaseURL  string `mapstructure:"base_url"` // URL de base du serveur
	} `mapstructure:"server"` // Sous-structure pour la configuration du serveur

	Database struct {
		Name     string `mapstructure:"name"`     // Nom de la base de donnéesAdd commentMore actions
	} `mapstructure:"database"` // Sous-structure pour la configuration de la base de données

	Analytics struct {
		BufferSize int `mapstructure:"buffer_size"` // Taille du tampon pour les données
		WorkerCount int `mapstructure:"worker_count"` // Nombre de workers pour traiter les données
	} `mapstructure:"analytics"` // Sous-structure pour la configuration des analytics

	Monitor struct {
		IntervalMinutes int `mapstructure:"interval_minutes"` // Intervalle de surveillance en minutes
	} `mapstructure:"monitor"` // Sous-structure pour la configuration du moniteur
}

// LoadConfig charge la configuration de l'application en utilisant Viper.
// Elle recherche un fichier 'config.yaml' dans le dossier 'configs/'.
// Elle définit également des valeurs par défaut si le fichier de config est absent ou incomplet.
func LoadConfig() (*Config, error) {
	// TODO Spécifie le chemin où Viper doit chercher les fichiers de config.
	// on cherche dans le dossier 'configs' relatif au répertoire d'exécution.
	viper.AddConfigPath("./configs") // Chemin relatif au répertoire d'exécution

	// TODO Spécifie le nom du fichier de config (sans l'extension).
	viper.SetConfigName("config")     // Nom du fichier de configuration sans l'extension

	// TODO Spécifie le type de fichier de config.
	viper.SetConfigType("yaml")

	// TODO : Définir les valeurs par défaut pour toutes les options de configuration.
	// Ces valeurs seront utilisées si les clés correspondantes ne sont pas trouvées dans le fichier de config
	// ou si le fichier n'existe pas.
	// server.port, server.base_url etc.
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.base_url", "http://localhost:8080")
	
	viper.SetDefault("database.name", "url_shortener.db")

	viper.SetDefault("analytics.buffer_size", 1000)
	viper.SetDefault("analytics.worker_count", 5)

	viper.SetDefault("monitor.interval_minutes", 5)

	// TODO : Lire le fichier de configuration.
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Erreur lors de la lecture du fichier de configuration: %v", err)
		return nil, err
	}

	// TODO 4: Démapper (unmarshal) la configuration lue (ou les valeurs par défaut) dans la structure Config.
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Printf("Erreur lors du démapping de la configuration: %v", err)
		return nil, err
	}

	// Log  pour vérifier la config chargée
	log.Printf("Configuration loaded: Server Port=%d, DB Name=%s, Analytics Buffer=%d, Monitor Interval=%dmin",
		cfg.Server.Port, cfg.Database.Name, cfg.Analytics.BufferSize, cfg.Monitor.IntervalMinutes)

	return &cfg, nil // Retourne la configuration chargée
}
