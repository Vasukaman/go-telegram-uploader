package bot

type MessageStore interface {
	SaveMessage(message string)
	GetMessage() (messaage string)
}

type InMemoryMessageStore struct {
	message string
}

func NewInMemoryMessageStore() *InMemoryMessageStore {
	return &InMemoryMessageStore{message: "No secrets in here."}
}

func (s *InMemoryMessageStore) SaveMessage(message string) {
	s.message = message
}

// There should be a error output here. But for now it's only 1 string.
// I do not delete the message, bc the idea is that it needs to be rewritten by user.
func (s *InMemoryMessageStore) GetMessage() string {

	message := s.message

	return message
}
