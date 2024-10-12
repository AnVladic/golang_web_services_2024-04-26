package utils

import (
	"math/rand"
	"net/http"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func SetJsonContentHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func Contains(list []string, filter string) bool {
	for _, v := range list {
		if v == filter {
			return true
		}
	}
	return false
}
