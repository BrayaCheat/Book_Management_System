package routes

import (
	"server/src/controllers"
	"server/src/middlewares"
	"github.com/gofiber/fiber/v2"
)

func AuthorRoute(app *fiber.App){

	//author
	author := app.Group("/api/author")
	author.Post("/", middlewares.AuthMiddleware, controllers.CreateAuthor)
	author.Get("/", controllers.ListAuthor)
	author.Get("/:id", controllers.GetAuthor)
	author.Put("/:id", middlewares.AuthMiddleware, controllers.UpdateAuthor)
	author.Delete("/:id", middlewares.AuthMiddleware, controllers.DeleteAuthor)

	//address
	address := app.Group("/api/address")
	address.Post("/", middlewares.AuthMiddleware, controllers.CreateAddress)
	address.Get("/", controllers.ListAddress)
	address.Get("/:id", controllers.GetAddress)
	address.Put("/:id", middlewares.AuthMiddleware, controllers.UpdateAddress)
	address.Delete("/:id", middlewares.AuthMiddleware, controllers.DeleteAddress)
}