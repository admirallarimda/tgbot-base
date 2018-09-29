package botbase

type PropertyStorage interface {
	GetProperty(name string, user UserID, chat ChatID) (string, error)
	SetPropertyForUser(name string, user UserID, value interface{}) error
	SetPropertyForChat(name string, chat ChatID, value interface{}) error
	SetPropertyForUserInChat(name string, user UserID, chat ChatID, value interface{}) error
}
