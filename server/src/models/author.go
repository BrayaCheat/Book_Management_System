package models

import (
	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	ID            uint           `gorm:"primaryKey"`
	Name          string         `json:"name"`
	Gender        string         `json:"gender"`
	Books         []Book         `gorm:"foreignKey:author_id;"`
	AuthorAddress AuthorAddress  `gorm:"foreignKey:author_id;"`
}
