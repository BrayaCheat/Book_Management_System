package models

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	ID             uint           `gorm:"primaryKey"`
	Title          string         `json:"title"`
	PubslishedDate string         `json:"published_date"`
	BookCover      string         `json:"book_cover"`
	AuthorID       uint           `gorm:"foreignKey:author_id"`
}
