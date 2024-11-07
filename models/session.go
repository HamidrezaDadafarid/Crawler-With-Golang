package models

type UserSession struct {
	ChatID  int64
	Filters map[string]string
	State   string
}

var userSessions = make(map[int64]*UserSession)

func GetUserSession(chatID int64) *UserSession {
	session, exists := userSessions[chatID]
	if !exists {
		session = &UserSession{
			ChatID:  chatID,
			Filters: make(map[string]string),
		}
		userSessions[chatID] = session
	}
	return session
}
