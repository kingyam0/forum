package handlers

import (
	"fmt"

	"forum/sqldb"
)

// GetRateCountOfPost ...
func GetLikeDislikeCountOfPost(postID int) *Plikedislike {
	rates := NewPostRating()
	if err := sqldb.DB.QueryRow("SELECT * FROM postlikedislike WHERE postID = ?", postID).
		Scan(&rates.PostID, &rates.UserID, &rates.Likecount, &rates.Dislikecount); err != nil {
		// It means nobody rated the post, likeCount and dislikeCount now is zero
		rates.Likecount = 0
		rates.Dislikecount = 0
	}
	return rates
}

// NewPostRating ...
func NewPostRating() *Plikedislike {
	return &Plikedislike{}
}

// IsUserRatePost ...
func IsThisUsersPost(uid, pid int) bool {
	var comp int
	// need to check all rated posts of the user
	res, err := sqldb.DB.Query("SELECT postID FROM lduserpost WHERE userID=?", uid)
	if err != nil {
		return false
	}
	defer res.Close()
	// Check postID is rated or not
	for res.Next() {
		if err := res.Scan(&comp); err != nil {
			fmt.Println(err.Error(), "IsThisUsersPost")
		}
		if comp == pid {
			return true
		}
		comp = 0
	}
	return false
}

// GetRatedPostsByUID ...
func GetRatedPostsByUID(uid int) []*Pitem {
	var posts []*Pitem
	postIDs := []int{}
	query, err := sqldb.DB.Query("SELECT postID FROM lduserpost WHERE userID=? AND kind=1 ORDER BY postID DESC", uid)
	if err != nil {
		fmt.Println(err.Error(), "GetRatedPostsByUID")

		return nil
	}
	defer query.Close()
	for query.Next() {
		var id int
		if err := query.Scan(&id); err != nil {
			fmt.Println(err.Error(), "GetRatedPostsByUID")
			return nil
		}
		// Save all rated post ID to array
		postIDs = append(postIDs, id)
	}
	// Get full info of post by PostID
	for _, postID := range postIDs {
		post := GetPostByPID(postID)
		if err != nil {
			fmt.Println(err.Error(), "GetRatedPostsByUID")
			return nil
		}
		posts = append(posts, post)
	}
	return posts
}

// AddLike ...
func AddLike(likeCnt, postID int) bool {
	stmnt, err := sqldb.DB.Prepare("UPDATE postlikedislike SET likecount=? WHERE postID=?")
	defer stmnt.Close()
	_, err = stmnt.Exec(likeCnt+1, postID)
	if err != nil {
		fmt.Println("update likecount error")
		return false
	}
	return true
}

// DeleteLike ...
func DeleteLike(likeCnt, postID int) bool {
	stmnt, err := sqldb.DB.Prepare("UPDATE postlikedislike SET likecount=? WHERE postID=?")
	defer stmnt.Close()
	_, err = stmnt.Exec(likeCnt-1, postID)
	if err != nil {
		fmt.Println("update likecount error")
		return false
	}
	return true
}

// AddDislike ...
func AddDislike(dislikeCnt, postID int) bool {
	stmnt, err := sqldb.DB.Prepare("UPDATE postlikedislike SET dislikecount=? WHERE postID=?")
	defer stmnt.Close()
	_, err = stmnt.Exec(dislikeCnt+1, postID)
	if err != nil {
		fmt.Println("update dislikecount error")
		return false
	}
	return true
}

// DeleteDislike ...
func DeleteDislike(dislikeCnt, postID int) bool {
	stmnt, err := sqldb.DB.Prepare("UPDATE postlikedislike SET dislikecount=? WHERE postID=?")
	defer stmnt.Close()
	_, err = stmnt.Exec(dislikeCnt-1, postID)
	if err != nil {
		fmt.Println("update dislikecount error")
		return false
	}
	return true
}

// DeleteLikeDislikeFromDB ...
func DeleteLikeDislikeFromDB(uid, pid int) bool {
	stmnt, err := sqldb.DB.Prepare("DELETE FROM lduserpost WHERE userID=? AND postID=?")
	defer stmnt.Close()
	_, err = stmnt.Exec(uid, pid)
	if err != nil {
		return false
	}
	// check posts rateCounts, if [0, 0] delete row from PostRating
	rates := GetLikeDislikeCountOfPost(pid)
	if rates.Dislikecount == 0 && rates.Likecount == 0 {
		stmnt, err := sqldb.DB.Prepare("DELETE FROM postlikedislike WHERE postID=?")
		defer stmnt.Close()
		_, err = stmnt.Exec(pid)
		if err != nil {
			return false
		}
	}

	return true
}

// UpdateRate ...
func UpdateLikeDislike(before, uid, pid int) bool {
	// change to like
	if before == 0 {

		stmnt, err := sqldb.DB.Prepare("UPDATE lduserpost SET kind=1 WHERE postID=? AND userID=?")
		defer stmnt.Close()
		_, err = stmnt.Exec(pid, uid)
		if err != nil {
			fmt.Println("update likecount error")
			return false
		}
		// change to dislike
	} else {
		stmnt, err := sqldb.DB.Prepare("UPDATE lduserpost SET kind=0 WHERE postID=? AND userID=?")
		defer stmnt.Close()
		_, err = stmnt.Exec(pid, uid)
		if err != nil {
			fmt.Println("update likecount error")
			return false
		}
	}
	return true
}

// AddRateToPost ...
func AddLikeDislikeToPost(l *Plikedislike, uid int) bool {
	/*Need to handle the situation:
	If user liked the post, DB will insert record with likeCount value only.
	And then for example, another user will dislike the post, DB will add another record
	with this postID and likeCount as null
	*/

	rate := GetLikeDislikeCountOfPost(l.PostID)
	// fmt.Println(rate.Likecount, rate.Dislikecount, "<-- post rates")
	var kind int

	if l.Like {
		kind = 1
		if rate.Likecount == 0 && rate.Dislikecount == 0 {
			stmnt, err := sqldb.DB.Prepare("INSERT INTO postlikedislike (postID, userID, likecount, dislikecount) VALUES (?, ?, ?, ?);")
			defer stmnt.Close()
			fmt.Println("stmnt addlikedislike kind= 1: ", stmnt)
			fmt.Println("stmnt addlikedislike kind= 1: ", err)

			_, err = stmnt.Exec(l.PostID, uid, rate.Likecount+1, rate.Dislikecount)

			if err != nil {
				fmt.Println("db Insert postlikedislike error", err)
				return false
			}

		} else {
			// Update column "likeCount" in the table
			if ok := AddLike(rate.Likecount, rate.PostID); !ok {
				return false
			}
		}
		// If dislike
	} else {
		kind = 0
		if rate.Likecount == 0 && rate.Dislikecount == 0 {
			stmnt, err := sqldb.DB.Prepare("INSERT INTO postlikedislike (postID, userID, likecount, dislikecount) VALUES (?, ?, ?, ?);")
			defer stmnt.Close()
			_, err = stmnt.Exec(l.PostID, uid, rate.Likecount, rate.Dislikecount+1)
			if err != nil {
				fmt.Println("db Insert postlikedislike dislike error")
				return false
			}
		} else {
			// Update column "dislikeCount" in the table
			if ok := AddDislike(rate.Dislikecount, rate.PostID); !ok {
				return false
			}
		}

	}
	stmnt, err := sqldb.DB.Prepare("INSERT INTO lduserpost (userID, postID, kind) VALUES (?, ?, ?)")
	_, err = stmnt.Exec(uid, l.PostID, kind)
	if err != nil {
		fmt.Println("AddLikeDislikeToPost error")
		return false
	}

	return true
}

// UpdateLikeDislikeOfPost ...
func UpdateLikeDislikeOfPost(rate *Plikedislike, uid int) bool {
	// 1) What user did now? (like or dislike)
	// 2) What user have done before?
	// if user 1) liked and 2) liked ---> Delete like from post
	// If 1) liked 2) dislike ---> Delete like and add dislike
	// If 1) disliked 2)disliked ---> Delete dislike
	// If 1) disliked 2) liked ---> Delete dislike and add like

	var before int
	if err := sqldb.DB.QueryRow("SELECT kind FROM lduserpost WHERE userID=? AND postID=?", uid, rate.PostID).
		Scan(&before); err != nil {
		fmt.Println("DeleteLikeDislikeFromPost error type")
		return false
	}
	rateCount := GetLikeDislikeCountOfPost(rate.PostID)
	rate.Dislikecount = rateCount.Dislikecount
	rate.Likecount = rateCount.Likecount

	// Scenarios
	if before == 1 && rate.Like {
		// delete like
		if ok := DeleteLike(rate.Likecount, rate.PostID); !ok {
			return false
		}
		// delete row from RateUserPost
		if ok := DeleteLikeDislikeFromDB(uid, rate.PostID); !ok {
			return false
		}

	} else if before == 1 && !(rate.Like) {
		// delete like, add dislike
		if ok := DeleteLike(rate.Likecount, rate.PostID); !ok {
			return false
		}
		if ok := AddDislike(rate.Dislikecount, rate.PostID); !ok {
			return false
		}
		// Update kind equal to '0'
		if ok := UpdateLikeDislike(before, uid, rate.PostID); !ok {
			return false
		}

	} else if before == 0 && !(rate.Like) {
		// delete dislike
		if ok := DeleteDislike(rate.Dislikecount, rate.PostID); !ok {
			return false
		}
		// delete row from lduserpost
		if ok := DeleteLikeDislikeFromDB(uid, rate.PostID); !ok {
			return false
		}
	} else if before == 0 && rate.Like {
		// delete dislike, add like
		if ok := DeleteDislike(rate.Dislikecount, rate.PostID); !ok {
			return false
		}
		if ok := AddLike(rate.Likecount, rate.PostID); !ok {
			return false
		}
		// Update kind equal to '1'
		if ok := UpdateLikeDislike(before, uid, rate.PostID); !ok {
			return false
		}
	}

	return true
}

// If server delete like, and after it becomes [0 0] need to delete row from PostRating
