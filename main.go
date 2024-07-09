package main

import (
	"BookAuthor_ManyToMany/database"
	"BookAuthor_ManyToMany/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	database.Connect()
	database.AutoMigrate()

	routes.Setup(app)

	app.Listen(":3000")
}
