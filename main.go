package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Ryanprw/simple-upload-file/core/utils"
	"github.com/Ryanprw/simple-upload-file/module/file"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using default values")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	uploadDir := "./public/uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Error creating upload directory: %v", err)
		}
	}

	config := fiber.Config{
		BodyLimit: 1024 * 1024 * 1024,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err == nil {
				return nil
			}
			code := fiber.ErrInternalServerError.Code
			message := fiber.ErrInternalServerError.Message

			if value, ok := err.(*fiber.Error); ok {
				code = value.Code
				message = value.Message
			}
			return c.Status(code).JSON(utils.BaseResponse{
				Message: message,
				Error:   nil,
			})
		},
	}

	app := fiber.New(config)

	app.Use(cors.New())

	app.Static("/", "./public")
	app.Static("/uploads", uploadDir)

	v1 := app.Group("/api/v1")
	v1.Route("/file", file.Route)

	localIP := "192.168.1.2"
	serverAddr := "0.0.0.0:" + port

	// Start Server
	fmt.Printf("Server running at: http://127.0.0.1:%s\n", port)
	fmt.Printf("Accessible via local network: http://%s:%s\n", localIP, port)
	log.Fatal(app.Listen(serverAddr))
}
