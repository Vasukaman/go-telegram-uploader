package bot

import (
	"log"

	"github.com/Vasukaman/go-telegram-uploader/internal/adapters/cloudinary"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot is the main application struct.
type Bot struct {
	api     *tgbotapi.BotAPI
	handler *MessageHandler
}

// New creates a new Bot instance.
func New(api *tgbotapi.BotAPI, handler *MessageHandler, cloudinaryClient *cloudinary.Client) *Bot {
	return &Bot{
		api:     api,
		handler: handler,
	}
}

// Start begins the bot's listening loop.
func (b *Bot) Start() {
	log.Printf("Authorized on account %s", b.api.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := b.api.GetUpdatesChan(updateConfig)

	// This is the main loop for receiving and handling updates.
	for update := range updates {
		// We only care about messages
		if update.Message == nil {
			continue
		}

		b.handler.Handle(update.Message)
	}
}
