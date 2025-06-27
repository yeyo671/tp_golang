package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"gorm.io/gorm" // Nécessaire pour la gestion spécifique de gorm.ErrRecordNotFound

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository" // Importe le package repository
)

// Définition du jeu de caractères pour la génération des codes courts.
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// LinkService est une structure qui fournit des méthodes pour la logique métier des liens.
// Elle détient linkRepo qui est une référence vers une interface LinkRepository.
type LinkService struct {
	linkRepo repository.LinkRepository // Interface pour accéder aux méthodes du repository
}

// NewLinkService crée et retourne une nouvelle instance de LinkService.
func NewLinkService(linkRepo repository.LinkRepository) *LinkService {
	return &LinkService{
		linkRepo: linkRepo,
	}
}

// GenerateShortCode génère un code court aléatoire d'une longueur spécifiée.
func (s *LinkService) GenerateShortCode(length int) (string, error) {
	// Génère un code court aléatoire sécurisé de la longueur spécifiée
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}

// CreateLink crée un nouveau lien raccourci.
func (s *LinkService) CreateLink(longURL string) (*models.Link, error) {
	const codeLength = 6
	const maxRetries = 5
	var shortCode string
	var err error

	// Tentatives pour générer un code unique
	for i := 0; i < maxRetries; i++ {
		var code string
		code, err = s.GenerateShortCode(codeLength)
		if err != nil {
			return nil, fmt.Errorf("failed to generate short code: %w", err)
		}

		// Vérifie si le code existe déjà
		_, err = s.linkRepo.GetLinkByShortCode(code)
		if err != nil {
			// Si l'erreur est 'record not found', le code est unique
			if errors.Is(err, gorm.ErrRecordNotFound) {
				shortCode = code
				break
			}
			// Autre erreur DB
			return nil, fmt.Errorf("database error checking short code uniqueness: %w", err)
		}
		// Si aucune erreur, le code existe déjà, on retente
	}

	if shortCode == "" {
		return nil, errors.New("could not generate a unique short code after several attempts")
	}

	// Crée une nouvelle instance du modèle Link
	link := &models.Link{
		ShortCode: shortCode,
		LongURL:   longURL,
		// CreatedAt sera géré automatiquement par GORM
	}

	// Persiste le nouveau lien dans la base de données via le repository
	if err := s.linkRepo.CreateLink(link); err != nil {
		return nil, fmt.Errorf("failed to persist new link: %w", err)
	}

	return link, nil
}

// GetLinkByShortCode récupère un lien via son code court.
func (s *LinkService) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	// Utilise le repository pour récupérer le lien par son code court
	link, err := s.linkRepo.GetLinkByShortCode(shortCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get link by shortcode: %w", err)
	}
	return link, nil
}

// GetLinkStats récupère les statistiques pour un lien donné (nombre total de clics).
func (s *LinkService) GetLinkStats(shortCode string) (*models.Link, int, error) {
	// Récupère le lien par son shortCode
	link, err := s.linkRepo.GetLinkByShortCode(shortCode)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get link by shortcode: %w", err)
	}

	// Compte le nombre de clics pour ce LinkID
	clicks, err := s.linkRepo.CountClicksByLinkID(link.ID)
	if err != nil {
		return link, 0, fmt.Errorf("failed to count clicks for link: %w", err)
	}

	return link, clicks, nil
}