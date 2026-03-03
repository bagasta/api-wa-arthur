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

	// DATA_DIR allows store.db to be written to a Docker volume (/app/data)
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "." // default: current directory
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8101"
	}

	InitDB(dbUrl)
	ConnectWhatsapp(dataDir)

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	SetupRoutes(app)

	log.Printf("WhatsApp Endpoint starting on :%s", port)
	log.Fatal(app.Listen(":" + port))
}
