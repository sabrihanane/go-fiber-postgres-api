package routes

import (
	"BookAuthor_ManyToMany/handlers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	bookGroup := app.Group("/book")
	bookGroup.Get("/get_books", handlers.GetBooks)
	bookGroup.Get("/get_book_by_id/:id", handlers.GetBookById)
	bookGroup.Get("/get_book_authors_by_id/:id", handlers.AuthorsOfaSpecificBook)
	bookGroup.Post("/create_books", handlers.CreateBook)
	bookGroup.Put("/update_book", handlers.UpdateBook)
	bookGroup.Delete("/delete_book/:id", handlers.DeleteBookById)
	bookGroup.Delete("/delete_book_and_associations_by_id/:id", handlers.DeleteBookAndAssociationsById)
	//bookGroup.Delete("/delete_book_and_associated_authors_by_id/:id", handlers.DeleteBookAndAssociatedAuthorsById)
	bookGroup.Post("/assign_author_to_book_by_ids", handlers.AssignAuthorToBookByIds)
	bookGroup.Post("/unassign_author_from_book_by_ids", handlers.UnassignAuthorFromBookByIds)

	athorGroup := app.Group("/author")
	athorGroup.Get("/get_authors", handlers.GetAuthors)
	athorGroup.Get("/get_author_by_id/:id", handlers.GetAuthorById)
	athorGroup.Post("/create_authors", handlers.CreateAuthor)
}
