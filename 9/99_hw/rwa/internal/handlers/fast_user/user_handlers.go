package fast_user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"rwa/internal/db"
	"rwa/internal/handlers"
	"rwa/internal/utils"
	"strconv"
	"time"
)

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUser struct {
	LoginUser
	Username string `json:"username"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		User RegisterUser `json:"user"`
	}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now()
	salt := utils.RandStringRunes(8)
	newUser := db.User{
		ID:        strconv.Itoa(len(h.Users)),
		Username:  reqBody.User.Username,
		Email:     reqBody.User.Email,
		CreatedAt: &now,
		UpdatedAt: &now,
		Password:  h.hashPass(reqBody.User.Password, salt),
		Salt:      salt,
	}
	h.UserMu.Lock()
	h.Users[newUser.ID] = &newUser
	h.UserMu.Unlock()

	token, err := h.SessionManager.Create(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var ur = struct {
		db.User
		Token string `json:"token"`
	}{
		User:  newUser,
		Token: token,
	}
	response := map[string]interface{}{"user": ur}
	handlers.JsonWrite(w, http.StatusCreated, response)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		User LoginUser `json:"user"`
	}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := FindByEmail(h.Users, reqBody.User.Email)
	if user.Email == "" {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if !bytes.Equal(user.Password, h.hashPass(reqBody.User.Password, user.Salt)) {
		http.Error(w, "Wrong password", http.StatusUnauthorized)
		return
	}

	token, err := h.SessionManager.Create(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var ur = struct {
		db.User
		Token string `json:"token"`
	}{
		User:  *user,
		Token: token,
	}
	response := map[string]interface{}{"user": ur}
	handlers.JsonWrite(w, http.StatusOK, response)
}

func (h *UserHandler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value("session").(*db.Session)
	if !ok {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	response := map[string]interface{}{"user": session.User}
	handlers.JsonWrite(w, http.StatusOK, response)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		User struct {
			Email string `json:"email"`
			Bio   string `json:"bio"`
		} `json:"user"`
	}
	session, ok := r.Context().Value("session").(*db.Session)
	if !ok {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.UserMu.Lock()
	defer h.UserMu.Unlock()

	delete(h.Users, session.User.Email)
	session.User.Email = reqBody.User.Email
	session.User.Bio = reqBody.User.Bio
	h.Users[session.User.ID] = session.User

	token, err := h.SessionManager.Create(session.User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var ur = struct {
		db.User
		Token string `json:"token"`
	}{
		User:  *session.User,
		Token: token,
	}
	response := map[string]interface{}{"user": ur}
	handlers.JsonWrite(w, http.StatusOK, response)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value("session").(*db.Session)
	if !ok {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	h.SessionManager.Destroy(session.User.Token)
	response := map[string]interface{}{"token": ""}
	handlers.JsonWrite(w, http.StatusOK, response)
}
