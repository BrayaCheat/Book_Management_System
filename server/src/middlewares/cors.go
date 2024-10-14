package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func EnableCors(app *fiber.App) error {
	app.Use(cors.New(
		cors.Config{
			AllowOrigins: "http://localhost:8081, http://localhost:8080",
        	AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
       	 	AllowHeaders: "Origin, Content-Type, Accept",
        	AllowCredentials: true,
		},
	))
	return nil
}