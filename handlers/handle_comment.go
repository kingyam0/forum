package handlers

import (
	"fmt"
	"log"
	"time"

	"forum/sqldb"
)

// GetCommentByID ...
func GetCommentByCID(cid int) *Comm {
	comment := NewComment()
	if err1 := sqldb.DB.QueryRow("SELECT commentID, postID, authorID, author, text, creationDate FROM comments WHERE commentID = ?", cid).
		Scan(&comment.CommentID, &comment.PostID, &comment.AuthorID, &comment.Author, &comment.Text, &comment.CreationDate); err1 != nil {
			log.Fatal("GetCommentByID: ", err1)
		return nil
	}
	return comment
}

// GetCommentsOfPost ...
func GetCommentsOfPost(pid int) []*Comm {
	rows, err2 := sqldb.DB.Query("SELECT * FROM comments WHERE postID = ?", pid)
	if err2 != nil {
		return nil
	}
	comments := []*Comm{}
	for rows.Next() {
		var tempComm *Comm = NewComment()
		if err2 = rows.Scan(&tempComm.CommentID, &tempComm.PostID, &tempComm.AuthorID, &tempComm.Author, &tempComm.Text, &tempComm.CreationDate); err2 != nil {

			fmt.Println(err2.Error(), "GetCommentsOfPost error: ", err2)
			return nil
		}
		comments = append(comments, tempComm)
	}
	rows.Close()
	return comments
}

// AddComment ...
func AddAComment(c *Comm) bool {
	stmt, err3 := sqldb.DB.Prepare("INSERT INTO comments (postID, authorID, author, text, creationDate) VALUES (?, ?, ?, ?, ?)")
	defer stmt.Close()
	//var tempComm *Comm
	_, err3 = stmt.Exec(c.PostID, c.AuthorID, c.Author, c.Text, time.Now().Format("2006-01-02 15:04:05"))
	if err3 != nil {
		log.Fatal("AddAComment: ", err3)
		return false
	}
	return true
}

func AddDataToComments(comments []Comm) []Comm {
	for i, comment := range comments {
		err4 := sqldb.DB.QueryRow("SELECT username FROM users WHERE userID=?",
			comment.AuthorID).Scan(&comments[i].Author)
		if err4 != nil {

			log.Fatal("AddDataToComments: ", err4)
		}

		// tempTimeArray := strings.Split(comment.Timestamp, "T")
		// comments[i].Timestamp = tempTimeArray[0]

	}
	return comments
}


func NewComment() *Comm {
	return &Comm{}
}

// GetRateCountOfComment ...
func GetRateCountOfComment(cid, pid int) *Clikedislike {
	rate := NewCommentRating()
	// fmt.Println(cid,pid)
	// fmt.Printf("SELECT likecount, dislikecount FROM commlikedislike WHERE commentID=%d AND postID=%d\n", cid, pid)
	if err5 := sqldb.DB.QueryRow("SELECT * FROM commlikedislike WHERE commentID=? AND postID=?", cid, pid).
		Scan(&rate.CommentID, &rate.UserID, &rate.PostID, &rate.Likecount, &rate.Dislikecount); err5 != nil {
			fmt.Println("GetRateCountOfComment: ", err5)
		// It means nobody rated the comment, likeCount and dislikeCount now is zero
		rate.Likecount = 0
		rate.Dislikecount = 0
	}
	// fmt.Println("hello",rate.Likecount,rate.Dislikecount)
	return rate
}


