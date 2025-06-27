package models

import "time"

// TODO : Créer la struct Link
// Link représente un lien raccourci dans la base de données.
// Les tags `gorm:"..."` définissent comment GORM doit mapper cette structure à une table SQL.
// ID qui est une primaryKey
// Shortcode : doit être unique, indexé pour des recherches rapide (voir doc), taille max 10 caractères
// LongURL : doit pas être null
// CreateAt : Horodatage de la créatino du lien

type Link struct {
	ID        uint   `gorm:"primaryKey"`        // Clé primaire
	Shortcode string `gorm:"uniqueIndex;size:10"` // Code court unique, indexé pour des recherches rapides, taille maximale de 10 caractères
	LongURL   string `gorm:"not null"` // URL longue, ne peut pas être nulle
	CreatedAt time.Time `gorm:"autoCreateTime"` // Horodatage de création, automatiquement défini par GORM
}