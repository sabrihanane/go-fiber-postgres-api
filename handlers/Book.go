package handlers

import (
	"BookAuthor_ManyToMany/database"
	"BookAuthor_ManyToMany/dto"
	"BookAuthor_ManyToMany/models"
	"BookAuthor_ManyToMany/validators"

	"fmt"

	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetBooks(c *fiber.Ctx) error {
	var books []models.Book
	database.DB.Preload("Authors").Find(&books)
	return c.JSON(books)
}

func GetBookById(c *fiber.Ctx) error {
	var book models.Book
	//var bookAuthorNames dto.BookDto
	id := c.Params("id")
	bookId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid id"})
	}
	result := database.DB.Preload("Authors").First(&book, bookId)
	fmt.Println("Book = ", book)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found"})
	}
	fmt.Println("Book authors = ", book.Authors)

	// Create the book dto with an empty array of authors iIDs
	bookAuthorNames := &dto.BookDto{
		ID:          book.ID,
		Title:       book.Title,
		AuthorNames: make([]string, len(book.Authors)),
	}

	fmt.Println("bookAuthorNames.AuthorNames = ", len(bookAuthorNames.AuthorNames))

	for i, author := range book.Authors {
		bookAuthorNames.AuthorNames[i] = author.Name
	}
	return c.Status(fiber.StatusOK).JSON(bookAuthorNames)
}

func CreateBook(c *fiber.Ctx) error {
	book := new(models.Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	// Validate the data
	if book.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Book title can not be empty!"})
	}
	if len(book.Title) < 4 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Book title can not be less than 4 charachters"})
	}
	if validators.IsNumeric(book.Title) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Book title can not be numeric"})
	}

	var bookModel models.Book
	result := database.DB.Where("title = ?", book.Title).First(&bookModel)

	if result.Error == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "A book with the same title already exists"})
	}

	database.DB.Create(&book)
	return c.JSON(book)
}

func AuthorsOfaSpecificBook(c *fiber.Ctx) error {
	bookId := c.Params("id")

	var book models.Book
	var authors []dto.AuthorDto

	if err := database.DB.First(&book, bookId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found"})
	}

	if err := database.DB.Model(&book).Select([]string{"id", "name"}).Association("Authors").Find(&authors); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load Authors"})
	}

	return c.JSON(authors)
}

// UpdateBook updates a book in the database by their ID
func UpdateBook(c *fiber.Ctx) error {
	book := new(models.Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	result := database.DB.Model(&models.Book{}).Where("id = ?", book.ID).Updates(book)
	if result.Error != nil {
		return c.Status(fiber.StatusNotModified).JSON(fiber.Map{"error": "Updating book failed"})
	}
	return c.JSON(book)
}

func DeleteBookById(c *fiber.Ctx) error {
	book := new(models.Book)
	var associationCount int64

	id := c.Params("id")
	if err := database.DB.First(&book, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	associationCount = database.DB.Model(&book).Association("Authors").Count()

	if associationCount >= 1 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "This book is associated with authors and can not be deleted "})
	}

	if err := database.DB.Delete(&book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete book",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "The book has been deleted"})
}

func DeleteBookAndAssociationsById(c *fiber.Ctx) error {
	// 1- get all the parameters
	book := new(models.Book)
	id := c.Params("id")          //  required Used to identify specific resources within the API. They are mandatory for the endpoint to make sense and are part of the URL structure.
	confirm := c.Query("confirm") // used to pass key-value pairs to the server, Typically used for filtering, searching, and modifying the request. They can be optional and do not change the URL structure

	// 2- check the existence of the book to be deleted
	if err := database.DB.Preload("Authors").First(&book, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	// 3- check if the user confirmed to delete the book and its associations
	if confirm != "yes" {
		return c.JSON(fiber.Map{
			"message":     " When deleting a book it automaticlly deletes the relationship with the assigned authors, Are you sure you want to delete all of those relationships?",
			"confirm_url": c.BaseURL() + c.Path() + "?confirm=yes",
		})
	}

	// 4- Delete the book associations
	err := database.DB.Model(&book).Association("Authors").Clear()
	if err != nil {
		return c.Status(fiber.StatusNotModified).JSON(fiber.Map{"error": "Failed to delete associated authors"})
	}

	// 5- Delete the book
	err1 := database.DB.Delete(&book).Error
	if err1 != nil {
		return c.Status(fiber.StatusNotModified).JSON(fiber.Map{"error": "Deleting book failed"})
	}

	// 6- Return deletion confirmation message to the user
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": "Book and its associations deleted succefully"})

}

func DeleteBookAndAssociatedAuthorsById(c *fiber.Ctx) error {
	// 1- get all the parameters
	var associations []models.BookAuthor
	book := new(models.Book)
	id := c.Params("id")          //  required Used to identify specific resources within the API. They are mandatory for the endpoint to make sense and are part of the URL structure.
	confirm := c.Query("confirm") // used to pass key-value pairs to the server, Typically used for filtering, searching, and modifying the request. They can be optional and do not change the URL structure

	//2- check the existence of the book to be deleted
	if err := database.DB.Preload("Authors").First(&book, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}
	fmt.Println("Book = ", book)

	// 3- check if the user confirmed to delete the book and its associations
	if confirm != "yes" {
		return c.JSON(fiber.Map{
			"message":     " When deleting a book it automaticlly deletes associated authors, Are you sure you want to delete all the associations and the associated authors?",
			"confirm_url": c.BaseURL() + c.Path() + "?confirm=yes",
		})
	}

	// Retrieve the book association array from the book/author mapping table
	if err := database.DB.Where("book_id = ?", id).Find(&associations).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}
	fmt.Println("Associations = ", associations)

	// Loop through the book associations array
	for _, association := range associations {
		var associationCount int64
		// Retrieve the author associations list and count the number of books associated with each of the authors
		if err := database.DB.Model(&models.BookAuthor{}).Where("author_id = ?", association.AuthorId).Count(&associationCount).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		// Check if the author has more than one book, if not, delete the author.
		if associationCount == 1 {
			if err := database.DB.Delete(&models.Author{}, association.AuthorId).Error; err != nil {
				return c.Status(fiber.StatusNotFound).JSON(err.Error())
			}
		}

		// Delete the association from author/book mapping table
		if err := database.DB.Where("book_id = ? AND author_id = ?", &association.BookId, &association.AuthorId).Delete(&models.BookAuthor{}).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		// Delete the book
		if err := database.DB.Delete(&book).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Success": "Book and its associations deleted succefully"})
}

func AssignAuthorToBookByIds(c *fiber.Ctx) error {
	bookAuthor := new(models.BookAuthor)
	if err := c.BodyParser(bookAuthor); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	var associationCount int64
	if err := database.DB.Model(&models.BookAuthor{}).Where("book_id = ? AND author_id = ?", &bookAuthor.BookId, &bookAuthor.AuthorId).Count(&associationCount).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
	}

	if associationCount > 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"Message": "This author is already assigned to that specific book",
		})
	}

	var book models.Book
	if err := database.DB.First(&book, bookAuthor.BookId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	var author models.Author
	if err := database.DB.First(&author, bookAuthor.AuthorId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Author not found",
		})
	}

	if err := database.DB.Model(&book).Association("Authors").Append(&author); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Author assigned to book successfully",
	})
}

func UnassignAuthorFromBookByIds(c *fiber.Ctx) error {
	bookAuthor := new(models.BookAuthor)

	if err := c.BodyParser(&bookAuthor); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	var associationCount int64
	if err := database.DB.Model(&models.BookAuthor{}).Where("book_id = ? AND author_id = ?", &bookAuthor.BookId, &bookAuthor.AuthorId).Count(&associationCount).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
	}

	if associationCount < 1 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"Message": "This author is not assigned to this book",
		})
	}

	var book models.Book
	if err := database.DB.First(&book, bookAuthor.BookId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	var author models.Author
	if err := database.DB.First(&author, bookAuthor.AuthorId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Author not found",
		})
	}

	database.DB.Model(&book).Association("Authors").Delete(&author)

	return c.JSON(fiber.Map{
		"message": "Author unassigned from book successfully",
	})
}
