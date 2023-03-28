package controllers

import(
	"fmt"
	"time"
	"net/http"
	"crypto/md5"
	"encoding/hex"
	//"github.com/gorilla/mux"
	//"github.com/gorilla/sessions"
	"github.com/google/uuid"
	_ "github.com/go-sql-driver/mysql"

)

var sessions = map[string]session{}

type session struct {
	username string
	expiry   time.Time
}

func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}


func Index(w http.ResponseWriter, r *http.Request) {
	
	if r.Method == "POST" {
		email := r.FormValue("email")
		formPassword := r.FormValue("password")

		DB := conectionDB()
		selectRecord, err := DB.Query("SELECT `passaword` FROM `users` WHERE `email` = ? LIMIT 1", email)
		if err != nil {
			front.ExecuteTemplate(w, "index", "Error in credentials.")
			return
		}
		selectRecord.Scan()
		for selectRecord.Next() {
			var password string
			err = selectRecord.Scan(&password)
			if err != nil {
				panic(err.Error())
			}

			hash := md5.New()
			defer hash.Reset()
			hash.Write([]byte(formPassword))
			formPassword = hex.EncodeToString(hash.Sum(nil))

			if formPassword == password{
				sessionToken := uuid.NewString()
				expiresAt := time.Now().Add(120 * time.Second)

				sessions[sessionToken] = session{
					username: email,
					expiry:   expiresAt,
				}

				http.SetCookie(w, &http.Cookie{
					Name:    "session_token",
					Value:   sessionToken,
					Expires: expiresAt,
				})

				http.Redirect(w, r, "/list", 301)
			}else{
				fmt.Println(password+" "+formPassword)
				front.ExecuteTemplate(w, "index", "Error in credentials.")
				return
			}
		}
		front.ExecuteTemplate(w, "index", "Error in credentials.")
		return
	}
	front.ExecuteTemplate(w, "index", nil)
}

func Logout(w http.ResponseWriter, r *http.Request){

	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Println("Session cookie dont exists")
			http.Redirect(w, r, "/", 301)
			return
		}
		fmt.Println("Session cookie dont exists")
		http.Redirect(w, r, "/", 301)
		return
	}
	sessionToken := c.Value

	// remove the users session from the session map
	delete(sessions, sessionToken)

	// We need to let the client know that the cookie is expired
	// In the response, we set the session token to an empty
	// value and set its expiry as the current time
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	fmt.Println("Logout")
	http.Redirect(w, r, "/", 301)


}



func SessionHealthCheck(w http.ResponseWriter, r *http.Request){
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Println("Session cookie dont exists")
			http.Redirect(w, r, "/", 301)
			return
		}
		fmt.Println("Session cookie dont exists")
		http.Redirect(w, r, "/", 301)
		return
	}
	sessionToken := c.Value

	userSession, exists := sessions[sessionToken]
	if !exists {
		fmt.Println("The session token is not present in session map")
		http.Redirect(w, r, "/", 301)
		return
	}

	if userSession.isExpired() {
		delete(sessions, sessionToken)
		fmt.Println("The session is present, but has expired")
		http.Redirect(w, r, "/", 301)
		return
	}
}

