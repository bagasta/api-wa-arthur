package main

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Agent Model mirror (dummy for relation)
type Agent struct {
	ID   string `gorm:"primaryKey;type:uuid"`
	Name string
}

// ApiKey Model mirror
type ApiKey struct {
	ID          string `gorm:"primaryKey;type:uuid"`
	AgentID     string `gorm:"type:uuid"`
	AccessToken string `gorm:"type:varchar"`
	IsActive    bool
}

// WhatsappSession stores user's current connection to an agent
type WhatsappSession struct {
	ChatID        string  `gorm:"primaryKey;type:varchar"`
	ActiveAgentID *string `gorm:"type:uuid"`
	ApiKey        *string `gorm:"type:varchar(512)"`
	ConnectedAt   time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

var DB *gorm.DB

func InitDB(dsn string) {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Fatal: Failed to connect to DB: %v", err)
	}

	// Auto Migrate exclusively for the WhatsappSession table
	if err := DB.AutoMigrate(&WhatsappSession{}); err != nil {
		log.Printf("Warn: AutoMigrate failed: %v", err)
	}

	log.Println("Database connection established.")
}
