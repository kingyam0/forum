package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"forum/sqldb"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var expiresAt time.Time

// loginHandler serves form for users to login with
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseGlob("static/templates/login.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error, Parsing Error", http.StatusInternalServerError)
		WarnMessage(w, "500 Internal Server Error, Parsing Error")
		return
	}
	// 	if errlogin != nil {
	// 		fmt.Printf("LoginHandler (ParseGlob) error: %+v/n", errlogin)
	// 	}
	fmt.Println("*******LoginHandler running*******")
	tpl.Execute(w, nil)
}

// loginAuthHandler authenticates user login
func LoginAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*******LoginAuthHandler running*******")

	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println("Logged in user:", username)

	// retrieve password from db to compare (hash) with user supplied password's hash
	var hash string
	stmt := "SELECT passwordhash FROM users WHERE Username = ?"
	row := sqldb.DB.QueryRow(stmt, username)
	err2 := row.Scan(&hash)
	if err2 != nil {
		fmt.Println("err2 selecting passwordhash in db by Username", err2)
		WarnMessage(w, "check username and password")
		return
	}

	// func CompareHashAndPassword(hashed password, password []byte) error
	errbcrypt := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// returns nil on success
	if errbcrypt == nil {
		stmtCurrentUer := "SELECT * FROM users WHERE Username = ?"
		rowCurrentUser := sqldb.DB.QueryRow(stmtCurrentUer, username)
		err3 := rowCurrentUser.Scan(&CurrentUser.UserID, &CurrentUser.Username, &CurrentUser.Email, &CurrentUser.passwordhash, &CurrentUser.CreationDate)
		if err3 != nil {
			fmt.Println("err3 rowCred.scan:", err3)
			WarnMessage(w, "error accessing db")
			return
		}

		err3 = IsUserAuthenticated(w, &CurrentUser)
		fmt.Println("this user logged in: ", err3)
		if err3 != nil {
			WarnMessage(w, "You are already logged in 🧐")
			fmt.Println("already logged in: ", err3)
			return
		}

		// login := &User{
		// 	UserID:   CurrentUser.UserID,
		// 	Username: CurrentUser.Username,
		// 	Email: CurrentUser.Email,
		// 	passwordhash: CurrentUser.passwordhash,
		// 	CreationDate: CurrentUser.CreationDate,
		// }

		// storing the cookie values in struct to access on other pages
		// user_session = Cookie{"session_token", sessionToken, expiresAt}

		// storing the cookie values in struct to access on other pages
		// user_session = Cookie{"session_token", sessionToken, expiresAt}

		// Create a new random session token
		sessionToken := uuid.NewV4().String()
		expiresAt = time.Now().Add(120 * time.Minute)

		// Finally, we set the client cookie for "session_token" as the session token we just generated
		// we also set an expiry time of 120 minutes
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			MaxAge:   900,
			Expires:  expiresAt,
			HttpOnly: true,
		})

		// storing the cookie values in struct to access on other pages
		user_session = Cookie{"session_token", sessionToken, expiresAt}

		insertsessStmt, err4 := sqldb.DB.Prepare("INSERT INTO session (sessionID, userID) VALUES (?, ?);")
		// fmt.Println("Session Token:", sessionToken)
		if err4 != nil {
			fmt.Println("err4 preparing statement:", err4)
			WarnMessage(w, "there was a problem logging in")
			return
		}
		defer insertsessStmt.Close()
		insertsessStmt.Exec(sessionToken, CurrentUser.UserID)

		// redirect user to index handler after successful login
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	fmt.Println("incorrect password")
	WarnMessage(w, "check username and password")
	return
}

// NewSession ...
func NewSession() *Session {
	return &Session{}
}

// AddSession ...
func AddSession(w http.ResponseWriter, sessionName string, user *User) {
	sessionToken := uuid.NewV4().String()
	expiresAt = time.Now().Add(120 * time.Minute)

	cookieSession := &http.Cookie{
		Name:     sessionName,
		Value:    sessionToken,
		MaxAge:   900,
		Expires:  expiresAt,
		HttpOnly: true,
	}

	http.SetCookie(w, cookieSession)
	if sessionName != "guest" {
		InsertSession(user, cookieSession)
	}
}

// InsertSession ...
func InsertSession(u *User, session *http.Cookie) *Session {
	cookie := NewSession()
	stmnt, err := sqldb.DB.Prepare("INSERT INTO session (sessionID, userID) VALUES (?, ?)")
	_, err = stmnt.Exec(session.Value, u.UserID)
	if err != nil {
		fmt.Println("AddSession error inserting into DB: ", err)
	}
	cookie.sessionName = session.Name
	cookie.sessionUUID = session.Value
	cookie.UserID = u.UserID
	return cookie
}

// User's cookie expires when browser is closed, delete the cookie from the database.
func DeleteSession(w http.ResponseWriter, cookieValue string) error {
	cookie := &http.Cookie{
		Name:     "Session_token",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	stmt, err := sqldb.DB.Prepare("DELETE FROM session WHERE sessionID=?;")
	defer stmt.Close()
	stmt.Exec(cookieValue)
	if err != nil {
		fmt.Println("DeleteSession err: ", err)
		return err
	}
	return nil
}

// logout handle
func Logout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/logout" {
		c, err := r.Cookie("session_token")
		if err != nil {
			AddSession(w, "guest", nil)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			fmt.Println("Logout error: ", err)
		}
		DeleteSession(w, c.Value)
		fmt.Println("user logged out")
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// NewUser ...
func NewUser() *User {
	return &User{}
}

// FindByUserID ...
func FindByUserID(UID int64) *User {
	u := NewUser()
	if err := sqldb.DB.QueryRow("SELECT userID, username, email, passwordhash FROM users WHERE userID = ?", UID).
		Scan(&u.UserID, &u.Username, &u.Email, &u.passwordhash); err != nil {
		fmt.Println("error find by user: ", err)
		return nil
	}
	return u
}

// GetUserByCookie ...
func GetUserByCookie(cookieValue string) *User {
	var userID int64
	if err := sqldb.DB.QueryRow("SELECT userID from session WHERE sessionID = ?", cookieValue).Scan(&userID); err != nil {
		return nil
	}
	u := FindByUserID(userID)
	return u
}

// IsUserAuthenticated ...
func IsUserAuthenticated(w http.ResponseWriter, u *User) error {
	var cookieValue string
	if err := sqldb.DB.QueryRow("SELECT sessionID FROM session WHERE userID = ?", u.UserID).Scan(&cookieValue); err != nil {
		return nil
	}
	if err := DeleteSession(w, cookieValue); err != nil {
		return err
	}
	return nil
}
