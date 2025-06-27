package workers

import (
	"log"

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository" // Nécessaire pour interagir avec le ClickRepository
)

// StartClickWorkers lance un pool de goroutines "workers" pour traiter les événements de clic.
// Chaque worker lira depuis le même 'clickEventsChan' et utilisera le 'clickRepo' pour la persistance.
func StartClickWorkers(workerCount int, clickEventsChan <-chan models.ClickEvent, clickRepo repository.ClickRepository) {
	log.Printf("Starting %d click worker(s)...", workerCount)
	for i := 0; i < workerCount; i++ {
		// Lance chaque worker dans sa propre goroutine.
		// Le channel est passé en lecture seule (<-chan) pour renforcer l'immutabilité du channel à l'intérieur du worker.
		go clickWorker(clickEventsChan, clickRepo)
	}
}

// clickWorker est la fonction exécutée par chaque goroutine worker.
// Elle tourne indéfiniment, lisant les événements de clic dès qu'ils sont disponibles dans le channel.
func clickWorker(clickEventsChan <-chan models.ClickEvent, clickRepo repository.ClickRepository) {
	for event := range clickEventsChan { // Boucle qui lit les événements du channel
		// Conversion du ClickEvent en modèle Click
		click := models.Click{
			LinkID:    event.LinkID,
			Timestamp: event.Timestamp,
			UserAgent: event.UserAgent,
			IPAddress: event.IPAddress,
		}

		// Persister le clic en base de données via le clickRepo
		err := clickRepo.CreateClick(&click)

		if err != nil {
			// Si une erreur se produit lors de l'enregistrement, logguez-la.
			log.Printf("ERROR: Failed to save click for LinkID %d (UserAgent: %s, IP: %s): %v",
				event.LinkID, event.UserAgent, event.IPAddress, err)

		} else {
			// Log optionnel pour confirmer l'enregistrement (utile pour le débogage)
			log.Printf("Click recorded successfully for LinkID %d", event.LinkID)
		}
	}
}