package api

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm" // Pour gérer gorm.ErrRecordNotFound
)

// ClickEventsChannel est le channel global (ou injecté) utilisé pour envoyer les événements de clic
// aux workers asynchrones. Il est bufferisé pour ne pas bloquer les requêtes de redirection.
var ClickEventsChannel chan models.ClickEvent

// SetupRoutes configure toutes les routes de l'API Gin et injecte les dépendances nécessaires
func SetupRoutes(router *gin.Engine, linkService *services.LinkService, bufferSize int, baseURL string) {
	// Le channel est initialisé ici.
	if ClickEventsChannel == nil {
		// Créer le channel ici (make), il doit être bufférisé
		// La taille du buffer doit être configurable via Viper (cfg.Analytics.BufferSize)
		ClickEventsChannel = make(chan models.ClickEvent, bufferSize)
	}

	// Route de Health Check , /health
	router.GET("/health", HealthCheckHandler)

	// Routes de l'API
	// Doivent être au format /api/v1/
	// POST /links
	// GET /links/:shortCode/stats
	api := router.Group("/api/v1")
	{
		api.POST("/links", CreateShortLinkHandler(linkService, baseURL))
		api.GET("/links/:shortCode/stats", GetLinkStatsHandler(linkService))
	}

	// Route de Redirection (au niveau racine pour les short codes)
	router.GET("/:shortCode", RedirectHandler(linkService))
}

// HealthCheckHandler gère la route /health pour vérifier l'état du service.
func HealthCheckHandler(c *gin.Context) {
	// Retourner simplement du JSON avec un StatusOK, {"status": "ok"}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// CreateLinkRequest représente le corps de la requête JSON pour la création d'un lien.
type CreateLinkRequest struct {
	LongURL string `json:"long_url" binding:"required,url"` // 'binding:required' pour validation, 'url' pour format URL
}

// CreateShortLinkHandler gère la création d'une URL courte.
func CreateShortLinkHandler(linkService *services.LinkService, baseURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateLinkRequest
		// Tente de lier le JSON de la requête à la structure CreateLinkRequest.
		// Gin gère la validation 'binding'.
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Appeler le LinkService (CreateLink pour créer le nouveau lien.
		link, err := linkService.CreateLink(req.LongURL)
		if err != nil {
			log.Printf("Error creating link: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create link"})
			return
		}

		// Retourne le code court et l'URL longue dans la réponse JSON.
		c.JSON(http.StatusCreated, gin.H{
			"short_code":     link.ShortCode,
			"long_url":       link.LongURL,
			"full_short_url": baseURL + "/" + link.ShortCode,
		})
	}
}

// RedirectHandler gère la redirection d'une URL courte vers l'URL longue et l'enregistrement asynchrone des clics.
func RedirectHandler(linkService *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Récupère le shortCode de l'URL avec c.Param
		shortCode := c.Param("shortCode")

		// Récupérer l'URL longue associée au shortCode depuis le linkService (GetLinkByShortCode)
		link, err := linkService.GetLinkByShortCode(shortCode)
		if err != nil {
			// Si le lien n'est pas trouvé, retourner HTTP 404 Not Found.
			// Utiliser errors.Is et l'erreur Gorm
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
				return
			}
			// Gérer d'autres erreurs potentielles de la base de données ou du service
			log.Printf("Error retrieving link for %s: %v", shortCode, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Créer un ClickEvent avec les informations pertinentes.
		clickEvent := models.ClickEvent{
			LinkID:    link.ID,
			Timestamp: time.Now(),
			UserAgent: c.GetHeader("User-Agent"),
			IPAddress: c.ClientIP(),
		}

		// Envoyer le ClickEvent dans le ClickEventsChannel avec le Multiplexage.
		// Utilise un `select` avec un `default` pour éviter de bloquer si le channel est plein.
		select {
		case ClickEventsChannel <- clickEvent:
			// Succès: l'événement a été envoyé
		default:
			// Pour le default, juste un message à afficher :
			log.Printf("Warning: ClickEventsChannel is full, dropping click event for %s.", shortCode)
		}

		// Effectuer la redirection HTTP 302 (StatusFound) vers l'URL longue.
		c.Redirect(http.StatusFound, link.LongURL)
	}
}

// GetLinkStatsHandler gère la récupération des statistiques pour un lien spécifique.
func GetLinkStatsHandler(linkService *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Récupère le shortCode de l'URL avec c.Param
		shortCode := c.Param("shortCode")

		// Appeler le LinkService pour obtenir le lien et le nombre total de clics.
		link, totalClicks, err := linkService.GetLinkStats(shortCode)
		if err != nil {
			// Gérer le cas où le lien n'est pas trouvé.
			// toujours avec l'erreur Gorm ErrRecordNotFound
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
				return
			}
			// Gérer d'autres erreurs
			log.Printf("Error retrieving stats for %s: %v", shortCode, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Retourne les statistiques dans la réponse JSON.
		c.JSON(http.StatusOK, gin.H{
			"short_code":   link.ShortCode,
			"long_url":     link.LongURL,
			"total_clicks": totalClicks,
		})
	}
}
