package dto

type BookDto struct {
	ID          uint     `json:"id"`
	Title       string   `json:"title"`
	AuthorNames []string `json:"authorNames"`
}
