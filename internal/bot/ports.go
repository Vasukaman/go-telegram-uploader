// in file: /internal/bot/ports.go
package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Messenger is an interface that describes sending messages and getting file URLs.
// This is all our handler needs to know.
type Messenger interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	GetFileDirectURL(fileID string) (string, error)
}

type MessageStore interface {
	SaveMessage(message string)
	GetMessage() (messaage string)
}

type ImageUploader interface {
	UploadImage(imageData []byte) (string, error)
}
