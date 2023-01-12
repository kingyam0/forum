package handlers

import (
	"html/template"
	"time"
)

var tpl *template.Template

type User struct {
	UserID       int
	Email        string
	Username     string
	passwordhash string
	CreationDate time.Time
}

var user_session Cookie

var CurrentUser User

var Warning struct {
	Warn string
}

// each session contains the username of the user and the time at which it expires
type Session struct {
	UserID      int
	username    string
	sessionName string
	sessionUUID string
	expiry      time.Time
}

type Cookie struct {
	Name    string
	Value   string
	Expires time.Time
}

type Pitem struct {
	PostID       int
	AuthorID     string
	Author       string
	Title        string
	Text         string
	Category1    string
	Category2    string
	Category3    string
	Category4    string
	Likes        int
	Dislikes     int
	CreationDate time.Time
	Comments     []Comm
}

type Category struct {
	CategoryID    int
	Catergoryname string
	PostID        Pitem
}

type Comm struct {
	CommentID    int
	PostID       int
	AuthorID     int
	Author       string
	Text         string
	Likes        int
	Dislikes     int
	CreationDate time.Time
}

type Plikedislike struct {
	PostID       int
	UserID       int
	Like         bool
	Likecount    int
	Dislikecount int
}

type Clikedislike struct {
	CommentID    int
	PostID       int
	UserID       int
	Like         bool
	Likecount    int
	Dislikecount int
}
