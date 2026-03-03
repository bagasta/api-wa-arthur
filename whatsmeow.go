package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var WAClient *whatsmeow.Client

func ConnectWhatsapp() {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New(context.Background(), "sqlite3", "file:store.db?_foreign_keys=on", dbLog)
	if err != nil {
		log.Fatalf("Fatal: failed to open whatsapp store: %v", err)
	}

	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		log.Fatalf("Fatal: failed to get device: %v", err)
	}

	clientLog := waLog.Stdout("Client", "INFO", true)
	WAClient = whatsmeow.NewClient(deviceStore, clientLog)
	WAClient.AddEventHandler(eventHandler)

	if WAClient.Store.ID == nil {
		qrChan, _ := WAClient.GetQRChannel(context.Background())
		err = WAClient.Connect()
		if err != nil {
			log.Fatalf("Fatal: failed to connect: %v", err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				log.Println("====================")
				log.Println("SCAN QR CODE THIS IN WHATSAPP:")
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				log.Println("====================")
			} else {
				log.Println("QR channel event:", evt.Event)
			}
		}
	} else {
		err = WAClient.Connect()
		if err != nil {
			log.Fatalf("Fatal: failed to connect: %v", err)
		}
		log.Println("Whatsapp already logged in.")
	}
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if v.Info.IsFromMe {
			return
		}

		// Only respond to personal/private chats, ignore group messages
		if v.Info.Chat.Server == "g.us" {
			return
		}
		// Handle only text messages
		var text string
		if v.Message.GetConversation() != "" {
			text = v.Message.GetConversation()
		} else if v.Message.GetExtendedTextMessage() != nil {
			text = v.Message.GetExtendedTextMessage().GetText()
		}

		if text == "" {
			return
		}

		chatID := v.Info.Chat.String()
		senderName := v.Info.PushName
		if senderName == "" {
			senderName = "User"
		}

		log.Printf("Received WA Message from %s: %s", chatID, text)

		webhookURL := os.Getenv("N8N_WEBHOOK_URL")
		if webhookURL == "" {
			log.Println("WARN: N8N_WEBHOOK_URL is not set. Directly processing internally...")
			// Dummy wrap the payload so we can manually call our internal handler if not using N8N
			go simulateN8nWebhook(chatID, senderName, text)
			return
		}

		// Forward to N8N webhook trigger
		payload := map[string]interface{}{
			"chat_id":    chatID,
			"first_name": senderName,
			"text":       text,
		}
		jsonData, _ := json.Marshal(payload)
		resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("ERROR: Failed to forward to N8N webhook: %v", err)
		} else {
			resp.Body.Close()
		}
	}
}

// simulateN8nWebhook allows the app to bypass N8N entirely if they haven't configured the Webhook Trigger yet.
func simulateN8nWebhook(chatID, firstName, text string) {
	payload := map[string]interface{}{
		"message": map[string]interface{}{
			"chat": map[string]interface{}{"id": chatID},
			"from": map[string]interface{}{"id": chatID, "first_name": firstName},
			"text": text,
		},
	}
	jsonData, _ := json.Marshal(payload)
	// POST to our own endpoint
	http.Post("http://127.0.0.0:8101/api/v1/whatsapp/webhook", "application/json", bytes.NewBuffer(jsonData))
}
