package dto

type AuthorDto struct {
	ID         uint     `json:"id"`
	Name       string   `json:"name"`
	BookTitles []string `json:"bookTitles"`
}
