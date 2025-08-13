package main

import (
	"log"

	"github.com/Vasukaman/go-telegram-uploader/internal/adapters/cloudinary"
	"github.com/Vasukaman/go-telegram-uploader/internal/bot" // Import our new bot package
	"github.com/Vasukaman/go-telegram-uploader/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

// TelegramAdapter makes our tgbotapi.BotAPI client implement our Messenger interface.
type TelegramAdapter struct {
	api *tgbotapi.BotAPI
}

// Send wraps the underlying tgbotapi Send method.
func (a *TelegramAdapter) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	return a.api.Send(c)
}

// GetFileDirectURL wraps the underlying tgbotapi method.
func (a *TelegramAdapter) GetFileDirectURL(fileID string) (string, error) {
	return a.api.GetFileDirectURL(fileID)
}

func main() {
	// 1. Load configuration from .env file.
	// This is the first thing we do.
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found, loading from system environment.")
	}

	// 2. Create a new config object which reads from the environment.
	cfg := config.New()

	// --- INITIALIZATION / WIRING ---

	// 1. Initialize the low-level Telegram client (the "driver").
	api, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Panicf("Failed to connect to Telegram: %v", err)
	}

	telegramAdapter := &TelegramAdapter{api: api}

	cloudinaryClient, err := cloudinary.NewClient(cfg.CloudinaryCloudName, cfg.CloudinaryAPIKey, cfg.CloudinaryAPISecret)
	if err != nil {
		log.Panicf("Failed to connect to Cloudinary: %v", err)
	}

	// 2. Create our message handler (the "business logic").
	//    We pass it the API client so it can send replies.
	messageHandler := bot.NewMessageHandler(telegramAdapter, cloudinaryClient)

	//imageUploadTest(cloudinaryClient)

	// 3. Create our main bot application (the "application core").
	//    We pass it the API client and our handler.
	telegramBot := bot.New(api, messageHandler, cloudinaryClient)

	// --- START THE APPLICATION ---
	telegramBot.Start()
}
