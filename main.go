package main

import (
	"fmt"
	"forum/database"
	"forum/handlers"
	"forum/sqldb"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	sqldb.ConnectDB()
	database.CreateDB()

	/*path := http.FileServer(http.Dir("/static/css")) //to handle the css folder
	http.Handle("/css/", http.StripPrefix("/css/", path))*/
	path2 := http.FileServer(http.Dir("static")) //to handle the assets folder
	http.Handle("/static/", http.StripPrefix("/static/", path2))

	// fmt.Println("forum starts")
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/loginauth", handlers.LoginAuthHandler)
	http.HandleFunc("/logout", handlers.Logout)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/registerauth", handlers.RegisterAuthHandler)
	http.HandleFunc("/addpost", handlers.AddPost)
	http.HandleFunc("/createpost", handlers.CreateAPost)
	http.HandleFunc("/filter", handlers.FilterHandle)
	http.HandleFunc("/post", handlers.PostView)
	http.HandleFunc("/rate", handlers.LikeDislikeHandler)

	fmt.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	sqldb.CloseDB()
}
