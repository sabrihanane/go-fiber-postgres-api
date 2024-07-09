package models

import "gorm.io/gorm"

type Author struct {
	gorm.Model
	Name  string `json:"name"`
	Books []Book `gorm:"many2many:book_authors;" json:"books"`
}
