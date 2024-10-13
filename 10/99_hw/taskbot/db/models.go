package db

type Chat struct {
	Id   int64
	User *User
}

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Chat *Chat  `json:"chat"`
}

type Task struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Author    *User  `json:"author"`
	Performer *User  `json:"performer"`
}
