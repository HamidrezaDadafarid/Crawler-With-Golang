package models

type UserSession struct {
	Filters Filters
	Email   string
	State   string
	ChatID  int64
}

var userSessions = make(map[int64]*UserSession)

func GetUserSession(chatID int64) *UserSession {
	session, exists := userSessions[chatID]
	if !exists {
		session = &UserSession{
			ChatID: chatID,
		}
		userSessions[chatID] = session
	}
	return session
}
