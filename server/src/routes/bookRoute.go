package routes

import (
	"server/src/controllers"
	"server/src/middlewares"
	"github.com/gofiber/fiber/v2"
)

func BookRoute(app *fiber.App){
	api := app.Group("/api/book")
	api.Post("/", middlewares.AuthMiddleware, controllers.CreateBook)
	api.Get("/", controllers.ListBooks)
	api.Put("/:id", middlewares.AuthMiddleware, controllers.UpdateBook)
	api.Delete("/:id", middlewares.AuthMiddleware, controllers.DeleteBook)
}