package handlers

import (
	"fmt"
	"net/http"
	"strconv"
)

// LikeDislikeHandler ... /rate?post_id= || /rate?comment_id=
func LikeDislikeHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/rate" {
		WarnMessage(w, "404 Not Found")
		// http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	var (
		Islike = true
		isPost = true
	)

	cookie, err := r.Cookie("session_token")
	if err != nil {
		WarnMessage(w, "400 Bad Request, You need to be authorized")
		// http.Error(w, "You need to authorize", http.StatusForbidden)
		return
	} else if err == nil {

	id, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil {
		id, err = strconv.Atoi(r.URL.Query().Get("comment_id"))
		if err != nil {
			WarnMessage(w, "400 Bad Request, You are walking the wrong way")
			// http.Error(w, "Invalid parameter", http.StatusBadRequest)
			return
		} else{
		isPost = false
		}
	}
	// if "id" is negative number, it will be dislike
	if id < 0 {
		Islike = false
		id *= -1
	}
	user := GetUserByCookie(cookie.Value)

	if isPost {
		like := NewPostRating()
		like.Like = Islike
		err1 := GetPostByPID(int(id))
	
		fmt.Println("err1: ", err1)
		like.PostID = int(id)
		// fmt.Println("like.PostID: ", like.PostID)


		// like.UID = user.ID // assign like to UserID --> to check user liked this post,
		// prevent multiple liking of the post
		// Need to return new rate count
		/* Check user liked this post
		if Yes, delete rate from the post
		*/
		isRated := IsThisUsersPost(user.UserID, int(id))
		// fmt.Println("ISThisUsersPost - uid: ", user.UserID)
		// fmt.Println("ISThisUsersPost - pid:: ", int(id))
		// fmt.Println("IsRated: ", isRated)
		if isRated {
			if ok := UpdateLikeDislikeOfPost(like, user.UserID); !ok {
				fmt.Println("DeleteRateOfPost error")
				WarnMessage(w, "500 Internal Server Error, Something went wrong")
				// http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
		} else {
			if ok := AddLikeDislikeToPost(like, user.UserID); !ok {
				WarnMessage(w, "Something went wrong")
				// http.Error(w, "Something went wrong", http.StatusInternalServerError)
				fmt.Println("AddRatePost() error")
				return
			}
		}
	
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}  else {
		like := NewCommentRating()
		like.Like = Islike
		comment := GetCommentByCID(int(id))
		fmt.Println("CommentbyID: ", comment)
		if err != nil {
			WarnMessage(w, "The comment not found")
			// http.Error(w, "The comment not found", http.StatusBadRequest)
			fmt.Println("GetCommentByID error")
			return
		}
		like.CommentID = int(id)
		like.PostID = comment.PostID
		isRated := IsUsersComment(user.UserID, like.PostID, like.CommentID)
		if isRated {
			fmt.Println("comment is rated")
			if ok := UpdateLikeDislikeOfComment(like, user.UserID); !ok {
				fmt.Println("UpdateRateOfComment error")
				WarnMessage(w, "Something went wrong")
				// http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
		} else {
			fmt.Println("comment is not rated")
			if ok := AddLikeDislikeToComment(like, user.UserID); !ok {
				fmt.Println(ok)
				WarnMessage(w, "Something went wrong")
				// http.Error(w, "Something went wrong", http.StatusInternalServerError)
				fmt.Println("AddRateToComment error")
				return
			}
		}
		postID := strconv.Itoa(int(like.PostID))
		http.Redirect(w, r, "/post?id="+postID, http.StatusSeeOther)
		}
	}
}


	

