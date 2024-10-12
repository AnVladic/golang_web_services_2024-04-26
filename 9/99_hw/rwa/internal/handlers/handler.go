package handlers

import (
	"encoding/json"
	"net/http"
	"rwa/internal/session"
	"rwa/internal/utils"
)

type RwaHandler struct {
	SessionManager session.Manager
}

func JsonWrite(w http.ResponseWriter, status int, response interface{}) {
	responseBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SetJsonContentHeader(w)
	w.WriteHeader(status)
	_, err = w.Write(responseBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
