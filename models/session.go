package models

type UserSession struct {
	ChatID  int64
	Filters Filters
	State   string
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
