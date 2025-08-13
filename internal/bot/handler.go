package bot

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MessageHandler handles the logic for incoming messages.
type MessageHandler struct {
	messenger Messenger
	uploader  ImageUploader
}

// NewMessageHandler creates a new handler.
func NewMessageHandler(messenger Messenger, uploader ImageUploader) *MessageHandler {
	return &MessageHandler{messenger: messenger, uploader: uploader}
}

// Handle takes an incoming message and routes it to the correct logic.
func (h *MessageHandler) Handle(message *tgbotapi.Message) {
	// Log everything for debugging
	log.Printf("Received: [%s] Caption: [%s] Text: [%s]", message.From.UserName, message.Caption, message.Text)

	// --- Step 1: Check for any media ---
	if message.Photo != nil || message.Animation != nil || message.Document != nil {
		h.handleMediaUpload(message)
		return // We're done
	}

	// --- Step 2: If no media, handle the /test command or give help ---
	if strings.HasPrefix(message.Text, "/test") {
		h.handleTestCommand(message)
	} else {
		h.sendMessage(tgbotapi.NewMessage(message.Chat.ID, "Hello! Send me a photo, GIF, or document, and I will upload it for you."))
	}
}

// handleMediaUpload figures out what kind of media was sent and gets the FileID.
func (h *MessageHandler) handleMediaUpload(message *tgbotapi.Message) {
	var fileID string

	if message.Photo != nil {
		fileID = message.Photo[len(message.Photo)-1].FileID
	} else if message.Animation != nil {
		fileID = message.Animation.FileID
	} else if message.Document != nil {
		// You can add more checks here, like for file type or size
		fileID = message.Document.FileID
	} else {
		h.sendMessage(tgbotapi.NewMessage(message.Chat.ID, "Sorry, I couldn't find a file to upload."))
		return
	}

	// Now that we have the FileID, we can start the upload process.
	h.processFileUpload(message.Chat.ID, fileID)
}

// processFileUpload contains the actual download/upload logic.
func (h *MessageHandler) processFileUpload(chatID int64, fileID string) {
	h.sendMessage(tgbotapi.NewMessage(chatID, "Media received! Uploading..."))

	go func() {
		fileURL, err := h.messenger.GetFileDirectURL(fileID)
		if err != nil {
			h.sendMessage(tgbotapi.NewMessage(chatID, "Sorry, I couldn't get the file from Telegram."))
			return
		}

		resp, err := http.Get(fileURL)
		if err != nil {
			h.sendMessage(tgbotapi.NewMessage(chatID, "Sorry, I couldn't download the file."))
			return
		}
		defer resp.Body.Close()

		imageData, err := io.ReadAll(resp.Body)
		if err != nil {
			h.sendMessage(tgbotapi.NewMessage(chatID, "Sorry, an error occurred while processing the file."))
			return
		}

		uploadLink, err := h.uploader.UploadImage(imageData)
		if err != nil {
			h.sendMessage(tgbotapi.NewMessage(chatID, "Sorry, the upload to the cloud failed."))
			log.Printf("Error uploading to cloud: %v", err)
			return
		}

		h.sendMessage(tgbotapi.NewMessage(chatID, "Upload complete! Link: "+uploadLink))
	}()
}

// handleTestCommand processes the "/test" command.
func (h *MessageHandler) handleTestCommand(message *tgbotapi.Message) {
	randomNumber := rand.Intn(100) // Generate a random number between 0 and 99
	replyText := fmt.Sprintf("Test successful! Your random number is: %d", randomNumber)

	reply := tgbotapi.NewMessage(message.Chat.ID, replyText)
	h.sendMessage(reply)
}

// sendMessage is a helper function to send messages and log errors.
func (h *MessageHandler) sendMessage(msg tgbotapi.Chattable) {
	if _, err := h.messenger.Send(msg); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}
