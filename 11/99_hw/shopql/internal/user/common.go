package user

import "golang.org/x/crypto/argon2"

func (h *Handler) hashPass(plainPassword, salt string) []byte {
	hashedPass := argon2.IDKey(
		[]byte(plainPassword), []byte(salt), 1, 64*1024, 4, 32)
	res := make([]byte, len(salt))
	copy(res, salt[:])
	return append(res, hashedPass...)
}
