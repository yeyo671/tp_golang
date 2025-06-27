package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"gorm.io/gorm" // Nécessaire pour la gestion spécifique de gorm.ErrRecordNotFound

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository" // Importe le package repository
)

// Définition du jeu de caractères pour la génération des codes courts.
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// TODO Créer la struct
// LinkService est une structure qui g fournit des méthodes pour la logique métier des liens.
// Elle détient linkRepo qui est une référence vers une interface LinkRepository.
// IMPORTANT : Le champ doit être du type de l'interface (non-pointeur).


// NewLinkService crée et retourne une nouvelle instance de LinkService.
func NewLinkService(linkRepo repository.LinkRepository) *LinkService {
	return &LinkService{
		linkRepo: linkRepo,
	}
}

// TODO Créer la méthode GenerateShortCode
// GenerateShortCode est une méthode rattachée à LinkService
// Elle génère un code court aléatoire d'une longueur spécifiée. Elle prend une longueur en paramètre et retourne une string et une erreur
// Il utilise le package 'crypto/rand' pour éviter la prévisibilité.
// Je vous laisse chercher un peu :) C'est faisable en une petite dizaine de ligne


// CreateLink crée un nouveau lien raccourci.
// Il génère un code court unique, puis persiste le lien dans la base de données.
func (s *LinkService) CreateLink(longURL string) (*models.Link, error) {
    const maxRetries = 5
    var shortCode string

    for i := 0; i < maxRetries; i++ {
        code, err := s.GenerateShortCode(6)
        if err != nil {
            return nil, fmt.Errorf("failed to generate short code: %w", err)
        }

        _, err = s.GetLinkByShortCode(code)
        if errors.Is(err, gorm.ErrRecordNotFound) {
            shortCode = code
            break
        } else if err != nil {
            return nil, fmt.Errorf("database error: %w", err)
        }
        // Collision, retry
    }

    if shortCode == "" {
        return nil, errors.New("could not generate unique short code after several attempts")
    }

    link := &models.Link{
        ShortCode: shortCode,
        LongURL:   longURL,
    }

    if err := s.repository.CreateLink(link); err != nil {
        return nil, fmt.Errorf("failed to save link: %w", err)
    }
	
    return link, nil
}

// GetLinkByShortCode récupère un lien via son code court.
// Il délègue l'opération de recherche au repository.
func (s *LinkService) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	// TODO : Récupérer un lien par son code court en utilisant s.linkRepo.GetLinkByShortCode.
	// Retourner le lien trouvé ou une erreur si non trouvé/problème DB.

}

// GetLinkStats récupère les statistiques pour un lien donné (nombre total de clics).
// Il interagit avec le LinkRepository pour obtenir le lien, puis avec le ClickRepository
func (s *LinkService) GetLinkStats(shortCode string) (*models.Link, int, error) {
	// TODO : Récupérer le lien par son shortCode


	// TODO 4: Compter le nombre de clics pour ce LinkID

	// TODO : on retourne les 3 valeurs
	return
}

