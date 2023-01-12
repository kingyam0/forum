package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"forum/sqldb"

	_ "github.com/mattn/go-sqlite3"
)

func AddPost(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseGlob("static/templates/posts.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error, Parsing Error", http.StatusInternalServerError)
		WarnMessage(w, "500 Internal Server Error, Parsing Error")
		return
	}
	fmt.Println("\n AddPost\n err1:", err)
	tpl.Execute(w, nil)
}

// this function allows a user that is registered to create a new post
func CreateAPost(w http.ResponseWriter, r *http.Request) {
	c, err3 := r.Cookie("session_token")
	if err3 != nil {
		fmt.Println("err3 in session_token:", err3)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	currentUser := GetUserByCookie(c.Value)
	// fmt.Println("\n currentuser:username:", currentUser.Username)

	if r.Method == "POST" {
		// Get all of the post's info.
		authorID := currentUser.UserID
		author := currentUser.Username
		title := r.FormValue("title")
		// postText := r.FormValue("postText")
		topics := r.FormValue("topics")
		category1 := r.FormValue("category1")
		category2 := r.FormValue("category2")
		category3 := r.FormValue("category3")
		category4 := r.FormValue("category4")
		content := r.FormValue("postInput")

		// fmt.Println("authorID:", authorID)
		// fmt.Println("author:", author)
		// fmt.Println("title:", title)
		// fmt.Println("category:", category)
		// fmt.Println("content:", content)
		// fmt.Println()

		if title == "" && category1 == "" && category2 == "" && category3 == "" && category4 == "" && content == "" {
			WarnMessage(w, "400 Bad Request, there was a problem creating your post")
			return
		}

		if category1 == "" && category2 == "" && category3 == "" && category4 == "" && content == "" {
			WarnMessage(w, "400 Bad Request, there was a problem creating your post")
			return
		}

		if topics == "" && category1 == "" && category2 == "" && category3 == "" && category4 == "" {
			WarnMessage(w, "400 Bad Request, there was a problem creating your post")
			return
		}
		if content == "" || title == "" {
			WarnMessage(w, "400 Bad Request, there was a problem creating your post")
			return
		}
		// fmt.Println("Adding a post to our DB")

		// uuid := uuid.NewV4().String()
		InsrtPostStmt, err4 := sqldb.DB.Prepare("INSERT INTO posts (authorID, author, title, text, category1, category2, category3, category4, creationDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);")

		// fmt.Println("InsrtPostStmt:", InsrtPostStmt)
		// fmt.Println("err4:", err4)
		// fmt.Println()

		if err4 != nil {
			// fmt.Println("err4 InsrtPostStmt for db:", err4)
			WarnMessage(w, "500 Internal Server Error, there was a problem creating your post")
		}
		defer InsrtPostStmt.Close()

		result, err5 := InsrtPostStmt.Exec(authorID, author, title, content, category1, category2, category3, category4, time.Now().Format("2006-01-02 15:04:05"))
		// fmt.Println(result)
		// fmt.Println(err5)
		rowsAff, _ := result.RowsAffected()
		lastIns, _ := result.LastInsertId()
		fmt.Println("rowsAff:", rowsAff)
		fmt.Println("err:", lastIns)
		// fmt.Println("err5:", err5)
		if err5 != nil {
			// fmt.Println("err5 inserting new post:", err5)
			WarnMessage(w, "500 Internal Server Error, there was a problem saving your post to the database")
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func NewPost() *Pitem {
	return &Pitem{}
}

// gel all posts to show on homepage
func ReadAllPosts() []*Pitem {
	var posts []*Pitem
	rows, errRAP := sqldb.DB.Query("SELECT * FROM posts;")
	if errRAP != nil {
		fmt.Println("errRAP: ", errRAP)
		return nil
	}
	for rows.Next() {
		var tempPost *Pitem = NewPost()
		err := rows.Scan(&tempPost.PostID, &tempPost.AuthorID, &tempPost.Author, &tempPost.Title, &tempPost.Text, &tempPost.Category1, &tempPost.Category2, &tempPost.Category3, &tempPost.Category4, &tempPost.CreationDate)
		if err != nil {
			fmt.Println("err: ", err)
		}
		posts = append(posts, tempPost)
	}
	rows.Close()
	return posts
}

// get posts by userID to show into my posts
func GetPostByUID(uid int) []*Pitem {
	var posts []*Pitem
	query, err1 := sqldb.DB.Query("SELECT * FROM posts WHERE authorID = ? ORDER BY postID DESC", uid)
	if err1 != nil {
		fmt.Println("GetPostByUID error", err1.Error())
		return nil
	}
	defer query.Close()
	for query.Next() {
		post := NewPost()
		if err := query.Scan(&post.PostID, &post.AuthorID, &post.Author, &post.Title, &post.Text, &post.Category1, &post.Category2, &post.Category3, &post.Category4, &post.CreationDate); err != nil {
			fmt.Println("GetPostByUID: ", err)
			return nil
		}
		posts = append(posts, post)
	}
	return posts
}

// get posts by category to show into filter by category
func GetPostByCategory(cat1, cat2, cat3, cat4 string) []*Pitem {
	var posts []*Pitem
	query, err2 := sqldb.DB.Query("SELECT postID, authorID, author, title, text, category1, category2, category3, category4, creationDate FROM posts WHERE category1 = ? OR category2 = ? OR category3 = ? OR category4 = ?", cat1, cat2, cat3, cat4)
	if err2 != nil {
		fmt.Println("err2: ", err2)
		return nil
	}
	defer query.Close()
	for query.Next() {
		post := NewPost()
		if err := query.Scan(&post.PostID, &post.AuthorID, &post.Author, &post.Title, &post.Text, &post.Category1, &post.Category2, &post.Category3, &post.Category4, &post.CreationDate); err != nil {
			fmt.Println("GetPostByCat: ", err)
			return nil
		}
		posts = append(posts, post)
	}
	fmt.Println(posts)
	return posts
}

// get posts by PostID for postview
func GetPostByPID(pid int) *Pitem {
	post := NewPost()
	if err := sqldb.DB.QueryRow("SELECT postID,authorID, author, title, text, category1, category2, category3, category4, creationDate FROM posts WHERE PostID = ?", pid).Scan(&post.PostID, &post.AuthorID, &post.Author, &post.Title, &post.Text, &post.Category1, &post.Category2, &post.Category3, &post.Category4, &post.CreationDate); err != nil {
		fmt.Println("GetPostByPID: ", err)
		return nil
	}
	post.PostID = pid
	return post
}

func PostView(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseGlob("static/templates/postview.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error, Parsing Error", http.StatusInternalServerError)
		WarnMessage(w, "500 Internal Server Error, Parsing Error")
		return
	}
	type PostAttr struct {
		Comments []*Comm
		CommLD   *Clikedislike
	}
	var singlePost struct {
		PostInfo []*PostAttr
		AuthUser *User
		Post     *Pitem
		PostLD   *Plikedislike
	}

	id, errID := strconv.Atoi(r.URL.Query().Get("id"))
	if errID != nil {
		fmt.Println("errID != nil")
		WarnMessage(w, "404 Not Found, Invalid input id")
		return
	}
	cookie, err1 := r.Cookie("session_token")
	if err1 != nil {
		// if user is guest
		postAttr := &PostAttr{}
		singlePost.Post = GetPostByPID(int(id))
		singlePost.PostLD = GetLikeDislikeCountOfPost(int(id))
		postAttr.Comments = GetCommentsOfPost(int(id))
		commentRate := NewCommentRating()
		for i := 0; i < len(postAttr.Comments); i++ {
			commentRate = GetRateCountOfComment(postAttr.Comments[i].CommentID, int(id))
			postAttr.Comments[i].Likes = commentRate.Likecount
			postAttr.Comments[i].Dislikes = commentRate.Dislikecount
		}
		// postAttr.Category, _ = sqldb.DB.GetThreadOfPost(int(id))
		singlePost.PostInfo = append(singlePost.PostInfo, postAttr)
		singlePost.Post.PostID = int(id)
		tpl.Execute(w, singlePost)
		singlePost.Post = nil
		singlePost.PostInfo = nil

	} else if err1 == nil {
		// if user is logged in
		postAttr := &PostAttr{}
		user := GetUserByCookie(cookie.Value)
		singlePost.AuthUser = user
		singlePost.Post = GetPostByPID(int(id))
		singlePost.PostLD = GetLikeDislikeCountOfPost(int(id))
		postAttr.Comments = GetCommentsOfPost(int(id))
		commentRate := NewCommentRating()
		for i := 0; i < len(postAttr.Comments); i++ {
			commentRate = GetRateCountOfComment(postAttr.Comments[i].CommentID, int(id))
			postAttr.Comments[i].Likes = commentRate.Likecount
			postAttr.Comments[i].Dislikes = commentRate.Dislikecount
		}
		// postAttr.Category, _ = GetCategoryOfPost(int(id))
		singlePost.PostInfo = append(singlePost.PostInfo, postAttr)
		singlePost.Post.PostID = int(id)

		// getting new comments from user
		if r.Method == "POST" {
			r.ParseForm()
			comment := NewComment()
			comment.Text = r.PostFormValue("comment")
			comment.CreationDate = time.Now()
			comment.PostID = int(id)
			comment.AuthorID = user.UserID
			comment.Author = user.Username

			if comment.Text == "" {
				WarnMessage(w, "400 Bad Request, there was a problem adding your comment")
				return
			}
			if ok := AddAComment(comment); !ok {
				fmt.Println("AddComment error: ", ok)
				WarnMessage(w, "500 Internal Server Error, Something went wrong")
				// http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}

			postID := strconv.Itoa(id)
			http.Redirect(w, r, "/post?id="+postID, http.StatusSeeOther)
		} else {
			tpl.Execute(w, singlePost)
		}
		singlePost.Post = nil
		singlePost.PostInfo = nil
		singlePost.AuthUser = nil
	}
}
