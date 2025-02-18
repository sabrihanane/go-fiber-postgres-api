package routes

import (
	"BookAuthor_ManyToMany/handlers"
	"BookAuthor_ManyToMany/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	userGroup := app.Group("/user")
	userGroup.Get("/login", handlers.LogIn)

	bookGroup := app.Group("/book", middleware.JWTProtected())
	bookGroup.Get("/get_books", handlers.GetBooks)
	bookGroup.Get("/get_book_by_id/:id", handlers.GetBookById)
	bookGroup.Get("/get_book_authors_by_id/:id", handlers.AuthorsOfaSpecificBook)
	bookGroup.Post("/create_books", handlers.CreateBook)
	bookGroup.Put("/update_book", handlers.UpdateBook)
	bookGroup.Delete("/delete_book/:id", handlers.DeleteBookById)
	bookGroup.Delete("/delete_book_and_associations_by_id/:id", handlers.DeleteBookAndAssociationsById)
	bookGroup.Delete("/delete_book_and_associated_authors_by_id/:id", handlers.DeleteBookAndAssociatedAuthorsById)
	bookGroup.Post("/assign_author_to_book_by_ids", handlers.AssignAuthorToBookByIds)
	bookGroup.Delete("/unassign_author_from_book_by_ids", handlers.UnassignAuthorFromBookByIds)

	authorGroup := app.Group("/author")
	authorGroup.Get("/get_authors", handlers.GetAuthors)
	authorGroup.Get("/get_author_by_id/:id", handlers.GetAuthorById)
	authorGroup.Post("/create_authors", handlers.CreateAuthor)
	authorGroup.Put("/update_author", handlers.UpdateAuthor)
	authorGroup.Delete("delete_author/:id", handlers.DeleteAuthorById)
	authorGroup.Delete("delete_author_and_associations_by_id/:id", handlers.DeleteAuthorAndAssociationsById)
	authorGroup.Delete("delete_author_and_associated_books_by_id/:id", handlers.DeleteAuthorAndAssociatedBooksById)
}
