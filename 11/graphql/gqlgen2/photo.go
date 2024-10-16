package gqlgen2

import (
	"log"
	"strconv"
)

type Photo struct {
	ID     uint `json:"id"`
	UserID uint `json:"-"`
	// User     *User  `json:"user"`
	URL      string `json:"url"`
	Comment  string `json:"comment"`
	Rating   int    `json:"rating"`
	Liked    bool   `json:"liked"`
	Followed bool   `json:"followed"`
}

func (ph *Photo) Id() string {
	log.Println("call Photo.ID method", ph.ID)
	return strconv.Itoa(int(ph.ID))
}
