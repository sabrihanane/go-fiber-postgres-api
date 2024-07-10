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

func GetAuthors(c *fiber.Ctx) error {
	var authors []models.Author
	database.DB.Preload("Books").Find(&authors)
	return c.JSON(authors)
}

func GetAuthorById(c *fiber.Ctx) error {
	var author models.Author
	id := c.Params("id")
	authorId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid id"})
	}

	result := database.DB.Preload("Books").First(&author, authorId)
	fmt.Println("Author = ", author)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Author not found"})
	}
	fmt.Println("Author books = ", author.Books)

	authorBookTitles := &dto.AuthorDto{
		ID:         author.ID,
		Name:       author.Name,
		BookTitles: make([]string, len(author.Books)),
	}

	fmt.Println("authorBookTitles = ", len(authorBookTitles.BookTitles))

	for i, book := range author.Books {
		authorBookTitles.BookTitles[i] = book.Title
	}
	return c.Status(fiber.StatusOK).JSON(authorBookTitles)

}

func CreateAuthor(c *fiber.Ctx) error {
	author := new(models.Author)
	var authorModel models.Author
	// c.BodyParser() method is used to parse the request body into the provided Go struct, if the parsing is successful it returns nil
	if err := c.BodyParser(author); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	if author.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Author name can not be empty!"})
	}

	if len(author.Name) < 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Author name can not be less than 2 charachters"})
	}
	if validators.IsNumeric(author.Name) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": " Author name can not be numeric"})
	}

	result := database.DB.Where("name = ?", author.Name).First(&authorModel)
	if result.Error == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "An author with the same name already exists"})
	}

	database.DB.Create(&author)
	return c.JSON(author)
}

func UpdateAuthor(c *fiber.Ctx) error {
	author := new(models.Author)
	if err := c.BodyParser(author); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	result := database.DB.Model(&models.Author{}).Where("id = ?", author.ID).Updates(author)
	if result.Error != nil {
		return c.Status(fiber.StatusNotModified).JSON(fiber.Map{"error": "Updating book failed"})
	}
	return c.JSON(author)
}

func DeleteAuthorById(c *fiber.Ctx) error {
	author := new(models.Author)
	id := c.Params("id")
	if err := database.DB.First(&author, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Author not found",
		})
	}

	if err := database.DB.Delete(&author).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete author",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "The author has been deleted"})
}

func DeleteAuthorAndAssociationsById(c *fiber.Ctx) error {
	author := new(models.Author)
	id := c.Params("id")
	confirm := c.Query("confirm")

	if err := database.DB.Preload("Books").First(&author, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Author not found",
		})
	}

	if confirm != "yes" {
		return c.JSON(fiber.Map{
			"message":     " When deleting an author it automaticlly deletes the relationship with the assigned books, Are you sure you want to delete all of those relationships?",
			"confirm_url": c.BaseURL() + c.Path() + "?confirm=yes",
		})
	}

	err := database.DB.Model(&author).Association("Books").Clear()
	if err != nil {
		return c.Status(fiber.StatusNotModified).JSON(fiber.Map{"error": "Failed to delete associated authors"})
	}

	err1 := database.DB.Delete(&author).Error
	if err1 != nil {
		return c.Status(fiber.StatusNotModified).JSON(fiber.Map{"error": "Deleting author failed"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": "Author and its associations deleted succefully"})

}

func DeleteAuthorAndAssociatedBooksById(c *fiber.Ctx) error {
	// 1- get all the parameters
	var associations []models.BookAuthor
	author := new(models.Author)
	id := c.Params("id")
	confirm := c.Query("confirm")

	//2- check the existence of the author to be deleted
	if err := database.DB.Preload("Books").First(&author, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Author not found",
		})
	}
	fmt.Println("Author = ", author)

	// 3- check if the user confirmed to delete the book and its associations
	if confirm != "yes" {
		return c.JSON(fiber.Map{
			"message":     " When deleting an author it automaticlly deletes associated books, Are you sure you want to delete all the associations and the associated books?",
			"confirm_url": c.BaseURL() + c.Path() + "?confirm=yes",
		})
	}

	// Retrieve the author association array from the book/author mapping table
	if err := database.DB.Where("author_id = ?", id).Find(&associations).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}
	fmt.Println("Associations = ", associations)

	// Loop through the book associations array
	for _, association := range associations {
		var associationCount int64
		// Retrieve the book associations list and count the number of authors associated with each of the books
		if err := database.DB.Model(&models.BookAuthor{}).Where("book_id = ?", association.BookId).Count(&associationCount).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}
		fmt.Println("count = ", associationCount)

		// Check if the book has more than one author, if not, delete the book.
		if associationCount == 1 {
			if err := database.DB.Delete(&models.Book{}, association.BookId).Error; err != nil {
				return c.Status(fiber.StatusNotFound).JSON(err.Error())
			}
		}

		// Delete the association from author/book mapping table
		if err := database.DB.Where("book_id = ? AND author_id = ?", &association.BookId, &association.AuthorId).Delete(&models.BookAuthor{}).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		// Delete the author
		if err := database.DB.Delete(&author).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Success": " and its associations deleted succefully"})
}
