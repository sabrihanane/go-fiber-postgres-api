package handlers

import (
	"BookAuthor_ManyToMany/models"
	"BookAuthor_ManyToMany/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func LogIn(c *fiber.Ctx) error {
	var users []models.User
	users = append(users, models.User{ID: 1, UserName: "hanane_sabri", Password: "hanane12345", FirstName: "Hanane", LastName: "Sabri"})
	users = append(users, models.User{ID: 2, UserName: "john_smith", Password: "smith12345", FirstName: "john", LastName: "Smith"})
	users = append(users, models.User{ID: 3, UserName: "james_walker", Password: "james12345", FirstName: "James", LastName: "Walker"})

	credentials := new(models.Credentials)
	if err := c.BodyParser(credentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	user := searchUser(users, credentials.UserName, credentials.Password)
	if user == nil {
		fmt.Println("User not found")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Username or password is incorrect",
		})
	}

	fmt.Println("User found: ", user)

	token, err := utils.GenerateToken(user.ID, user.UserName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.JSON(fiber.Map{
		"accessToken": token,
	})
}

func searchUser(users []models.User, username string, password string) *models.User {
	for _, user := range users {
		if user.UserName == username && user.Password == password {
			return &user
		}
	}
	return nil
}
