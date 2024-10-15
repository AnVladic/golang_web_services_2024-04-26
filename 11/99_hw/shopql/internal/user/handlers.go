package user

import (
	"encoding/json"
	"hw11_shopql/internal"
	"net/http"
)

type Handler struct {
	SessionManager internal.SessionManager
}

func SetJsonContentHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func JsonWrite(w http.ResponseWriter, status int, response interface{}) {
	responseBytes, err := json.Marshal(map[string]interface{}{"body": response})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SetJsonContentHeader(w)
	w.WriteHeader(status)
	_, err = w.Write(responseBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	type RegisterUser struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var reqBody struct {
		User RegisterUser `json:"user"`
	}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newUser := internal.User{
		Email:    reqBody.User.Email,
		Password: h.hashPass(reqBody.User.Password, ""),
		Cart:     make([]*internal.CartItem, 0),
	}
	token, err := h.SessionManager.Create(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	response := map[string]interface{}{"token": token}
	JsonWrite(w, http.StatusOK, response)
}
