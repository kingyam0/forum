package handlers

import (
	"fmt"

	"forum/sqldb"
)

// NewCommentRating ...
func NewCommentRating() *Clikedislike {
	return &Clikedislike{}
}



// IsUserRateComment ...
func IsUsersComment(uid, pid, cid int) bool {
	var comp int
	query, err := sqldb.DB.Query("SELECT commentID FROM ldusercomment WHERE userID=? AND postID=?", uid, pid)
	if err != nil {
		return false
	}
	defer query.Close()
	for query.Next() {
		if err := query.Scan(&comp); err != nil {
			fmt.Println(err.Error(), "IsUserRateComment")
		}
		if comp == cid {
			return true
		}
		comp = 0
	}
	return false
}

// AddLikeToComment ..
func AddLikeToComment(likeCnt, cid, pid int) bool {
	stmt, err := sqldb.DB.Prepare("UPDATE commlikedislike SET likecount=? WHERE commentID=? AND postID=? ")
	defer stmt.Close()
	_, err = stmt.Exec(likeCnt+1, cid, pid)
	if err != nil {
		fmt.Println("update comment likecount error", err.Error())
		return false
	}
	return true
}

// DeleteLikeFromComment ...
func DeleteLikeFromComment(likeCnt, cid, pid int) bool {
	stmt, err := sqldb.DB.Prepare("UPDATE commlikedislike SET likecount=? WHERE commentID=? AND postID=?")
	defer stmt.Close()
	_, err = stmt.Exec(likeCnt-1, cid, pid)
	if err != nil {
		fmt.Println("update comment likecount error", err.Error())
		return false
	}
	return true
}

// AddDislikeToComment ...
func AddDislikeToComment(dislikeCnt, cid, pid int) bool {
	stmt, err := sqldb.DB.Prepare("UPDATE commlikedislike SET dislikecount=? WHERE commentID=? AND postID=?")
	defer stmt.Close()
	_, err = stmt.Exec(dislikeCnt+1, cid, pid)
	if err != nil {
		fmt.Println("update comment dislikecount error", err.Error())
		return false
	}
	return true
}

func DeleteDislikeFromComment(dislikeCnt, cid, pid int) bool {
	stmt, err := sqldb.DB.Prepare("UPDATE commlikedislike SET dislikecount=? WHERE commentID=? AND postID=?")
	defer stmt.Close()
	_, err = stmt.Exec(dislikeCnt-1, cid, pid)
	if err != nil {
		fmt.Println("update comment dislikecount error", err.Error())
		return false
	}
	return true
}

// AddRateToComment ...
func AddLikeDislikeToComment(l *Clikedislike, uid int) bool {
	rate := GetRateCountOfComment(l.CommentID, l.PostID)
	var kind int // like or dislike
	if l.Like {
		kind = 1
		if rate.Likecount == 0 && rate.Dislikecount == 0 {
			stmnt, err := sqldb.DB.Prepare("INSERT INTO commlikedislike (commentID, postID, userID, likecount, dislikecount) VALUES (?, ?, ?, ?, ?)")
			defer stmnt.Close()
			_, err = stmnt.Exec(l.CommentID, l.PostID, uid, rate.Likecount+1, rate.Dislikecount)
			if err != nil {
				fmt.Println("db Insert CommentLikes error", err.Error())
				return false
			}
		} else {
			// Update column "likeCount" in the table
			if ok := AddLikeToComment(rate.Likecount, l.CommentID, l.PostID); !ok {
				return false
			}
		}
	} else {
		kind = 0
		if rate.Likecount == 0 && rate.Dislikecount == 0 {
			stmnt, err := sqldb.DB.Prepare("INSERT INTO commlikedislike (commentID, postID, userID, likecount, dislikecount) VALUES (?, ?, ?, ?, ?)")
			defer stmnt.Close()
			_, err = stmnt.Exec(l.CommentID, l.PostID, uid, rate.Likecount, rate.Dislikecount+1)
			if err != nil {
				fmt.Println("db Insert CommentLikes error", err.Error())
				return false
			}
		} else {
			// Update column "dislikeCount" in the table
			if ok := AddDislikeToComment(rate.Dislikecount, l.CommentID, l.PostID); !ok {
				return false
			}
		}
	}

	stmnt, err := sqldb.DB.Prepare("INSERT INTO ldusercomment (commentID, postID, userID, kind) VALUES (?, ?, ?, ?)")
	defer stmnt.Close()
	_, err = stmnt.Exec(l.CommentID, l.PostID, uid, kind)
	if err != nil {
		fmt.Println("AddLikeDislikeToUserComment error")
		return false
	}
	return true
}

// DeleteCommentRateFromDB ...
func DeleteCommentLikeDislikeFromDB(uid, cid, pid int) bool {
	stmnt, err := sqldb.DB.Prepare("DELETE FROM ldusercomment WHERE commentID =? AND userID=? AND postID=?")
	defer stmnt.Close()
	_, err = stmnt.Exec(cid, uid, pid)
	if err != nil {
		return false
	}
	// check posts rateCounts, if [0, 0] delete row from CommentRating
	rates := GetRateCountOfComment(cid, pid)
	if rates.Dislikecount == 0 && rates.Likecount == 0 {
		stmnt, err := sqldb.DB.Prepare("DELETE FROM ldusercomment WHERE commentID=? AND postID=?")
		defer stmnt.Close()
		_, err = stmnt.Exec(cid, pid)
		if err != nil {
			return false
		}
	}
	return true
}

// EditCommentRate ...
func EditCommentLikeDislike(before, uid, cid, pid int) bool {
	// change to like
	if before == 0 {

		stmnt, err := sqldb.DB.Prepare("UPDATE ldusercomment SET kind=1 WHERE commentID=? AND postID=? AND userID=?")
		defer stmnt.Close()
		_, err = stmnt.Exec(cid, pid, uid)
		if err != nil {
			fmt.Println("update comment likecount error", err.Error())
			return false
		}
		// change to dislike
	} else {
		stmnt, err := sqldb.DB.Prepare("UPDATE ldusercomment SET kind=0 WHERE commentID=? AND postID=? AND userID=?")
		defer stmnt.Close()
		_, err = stmnt.Exec(cid, pid, uid)
		if err != nil {
			fmt.Println("update comment likecount error", err.Error())
			return false
		}
	}
	return true
}

// UpdateRateOfComment ...
func UpdateLikeDislikeOfComment(rate *Clikedislike, uid int) bool {
	// 1) What user did now? (like or dislike)
	// 2) What user have done before?
	// if user 1) liked and 2) liked ---> Delete like from post
	// If 1) liked 2) dislike ---> Delete like and add dislike
	// If 1) disliked 2)disliked ---> Delete dislike
	// If 1) disliked 2) liked ---> Delete dislike and add like

	var before int
	if err := sqldb.DB.QueryRow("SELECT kind FROM ldusercomment WHERE userID=? AND postID=? AND commentID=?", uid, rate.PostID, rate.CommentID).
		Scan(&before); err != nil {
		fmt.Println("Select kind rate of comment error type")
		return false
	}
	rateCount := GetRateCountOfComment(rate.CommentID, rate.PostID)
	rate.Dislikecount = rateCount.Dislikecount
	rate.Likecount = rateCount.Likecount

	// Scenarios
	if before == 1 && rate.Like {
		// delete like & delete row from lduserpost table
		if ok := DeleteLikeFromComment(rate.Likecount, rate.CommentID, rate.PostID); !ok {
			return false
		}
		// delete row from ldUserPost table
		if ok := DeleteCommentLikeDislikeFromDB(uid, rate.CommentID, rate.PostID); !ok {
			return false
		}

	} else if before == 1 && !(rate.Like) {
		// delete like, add dislike & update kind = 0 from lduserpost table
		if ok := DeleteLikeFromComment(rate.Likecount, rate.CommentID, rate.PostID); !ok {
			return false
		}
		// add dislike 
		if ok := AddDislikeToComment(rate.Dislikecount, rate.CommentID, rate.PostID); !ok {
			return false
		}
		// Update kind equal to '0' into lduserpost table
		if ok := EditCommentLikeDislike(before, uid, rate.CommentID, rate.PostID); !ok {
			return false
		}

	} else if before == 0 && !(rate.Like) {
		// delete dislike
		if ok := DeleteDislikeFromComment(rate.Dislikecount, rate.CommentID, rate.PostID); !ok {
			return false
		}
		// // delete row from ldusercomment
		if ok := DeleteCommentLikeDislikeFromDB(uid, rate.CommentID, rate.PostID); !ok {
			return false
		}
	} else if before == 0 && rate.Like {
		// delete dislike, add like
		if ok := DeleteDislikeFromComment(rate.Dislikecount, rate.CommentID, rate.PostID); !ok {
			return false
		}
		if ok := AddLikeToComment(rate.Likecount, rate.CommentID, rate.PostID); !ok {
			return false
		}
		// Update kind equal to '1'
		if ok := EditCommentLikeDislike(before, uid, rate.CommentID, rate.PostID); !ok {
			return false
		}
	}

	return true
}


