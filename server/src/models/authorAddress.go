package models

import (
	"gorm.io/gorm"
)

type AuthorAddress struct {
	gorm.Model
	ID        uint           `gorm:"primaryKey"`
	Street    string         `json:"street"`
	Town      string         `json:"town"`
	City      string         `json:"city"`
	Country   string         `json:"country"`
	AuthorID  uint           `gorm:"foreignKey:author_id"`
}
