package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("WARN: Error loading .env file, relying on system env variables")
	}

	LangchainURL = os.Getenv("LANGCHAIN_API_URL")
	if LangchainURL == "" {
		LangchainURL = "http://localhost:3000" // default fallback
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatalf("Fatal: DATABASE_URL is not set")
	}

	InitDB(dbUrl)
	ConnectWhatsapp()

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	SetupRoutes(app)

	log.Println("WhatsApp Endpoint starting on :8101")
	log.Fatal(app.Listen(":8101"))
}
