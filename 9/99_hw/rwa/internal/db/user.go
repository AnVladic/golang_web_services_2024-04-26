package db

import "time"

type User struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	Username  string     `json:"username,omitempty"`
	Bio       string     `json:"bio,omitempty"`
	Image     string     `json:"image,omitempty"`
	Token     string     `json:"-"`
	Following bool

	Password []byte `json:"-"`
	Salt     string `json:"-"`
}
