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
