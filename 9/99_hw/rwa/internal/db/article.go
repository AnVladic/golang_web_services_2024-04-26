package db

import "time"

type Article struct {
	Author         User      `json:"author"`
	Body           string    `json:"body"`
	CreatedAt      time.Time `json:"createdAt"`
	Description    string    `json:"description"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount int       `json:"favoritesCount"`
	Slug           string    `json:"slug"`
	TagList        []string  `json:"tagList"`
	Title          string    `json:"title"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
