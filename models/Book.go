package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title   string   `json:"title"`
	Authors []Author `gorm:"many2many:book_authors;" json:"authors"`
}
