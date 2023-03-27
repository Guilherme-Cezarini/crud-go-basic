package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"crypto/md5"
	"encoding/hex"

	"gopkg.in/go-playground/validator.v9"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func conectionDB() (conection *sql.DB) {
	Driver := "mysql"
	User := "root"
	Password := "Gu1n0m0@"
	Database := "go-crud"

	con, err := sql.Open(Driver, User+":"+Password+"@tcp(127.0.0.1)/"+Database)
	if err != nil {
		panic(err.Error())
	}

	return con

}

var front = template.Must(template.ParseGlob("front/*"))

type User struct {
	Id       int
	Name     string `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/update", Update)

	fmt.Println("Server on...")
	http.ListenAndServe(":8000", nil)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("id")
	DB := conectionDB()
	deleteRecord, err := DB.Prepare("DELETE FROM `users` WHERE `id` = ?")
	if err != nil {
		panic(err.Error())
	}
	deleteRecord.Exec(userId)
	http.Redirect(w, r, "/", 301)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("id")
	DB := conectionDB()
	selectRecord, err := DB.Query("SELECT * FROM `users` WHERE `id` = ?", userId)
	if err != nil {
		panic(err.Error())
	}
	selectRecord.Scan()
	user := User{}
	for selectRecord.Next() {
		var id int
		var name, email, password string
		err = selectRecord.Scan(&id, &name, &email, &password)
		if err != nil {
			panic(err.Error())
		}
		user.Id = id
		user.Name = name
		user.Email = email
		user.Password = password
	}

	front.ExecuteTemplate(w, "edit", user)
}

func Index(w http.ResponseWriter, r *http.Request) {

	DB := conectionDB()
	records, err := DB.Query("SELECT * FROM `users`")
	if err != nil {
		panic(err.Error())
	}
	user := User{}
	arrayUser := []User{}

	for records.Next() {
		var id int
		var name, email, passaword string
		err = records.Scan(&id, &name, &email, &passaword)
		if err != nil {
			panic(err.Error())
		}
		user.Id = id
		user.Name = name
		user.Email = email

		arrayUser = append(arrayUser, user)

	}

	front.ExecuteTemplate(w, "index", arrayUser)
}

func Create(w http.ResponseWriter, r *http.Request) {

	front.ExecuteTemplate(w, "create", nil)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		if len(password) == 0 {
			UpdateRecordWithoutPassword(id, name, email)
		}else{
			UpdateRecordWithPassword(id, name, email, password)
		}

		http.Redirect(w, r, "/", 301)
	}
}

func UpdateRecordWithPassword(id string, name string, email string, password string){
	DB := conectionDB()
	updateRecord, err := DB.Prepare("UPDATE `users` SET `name` = ?, `email` = ?, `passaword` = ? WHERE `id` = ?")
	if err != nil {
		panic(err.Error())
	}
	hash := md5.New()
	defer hash.Reset()
	hash.Write([]byte(password))
	password = hex.EncodeToString(hash.Sum(nil))
	
	updateRecord.Exec(name, email, password, id)
}

func UpdateRecordWithoutPassword(id string, name string, email string) {
	DB := conectionDB()
	updateRecord, err := DB.Prepare("UPDATE `users` SET `name` = ?, `email` = ? WHERE `id` = ?")
	updateRecord.Exec(name, email, id)
	if err != nil {
		panic(err.Error())
	}
}

func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		v := validator.New()
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		
		hash := md5.New()
		defer hash.Reset()
		hash.Write([]byte(password))
		password = hex.EncodeToString(hash.Sum(nil))
		user := User{
			Name:     name,
			Email:    email,
			Password: password,
		}
		
		fmt.Println(user)
		validationError := v.Struct(user)
		if validationError != nil {
			fmt.Println(validationError)
			front.ExecuteTemplate(w, "create", "erro")
			return

		}

		DB := conectionDB()
		insertRegistros, err := DB.Prepare("INSERT INTO `users` (name,email,passaword) VALUES (?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insertRegistros.Exec(name, email, password)

		http.Redirect(w, r, "/", 301)

	}
}

func flashMessage(c *gin.Context, message string) {
	session := sessions.Default(c)
	session.AddFlash(message)
	if err := session.Save(); err != nil {
		log.Printf("error in flashMessage saving session: %s", err)
	}
}
