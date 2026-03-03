package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

// Langchain API URL
var LangchainURL string

type SessionResponse struct {
	ChatID        string  `json:"chat_id"`
	ActiveAgentID *string `json:"active_agent_id"`
}

func SetupRoutes(app *fiber.App) {
	app.Get("/api/v1/whatsapp/sessions/:id", getSession)
	app.Post("/api/v1/whatsapp/send", sendWhatsAppMessageRaw)
	app.Post("/api/v1/whatsapp/webhook", processWebhook)
}

func getSession(c *fiber.Ctx) error {
	chatID := c.Params("id")
	var session WhatsappSession
	err := DB.Where("chat_id = ?", chatID).First(&session).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(200).JSON(fiber.Map{"error": "Session not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "DB error"})
	}

	return c.Status(200).JSON(SessionResponse{
		ChatID:        session.ChatID,
		ActiveAgentID: session.ActiveAgentID,
	})
}

func sendWhatsAppMessageRaw(c *fiber.Ctx) error {
	var body struct {
		ChatID string `json:"chat_id"`
		Text   string `json:"text"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid format"})
	}
	err := sendWAMessage(body.ChatID, body.Text)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"success": true})
}

// simulate the n8n webhook parsing
type WebhookBody struct {
	Message struct {
		Chat struct {
			ID string `json:"id"`
		} `json:"chat"`
		From struct {
			ID        string `json:"id"`
			FirstName string `json:"first_name"`
		} `json:"from"`
		Text string `json:"text"`
	} `json:"message"`
}

func processWebhook(c *fiber.Ctx) error {
	var body WebhookBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).SendString("Invalid format")
	}

	text := strings.TrimSpace(body.Message.Text)
	if text == "" {
		return c.Status(200).SendString("OK")
	}

	chatID := body.Message.Chat.ID

	if strings.HasPrefix(strings.ToLower(text), "/start") {
		msg := "👋 *Selamat datang di AI Agent Bot untuk WhatsApp!*\n\n" +
			"Bot ini memungkinkan kamu berbicara langsung dengan AI Agent.\n\n" +
			"*Perintah yang tersedia:*\n" +
			"• `/connect <AGENT_ID>` – Hubungkan ke agent\n" +
			"• `/disconnect` – Putus koneksi agent\n" +
			"• `/status` – Cek agent yang aktif\n\n" +
			"Kirimkan Agent ID saja untuk terhubung."
		go sendWAMessage(chatID, msg)
		return c.Status(200).SendString("OK")
	}

	if strings.HasPrefix(strings.ToLower(text), "/connect") {
		go sendWAMessage(chatID, "⏳ Memeriksa Agent dan API Key di Database...")
		parts := strings.SplitN(text, " ", 3)
		if len(parts) < 2 {
			go sendWAMessage(chatID, "⚠️ Format salah. Gunakan:\n`/connect <AGENT_ID>`\n\nContoh:\n`/connect 3fa85f64-5717-4562-b3fc-2c963f66afa6`")
			return c.Status(200).SendString("OK")
		}

		rawID := strings.TrimSpace(parts[1])
		var apiKey string

		if len(rawID) < 10 {
			go sendWAMessage(chatID, fmt.Sprintf("❌ Agent ID tidak valid: %s", rawID))
			return c.Status(200).SendString("OK")
		}

		// Pull API Key from DB
		var apiKeyRecord ApiKey
		err := DB.Where("agent_id = ? AND is_active = ?", rawID, true).First(&apiKeyRecord).Error
		if err != nil {
			go sendWAMessage(chatID, fmt.Sprintf("❌ Tidak menemukan API Key di database untuk Agent %s", rawID))
			return c.Status(200).SendString("OK")
		}
		apiKey = apiKeyRecord.AccessToken

		// Verify Agent exists in Langchain API
		url := fmt.Sprintf("%s/api/v1/agents/%s", LangchainURL, rawID)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+apiKey)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			go sendWAMessage(chatID, "❌ Gagal memvalidasi agent. Layanan down.")
			return c.Status(200).SendString("OK")
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			go sendWAMessage(chatID, fmt.Sprintf("❌ Validasi gagal (HTTP %d). Agent mungkin tidak ada atau API Key expired.", resp.StatusCode))
			return c.Status(200).SendString("OK")
		}

		// Upsert DB
		var session WhatsappSession
		result := DB.Where("chat_id = ?", chatID).First(&session)

		idStr := rawID
		if result.Error == gorm.ErrRecordNotFound {
			session = WhatsappSession{
				ChatID:        chatID,
				ActiveAgentID: &idStr,
				ApiKey:        &apiKey,
				ConnectedAt:   time.Now(),
			}
			DB.Create(&session)
		} else {
			session.ActiveAgentID = &idStr
			session.ApiKey = &apiKey
			session.ConnectedAt = time.Now()
			DB.Save(&session)
		}

		go sendWAMessage(chatID, "✅ *Terhubung dengan Agent!* Hubungi agent sekarang dengan mengetik sapaan.")
		return c.Status(200).SendString("OK")
	}

	if strings.HasPrefix(strings.ToLower(text), "/disconnect") {
		var session WhatsappSession
		err := DB.Where("chat_id = ?", chatID).First(&session).Error
		if err == nil {
			session.ActiveAgentID = nil
			session.ApiKey = nil
			DB.Save(&session)
			go sendWAMessage(chatID, "🔌 Koneksi dengan agent telah diputus.\nGunakan `/connect <AGENT_ID>` untuk terhubung kembali.")
		} else {
			go sendWAMessage(chatID, "ℹ️ Tidak ada agent yang terhubung saat ini.")
		}
		return c.Status(200).SendString("OK")
	}

	if strings.HasPrefix(strings.ToLower(text), "/status") {
		var session WhatsappSession
		err := DB.Where("chat_id = ?", chatID).First(&session).Error
		if err == nil && session.ActiveAgentID != nil {
			go sendWAMessage(chatID, fmt.Sprintf("🔗 Saat ini kamu terhubung dengan Agent ID:\n`%s`", *session.ActiveAgentID))
		} else {
			go sendWAMessage(chatID, "ℹ️ Belum ada agent yang terhubung.\nGunakan `/connect <AGENT_ID>` untuk mulai.")
		}
		return c.Status(200).SendString("OK")
	}

	// Normal chat
	var session WhatsappSession
	err := DB.Where("chat_id = ?", chatID).First(&session).Error
	if err != nil || session.ActiveAgentID == nil || session.ApiKey == nil {
		go sendWAMessage(chatID, "⚠️ Belum terhubung ke agent. Gunakan:\n`/connect <AGENT_ID>`")
		return c.Status(200).SendString("OK")
	}

	sessionStr := fmt.Sprintf("wa_%s", chatID)
	go func() {
		// --- Typing Indicator Loop ---
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			jid, _ := types.ParseJID(chatID)
			for {
				select {
				case <-ctx.Done():
					if WAClient != nil && WAClient.IsConnected() {
						WAClient.SendChatPresence(context.Background(), jid, types.ChatPresencePaused, types.ChatPresenceMediaText)
					}
					return
				default:
					if WAClient != nil && WAClient.IsConnected() {
						WAClient.SendChatPresence(context.Background(), jid, types.ChatPresenceComposing, types.ChatPresenceMediaText)
					}
					time.Sleep(5 * time.Second)
				}
			}
		}()
		// --- End Typing Indicator Loop ---

		reply, active := callLangchainExecute(*session.ActiveAgentID, *session.ApiKey, text, sessionStr)

		cancel() // Stop typing indicator

		if !active {
			session.ActiveAgentID = nil
			session.ApiKey = nil
			DB.Save(&session)
			sendWAMessage(chatID, "❌ Agent tidak lagi tersedia. Gunakan `/connect` untuk menghubungkan ulang.")
			return
		}
		if reply != "" {
			sendWAMessage(chatID, reply)
		}
	}()

	return c.Status(200).SendString("OK")
}

func callLangchainExecute(agentID, apiKey, inputText, sessionID string) (string, bool) {
	url := fmt.Sprintf("%s/api/v1/agents/%s/execute", LangchainURL, agentID)

	payload := map[string]interface{}{
		"input":      inputText,
		"session_id": sessionID,
		"parameters": map[string]interface{}{},
	}
	jsonData, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("ERROR calling langchain: %v", err)
		return "⏱️ Waktu tunggu habis atau terjadi kesalahan. Coba lagi nanti.", true
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return "⛔ Token limit agent ini telah habis. Hubungi admin untuk menambah kuota.", true
	}
	if resp.StatusCode == 404 {
		return "", false
	}
	if resp.StatusCode == 401 {
		return "🔑 API key tidak valid atau sudah expired.\nCoba reconnect kembali dengan `/connect <AGENT_ID>`", true
	}
	if resp.StatusCode != 200 {
		return fmt.Sprintf("❌ Error dari AI service (%d). Coba lagi nanti.", resp.StatusCode), true
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if val, ok := result["response"].(string); ok {
		return val, true
	}
	if val, ok := result["output"].(string); ok {
		return val, true
	}
	return "🤔 Agent tidak memberikan respons.", true
}

func sendWAMessage(chatID string, text string) error {
	if WAClient == nil || !WAClient.IsConnected() {
		return fmt.Errorf("WAClient is not connected")
	}
	jid, err := types.ParseJID(chatID)
	if err != nil {
		return err
	}

	msg := &waE2E.Message{
		Conversation: proto.String(text),
	}

	_, err = WAClient.SendMessage(context.Background(), jid, msg)
	return err
}
