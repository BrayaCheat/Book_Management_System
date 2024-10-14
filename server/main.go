package main

import (
	"log"
	"os"
	"os/signal"
	"server/src/configs"
	"server/src/routes"
	"syscall"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Initialize the database
	if err := configs.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// Create a new Fiber instance
	app := fiber.New()
	
	// Set up routes
	app.Use(cors.New())
	routes.SetupRoutes(app)

	app.Options("/*", func (c *fiber.Ctx) error  {
		return c.Status(200).JSON(
			fiber.Map{
				"message": "Enable preflight",
			},
		)
	})

	// Graceful shutdown
	go func() {
		if err := LoadPort(app); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Listen for OS signals to gracefully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	log.Println("Gracefully shutting down...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Error during server shutdown: %v", err)
	}

	log.Println("Server shutdown complete.")
}

func LoadPort(app *fiber.App) error {
	port := os.Getenv("PORT")
	err := app.Listen("0.0.0.0:" + port)
	if err != nil {
		log.Fatal("Port issue")
	}

	return nil
}
