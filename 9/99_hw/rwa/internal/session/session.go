package session

import (
	"crypto/rand"
	"encoding/hex"
	"rwa/internal/db"
	"sync"
)

type Manager struct {
	Sessions map[string]db.Session
	Mu       *sync.Mutex
}

func generateSessionID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (m *Manager) Create(user *db.User) (string, error) {
	sessionId, err := generateSessionID()
	if err != nil {
		return "", err
	}

	m.Destroy(user.Token)

	m.Mu.Lock()
	m.Sessions[sessionId] = db.Session{User: user}
	m.Mu.Unlock()
	user.Token = sessionId
	return sessionId, nil
}

func (m *Manager) Destroy(sessionId string) {
	m.Mu.Lock()
	delete(m.Sessions, sessionId)
	m.Mu.Unlock()
}
