package tgbotbase

type PropertyValue struct {
	Value string
	User  UserID
	Chat  ChatID
}

type PropertyStorage interface {
	GetProperty(name string, user UserID, chat ChatID) (string, error)
	SetPropertyForUser(name string, user UserID, value interface{}) error
	SetPropertyForChat(name string, chat ChatID, value interface{}) error
	SetPropertyForUserInChat(name string, user UserID, chat ChatID, value interface{}) error
	GetEveryHavingProperty(name string) ([]PropertyValue, error)
}
