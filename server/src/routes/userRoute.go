package routes

import (
	"server/src/controllers"
	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App){
	api := app.Group("/api/auth")
	api.Post("/login", controllers.Login)
	api.Post("/register", controllers.Register)
}