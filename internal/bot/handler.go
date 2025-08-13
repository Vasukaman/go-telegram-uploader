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

type ImageUploader interface {
	UploadImage(imageData []byte) (string, error)
}

// MessageHandler handles the logic for incoming messages.
type MessageHandler struct {
	messenger Messenger
	store     MessageStore
	uploader  ImageUploader
}

// NewMessageHandler creates a new handler.
func NewMessageHandler(messenger Messenger, store MessageStore, uploader ImageUploader) *MessageHandler {
	return &MessageHandler{messenger: messenger, store: store, uploader: uploader}
}

// Handle takes an incoming message and routes it to the correct logic.
func (h *MessageHandler) Handle(message *tgbotapi.Message) {
	log.Printf("Received: [%s] %s", message.From.UserName, message.Text)

	// Command routing
	switch {
	case message.Caption == "/upload": //If it's true, then we definetly have file attached AND command is correct. Mb not the cleanest way to do this, but i will keep it
		h.handleMediaUpload(message)
	case message.Text == "/upload":
		h.sendMessage(tgbotapi.NewMessage(message.Chat.ID, "To use this command, please send a photo or GIF with '/upload' as the caption."))
	case strings.HasPrefix(message.Text, "/test"):
		h.handleTestCommand(message)
	case strings.HasPrefix(message.Text, "/secret"):
		h.handleSecretCommand(message)
	default:
		h.handleUnknownCommand(message)
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

func (h *MessageHandler) handleSecretCommand(message *tgbotapi.Message) {

	newSecretMessage := strings.TrimSpace(strings.TrimPrefix(message.Text, "/secret"))
	hiddenMessage := h.store.GetMessage()

	replyText := fmt.Sprintf(`You unravel a secret message. It says:
	%s`, hiddenMessage)

	h.store.SaveMessage(newSecretMessage)
	reply := tgbotapi.NewMessage(message.Chat.ID, replyText)
	h.sendMessage(reply)
}

// handleUnknownCommand replies to any command that is not recognized.
func (h *MessageHandler) handleUnknownCommand(message *tgbotapi.Message) {
	reply := tgbotapi.NewMessage(message.Chat.ID, "Unknown command. Try /test.")
	h.sendMessage(reply)
}

// sendMessage is a helper function to send messages and log errors.
func (h *MessageHandler) sendMessage(msg tgbotapi.Chattable) {
	if _, err := h.messenger.Send(msg); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}
