package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

/* Filter Handle
1) created posts --> /filter?section=my_posts
2) liked posts   --> /filter?section=rated
3) by categories --> /filter?section=<request>
*/

func FilterHandle(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseGlob("static/templates/filter.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error, Parsing Error", http.StatusInternalServerError)
		WarnMessage(w, "500 Internal Server Error, Parsing Error")
		return
	}

	type PostRaw struct {
		Post   *Pitem
		PostLD *Plikedislike
	}
	var filter struct {
		Section    string
		AuthUser   *User
		PostScroll []*PostRaw
	}

	if r.URL.Path != "/filter" {
		WarnMessage(w, "404 Not Found")
		return
	}

	// if user authenticated
	section := r.URL.Query().Get("section")
	if section == "my_posts" {
		c, err := r.Cookie("session_token")
		if err != nil {
			WarnMessage(w, "400 Bad Request, You need to be authorised")
			return
		}
		user := GetUserByCookie(c.Value)
		// 1) search posts by userID
		posts := GetPostByUID(user.UserID)

		for _, post := range posts {
			result := &PostRaw{}
			result.Post = post
			result.PostLD = GetLikeDislikeCountOfPost(post.PostID)
			filter.PostScroll = append(filter.PostScroll, result)
		}
		filter.Section = "My Posts"
		filter.AuthUser = user
		tpl.Execute(w, filter)
		filter.Section = ""
		filter.AuthUser = nil
		filter.PostScroll = nil

	} else if section == "liked" {

		//  2) SEARCH BY LIKED POSTS-->else if section == liked
		c, err := r.Cookie("session_token")
		if err != nil {
			WarnMessage(w, "400 Bad Request, You need to be authorized")
			return
		}
		user := GetUserByCookie(c.Value)

		posts := GetRatedPostsByUID(user.UserID)

		for _, post := range posts {
			result := &PostRaw{}
			result.Post = post
			result.PostLD = GetLikeDislikeCountOfPost(post.PostID)
			filter.PostScroll = append(filter.PostScroll, result)

		}
		filter.Section = "Liked posts"
		filter.AuthUser = user

		tpl.ExecuteTemplate(w, "filter.html", filter)
		// filter.Section = ""
		// filter.AuthUser = nil
		filter.PostScroll = nil

		return
	}

	//  3) search by category
	if section == "travel" || section == "currentaffairs" || section == "sports" || section == "hobby" {

		c, err := r.Cookie("session_token")
		if err == nil {
			user := GetUserByCookie(c.Value)

			filter.AuthUser = user
			fmt.Println(section)
		}
		cat1 := r.URL.Query().Get("section")
		cat2 := r.URL.Query().Get("section")
		cat3 := r.URL.Query().Get("section")
		cat4 := r.URL.Query().Get("section")
		fmt.Println(cat1)
		fmt.Println(cat2)
		fmt.Println(cat4)

		posts := GetPostByCategory(cat1, cat2, cat3, cat4)
		fmt.Println("all categories: ", cat1, cat2, cat3, cat4)
		for _, post := range posts {
			result := &PostRaw{}
			result.Post = post
			result.PostLD = GetLikeDislikeCountOfPost(post.PostID)
			filter.PostScroll = append(filter.PostScroll, result)
		}
		filter.Section = ""
		tpl.Execute(w, filter)
		filter.AuthUser = nil
		filter.PostScroll = nil
	}
}
