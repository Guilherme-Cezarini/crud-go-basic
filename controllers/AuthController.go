package controllers

import(
	"fmt"
	"net/http"
	"crypto/md5"
	"encoding/hex"
	//"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/go-sql-driver/mysql"
)

var store = sessions.NewCookieStore([]byte("my_secret_key"))


func Index(w http.ResponseWriter, r *http.Request) {
	
	if r.Method == "POST" {
		email := r.FormValue("email")
		formPassword := r.FormValue("password")

		DB := conectionDB()
		selectRecord, err := DB.Query("SELECT `passaword` FROM `users` WHERE `email` = ? LIMIT 1", email)
		selectRecord.Scan()
		for selectRecord.Next() {
			fmt.Println(formPassword)
			var password string
			err = selectRecord.Scan(&password)
			if err != nil {
				panic(err.Error())
			}

			hash := md5.New()
			defer hash.Reset()
			hash.Write([]byte(formPassword))
			formPassword = hex.EncodeToString(hash.Sum(nil))

			session, _ := store.Get(r, "session.id")
			if formPassword == password{
				session.Values["authenticated"] = true
				session.Save(r, w)

				fmt.Println("senha certa")
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

