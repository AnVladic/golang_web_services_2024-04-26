package fast_user

import (
	"golang.org/x/crypto/argon2"
	"rwa/internal/db"
	"rwa/internal/handlers"
	"sync"
)

type UserHandler struct {
	Users  map[string]*db.User
	UserMu sync.Mutex
	handlers.RwaHandler
}

func FindByEmail(users map[string]*db.User, email string) *db.User {
	for _, u := range users {
		if u.Email == email {
			return u
		}
	}
	return nil
}

func (h *UserHandler) hashPass(plainPassword, salt string) []byte {
	hashedPass := argon2.IDKey(
		[]byte(plainPassword), []byte(salt), 1, 64*1024, 4, 32)
	res := make([]byte, len(salt))
	copy(res, salt[:])
	return append(res, hashedPass...)
}
