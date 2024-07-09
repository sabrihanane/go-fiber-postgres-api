package models

type BookAuthor struct {
	BookId   uint `json:"book_id" gorm:"book_id"`
	AuthorId uint `json:"author_id" gorm:"author_id"`
}
