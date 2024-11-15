package models

type UserSession struct {
	Filters         Filters
	Email           string
	State           string
	ChatID          int64
	IsAuthenticated bool
}

var userSessions = make(map[int64]*UserSession)

func GetUserSession(chatID int64) *UserSession {
	session, exists := userSessions[chatID]
	if !exists {
		session = &UserSession{
			Filters:         Filters{},
			Email:           "",
			State:           "",
			ChatID:          chatID,
			IsAuthenticated: false,
		}
		userSessions[chatID] = session
	}
	return session
}
