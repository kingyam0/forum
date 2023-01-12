package handlers

import (
	"html/template"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Need to create structure to show array of Users, Posts, Comments, Categories for arranging them in HTML
	type PostRaw struct {
		Post   *Pitem
		PostLD *Plikedislike
	}

	var mainPage struct {
		AuthUser   *User
		PostScroll []*PostRaw
	}
	var post *Pitem

	tpl, err := template.ParseGlob("static/templates/home.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error, Parsing Error", http.StatusInternalServerError)
		WarnMessage(w, "500 Internal Server Error, Parsing Error")
		return
	}
	if r.URL.Path != "/" {
		WarnMessage(w, "404 Not Found")
		return
	}

	var AllPosts []*Pitem = ReadAllPosts()

	if user_session.Name == "session_token" {
		c, err := r.Cookie("session_token")
		if err != nil {
			WarnMessage(w, "400 Bad Request, Session expired")
			return
		}
		// if User is authenticated
		user := GetUserByCookie(c.Value)

		mainPage.AuthUser = user
		for _, post = range AllPosts {
			auth := &PostRaw{}
			auth.Post = post
			auth.PostLD = GetLikeDislikeCountOfPost(post.PostID)
			mainPage.PostScroll = append(mainPage.PostScroll, auth)
		}
		tpl.ExecuteTemplate(w, "home.html", mainPage)
	} else {
		for _, post = range AllPosts {
			guest := PostRaw{}
			guest.Post = post
			guest.PostLD = GetLikeDislikeCountOfPost(post.PostID)
			mainPage.PostScroll = append(mainPage.PostScroll, &guest)
		}
		tpl.Execute(w, mainPage)
	}
}

// WarnMessage ...
func WarnMessage(w http.ResponseWriter, warn string) {
	tpl, err := template.ParseFiles("static/templates/error.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error, Parsing Error", http.StatusInternalServerError)
		return
	}
	Warning.Warn = warn
	tpl.Execute(w, Warning)
}
