package internal

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

type SessionManager struct {
	Sessions map[string]*Session
	Mu       *sync.Mutex
}

func CreateSessionManager() *SessionManager {
	return &SessionManager{
		Sessions: make(map[string]*Session),
		Mu:       &sync.Mutex{},
	}
}

func generateSessionID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (m *SessionManager) Create(user *User) (string, error) {
	sessionId, err := generateSessionID()
	if err != nil {
		return "", err
	}

	m.Destroy(user.Token)

	m.Mu.Lock()
	m.Sessions[sessionId] = &Session{User: user}
	m.Mu.Unlock()
	user.Token = sessionId
	return sessionId, nil
}

func (m *SessionManager) Destroy(sessionId string) {
	m.Mu.Lock()
	delete(m.Sessions, sessionId)
	m.Mu.Unlock()
}
