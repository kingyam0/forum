package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"strings"
	"time"
	"unicode"

	"forum/sqldb"

	"golang.org/x/crypto/bcrypt"
)

// registerHandler serves form for registering new users
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseGlob("static/templates/register.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error, Parsing Error", http.StatusInternalServerError)
		WarnMessage(w, "500 Internal Server Error, Parsing Error")
		return
	}
	fmt.Println("*******RegisterHandler running*******")
	tpl.Execute(w, nil)
}

// if errreg != nil {
// 	fmt.Printf("RegisterHandler (ParseGlob) error: %+v/n", errreg)
// }

// registerAuthHandler creates new user in database
func RegisterAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("****RegisterAuthRunning****")
	/*
		1. check username criteria
		2. check password criteria
		3. check if email address is valid
		4. check if username & email already exists in database
		5. create bcrypt hash from password
		6. insert username, email and password hash into database
	*/

	r.ParseForm()

	username := r.FormValue("username")
	fmt.Println(username)

	// 1. Check username
	// 1.1 Check username for alphanumeric characters
	nameAlphaNumeric := true
	for _, char := range username {
		// func IsLetter(r rune) bool, func IsNumber(r rune) bool
		// if unicode.IsLetter(char) == true && unicode.IsNumber(char) == true {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			nameAlphaNumeric = false
		}
	}
	// 1.2 Check username Length is 5 < x < 15
	var nameLength bool
	if len(username) >= 5 && len(username) <= 15 {
		nameLength = true
	}

	if !nameAlphaNumeric || !nameLength {
		WarnMessage(w, "400 Bad Request, please check username and password criteria")
		return
	}

	// 2 Check password criteria valid given all conditions are met
	password := r.FormValue("password")

	fmt.Println("password:", password, "\npswdLength:", len(password))
	// 2.1 variable that must pass for password creation criteria
	var pswdLowercase, pswdUppercase, pswdNumber, pswdSpecial, pswdLength, pswdNoSpaces bool
	pswdNoSpaces = true
	for _, char := range password {
		switch {
		// func IsLower(r rune) bool
		case unicode.IsLower(char):
			pswdLowercase = true
		// func IsUpper(r rune) bool
		case unicode.IsUpper(char):
			pswdUppercase = true
		// func IsNumber(r rune) bool
		case unicode.IsNumber(char):
			pswdNumber = true
		// func IsPunct(r rune) bool, IsSymbol(r rune) bool
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			pswdSpecial = true
		// func IsSpace(r rune) bool, type rune = int32
		case unicode.IsSpace(int32(char)):
			pswdNoSpaces = false
		}
	}

	// 2.2 Check password length is correct 8 < x < 20
	if len(password) > 8 && len(password) < 20 {
		pswdLength = true
	}
	fmt.Println("pswdLowercase:", pswdLowercase, "\npswdUppercase:", pswdUppercase, "\npswdNumber:", pswdNumber, "\npswdSpecial:", pswdSpecial, "\npswdLength", pswdLength, "\npswdNoSpaces:", pswdNoSpaces)
	if !pswdLowercase || !pswdUppercase || !pswdNumber || !pswdSpecial || !pswdLength || !pswdNoSpaces {
		WarnMessage(w, "400 Bad Request, please check username and password criteria")
		return
	}

	// 3. check if email address is valid
	email := r.FormValue("email")
	fmt.Println(email)
	i := strings.Index(email, "@")
	fmt.Println("i:", i)

	domain := email[i+1:]
	fmt.Println("Domain: ", domain)

	_, err2 := net.LookupMX(domain)
	// _, err3 := mail.ParseAddress(email)
	if err2 != nil {
		fmt.Println("err2 invalid email", err2)
		WarnMessage(w, "400 Bad Request, Please enter a valid email address")
		return

	}

	// 4.1 Check if email already exists
	emailStmt := "SELECT userID FROM users WHERE email = ?"
	rowE := sqldb.DB.QueryRow(emailStmt, email)
	var eID string
	errE := rowE.Scan(&eID)
	if errE != sql.ErrNoRows {
		fmt.Println("email already exists, errE:", errE)
		WarnMessage(w, "400 Bad Request, email already registered, try logging in instead")
		return
	}

	// 4.2 Check if username already exists for availability
	stmt := "SELECT userID FROM users WHERE username = ?"
	row := sqldb.DB.QueryRow(stmt, username)
	var uID string
	errU := row.Scan(&uID)
	if errU != sql.ErrNoRows {
		fmt.Println("username already exists, errU:", errU)
		WarnMessage(w, "400 Bad Request, username already taken")
		return
	}

	// 5 Create bcrypt hash from password
	var hash []byte
	// func GenerateFromPassword(password []byte, cost int) ([]byte, error)
	hash, err4 := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err4 != nil {
		fmt.Println("bcrypt err4:", err4)
		WarnMessage(w, "500 Internal Error, there was a problem registering account")
		return
	}
	fmt.Println("hash:", hash)
	fmt.Println("string(hash):", string(hash))

	// 	6. insert username and password hash in database

	// func (db *DB) Prepare(query string) (*stmt, error)
	// var insertStmt *sql.Stmt
	insertStmt, err5 := sqldb.DB.Prepare("INSERT INTO users (username, email, passwordhash, creationDate) VALUES (?, ?, ?, ?);")
	fmt.Println(insertStmt)
	if err5 != nil {
		fmt.Println("err5 preparing statement:", err5)
		WarnMessage(w, "400 Bad Request, there was a problem registering account")
		return
	}
	defer insertStmt.Close()

	// var result sql.Result
	// func (s *Stmt) Exec(args ...interface{}) (Result, error)
	result, err6 := insertStmt.Exec(username, email, hash, time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(result)

	// checking if the result has been added and the last inserted row
	rowsAff, _ := result.RowsAffected()
	lastIns, _ := result.LastInsertId()
	fmt.Println("rowsAff:", rowsAff)
	fmt.Println("lastIns:", lastIns)
	if err6 != nil {
		fmt.Println("error inserting new user", err6)
		WarnMessage(w, "400 Bad Request, there was a problem registering account")
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
	// tpl.Execute(w, "Congrats, your account has been successfully created, Please Log in")
}
