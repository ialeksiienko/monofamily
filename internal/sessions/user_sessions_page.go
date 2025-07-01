package sessions

import (
	"main-service/internal/entities"
	"sync"
)

type UPSessions struct {
	Page     int
	Families []entities.Family
}

var (
	userPageSessions      = make(map[int64]*UPSessions)
	userPageSessionsMutex sync.RWMutex
)

func SetUserPageSession(userID int64, session *UPSessions) {
	userPageSessionsMutex.Lock()
	defer userPageSessionsMutex.Unlock()

	userPageSessions[userID] = session
}

func GetUserPageSession(userID int64) (*UPSessions, bool) {
	userPageSessionsMutex.RLock()
	defer userPageSessionsMutex.RUnlock()

	session, exists := userPageSessions[userID]
	if !exists {
		return nil, false
	}
	return session, true
}
