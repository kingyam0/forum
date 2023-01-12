package database

import (
	"forum/sqldb"
)

// TABVLES TO CREATE
// USERS, COOKIES, POSTS, COMMENTS, CATEGORIES & PRIVATES CHAT

func CreateDB() {
	// user table
	sqldb.DB.Exec(`CREATE TABLE IF NOT EXISTS "users" (
				"userID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
				"username" TEXT NOT NULL UNIQUE,
				"email" TEXT NOT NULL UNIQUE, 
				"passwordhash" BLOB NOT NULL, 
				"creationDate" TIMESTAMP
				); `)
	// post table
	sqldb.DB.Exec(`CREATE TABLE IF NOT EXISTS "posts" ( 
				"postID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
				"authorID" INTEGER NOT NULL,
				"author" TEXT NOT NULL,
				"title" TEXT NOT NULL, 
				"text" TEXT NOT NULL, 
				"category1" TEXT NOT NULL,
				"category2" TEXT NOT NULL,
				"category3" TEXT NOT NULL,
				"category4" TEXT NOT NULL,
				"creationDate" TIMESTAMP,
				FOREIGN KEY(authorID)REFERENCES users(userID)
				);`)
	// category table TABLE NOT USED YET
	sqldb.DB.Exec(`CREATE TABLE IF NOT EXISTS "category" (
				"postID" TEXT REFERENCES post(postID), 
				"travel" INTEGER,
				"currentaffairs" INTEGER,
				"sports" INTEGER,
				"hobby" INTEGER
				);`)
	// //categories table
	// sqldb.DB.Exec(`CREATE TABLE IF NOT EXISTS "categories" (
	// 			"categoryID" INTEGER PRIMARY KEY,
	// 			"categoryname" TEXT);`)
	// comments table
	sqldb.DB.Exec(`CREATE TABLE IF NOT EXISTS "comments" ( 
				"commentID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
				"postID" INTEGER NOT NULL,
				"authorID" INTEGER NOT NULL,
				"author" TEXT NOT NULL,
				"text" TEXT NOT NULL, 
				"creationDate" TIMESTAMP,
				FOREIGN KEY(postID)REFERENCES posts(postID),
				FOREIGN KEY(authorID)REFERENCES users(userID)
				);`)
	// postlikedislike table
	sqldb.DB.Exec(`CREATE TABLE IF NOT EXISTS "postlikedislike" (
				"postID" INTEGER NOT NULL,
				"userID" INTEGER NOT NULL,
				"likecount" INTEGER NOT NULL, 
				"dislikecount" INTEGER NOT NULL,
				FOREIGN KEY(postID)REFERENCES posts(postID),
				FOREIGN KEY(userID)REFERENCES users(userID)
				);`)
	// commentlikedislike table
	sqldb.DB.Exec(`CREATE TABLE IF NOT EXISTS "commlikedislike" (
				"commentID" INTEGER NOT NULL,
				"userID" INTEGER NOT NULL,
				"postID" INTEGER NOT NULL,
				"likecount" INTEGER NOT NULL, 
				"dislikecount" INTEGER NOT NULL,
				FOREIGN KEY(commentID)REFERENCES comments(commentID),
				FOREIGN KEY(userID)REFERENCES users(userID),
				FOREIGN KEY(postID)REFERENCES posts(postID)
				);`)
	// lduserpost table
	sqldb.DB.Exec(`CREATE TABLE IF NOT EXISTS "lduserpost" (
				"userID" INTEGER NOT NULL,
				"postID" INTEGER NOT NULL,
				"kind" INTEGER NOT NULL,
				FOREIGN KEY(userID)REFERENCES users(userID),
				FOREIGN KEY(postID)REFERENCES posts(postID)
				);`)
	// ldusercomment table
	sqldb.DB.Exec(`CREATE TABLE IF NOT EXISTS "ldusercomment" (
				"commentID" INTEGER NOT NULL,
				"userID" INTEGER NOT NULL,
				"postID" INTEGER NOT NULL,
				"kind" INTEGER NOT NULL,
				FOREIGN KEY(commentID)REFERENCES comments(commentID),
				FOREIGN KEY(userID)REFERENCES users(userID),
				FOREIGN KEY(postID)REFERENCES posts(postID)
				);`)
	// sessions table
	sqldb.DB.Exec(`CREATE TABLE IF NOT EXISTS "session" ( 
				"sessionID" STRING NOT NULL PRIMARY KEY, 
				"userID" INTEGER NOT NULL,
				FOREIGN KEY(userID)REFERENCES users(userID)
				);`)
}
