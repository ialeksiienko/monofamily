package sessions

import (
	"main-service/internal/entities"
	"sync"
)

type UserState struct {
	Family *entities.Family
}

var (
	usersState = make(map[int64]*UserState)
	mu         sync.RWMutex
)

func SetUserState(userID int64, uState *UserState) {
	mu.Lock()
	defer mu.Unlock()

	usersState[userID] = uState
}

func GetUserState(userID int64) (*UserState, bool) {
	mu.RLock()
	defer mu.RUnlock()

	s, exists := usersState[userID]
	if !exists {
		return nil, false
	}
	return s, true
}

func DeleteUserState(userID int64) {
	mu.Lock()
	defer mu.Unlock()
	delete(usersState, userID)
}
