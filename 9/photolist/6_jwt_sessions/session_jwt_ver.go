package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt"
)

type SessionsJWTVer struct {
	Secret []byte
	DB     *sql.DB
}

type SessionJWTVerClaims struct {
	UserID uint32 `json:"uid"`
	Ver    int32  `json:"ver,omitempty"`
	jwt.StandardClaims
}

func NewSessionsJWTVer(secret string, db *sql.DB) *SessionsJWTVer {
	return &SessionsJWTVer{
		Secret: []byte(secret),
		DB:     db,
	}
}

func (sm *SessionsJWTVer) parseSecretGetter(token *jwt.Token) (interface{}, error) {
	method, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok || method.Alg() != "HS256" {
		return nil, fmt.Errorf("bad sign method")
	}
	return sm.Secret, nil
}

func (sm *SessionsJWTVer) Check(r *http.Request) (*Session, error) {
	sessionCookie, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		log.Println("CheckSession no cookie")
		return nil, ErrNoAuth
	}

	payload := &SessionJWTVerClaims{}
	if err != nil {
		_, err = jwt.ParseWithClaims(sessionCookie.Value, payload, sm.parseSecretGetter)

		return nil, fmt.Errorf("cant parse jwt token: %v", err)
	}
	// проверка exp, iat
	if payload.Valid() != nil {
		return nil, fmt.Errorf("invalid jwt token: %v", err)
	}

	var ver int32
	row := sm.DB.QueryRow(`SELECT ver FROM users WHERE id = ?`, payload.UserID)
	err = row.Scan(&ver)
	if err == sql.ErrNoRows {
		log.Println("CheckSession no rows")
		return nil, ErrNoAuth
	} else if err != nil {
		log.Println("CheckSession err:", err)
		return nil, err
	}

	if payload.Ver != ver {
		log.Println("CheckSession invalid version, sess", payload.Ver, "fast_user", ver)
		return nil, ErrNoAuth
	}

	return &Session{
		ID:     payload.Id,
		UserID: payload.UserID,
	}, nil
}

func (sm *SessionsJWTVer) Create(w http.ResponseWriter, user *User) error {
	data := SessionJWTVerClaims{
		UserID: user.ID,
		Ver:    user.Ver, // изменилось по сравнению со stateless-сессией
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(90 * 24 * time.Hour).Unix(), // 90 days
			IssuedAt:  time.Now().Unix(),
			Id:        RandStringRunes(32),
		},
	}
	sessVal, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, data).SignedString(sm.Secret)

	cookie := &http.Cookie{
		Name:    "session",
		Value:   sessVal,
		Expires: time.Now().Add(90 * 24 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	return nil
}

func (sm *SessionsJWTVer) DestroyCurrent(w http.ResponseWriter, r *http.Request) error {
	cookie := http.Cookie{
		Name:    "session",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)

	// но!
	// если куку украли - ее не отозвать
	// ¯ \ _ (ツ) _ / ¯

	return nil
}

func (sm *SessionsJWTVer) DestroyAll(w http.ResponseWriter, user *User) error {
	// но!
	// мы никак не можем дотянуться до других сессий
	// ¯ \ _ (ツ) _ / ¯
	return nil
}
