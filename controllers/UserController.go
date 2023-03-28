package controllers
import(
	"fmt"
	"database/sql"
	"net/http"
	"os"
	"text/template"
	"crypto/md5"
	"encoding/hex"
	"sistema/database/models"
	

	"gopkg.in/go-playground/validator.v9"

	_ "github.com/go-sql-driver/mysql"
)


var front = template.Must(template.ParseGlob("front/*"))

func conectionDB() (conection *sql.DB) {
	Driver := "mysql"
	User := os.Getenv("DB_USER")
	Password := os.Getenv("DB_PASSWORD")
	Database := os.Getenv("DB_DATABASE")
	fmt.Println(User + " - " + Password + " - " + Database)

	con, err := sql.Open(Driver, User+":"+Password+"@tcp(127.0.0.1)/"+Database)
	if err != nil {
		panic(err.Error())
	}

	return con

}

func Create(w http.ResponseWriter, r *http.Request) {
	SessionHealthCheck(w,r)
	if r.Method == "POST" {
		v := validator.New()
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		age := r.FormValue("age")
		
		hash := md5.New()
		defer hash.Reset()
		hash.Write([]byte(password))
		password = hex.EncodeToString(hash.Sum(nil))
		user := models.User{
			Name:     name,
			Age: 	  age,
			Email:    email,
			Password: password,
		}
		
		fmt.Println(user)
		validationError := v.Struct(user)
		if validationError != nil {
			fmt.Println(validationError)
			front.ExecuteTemplate(w, "create", validationError)
			return

		}
		models.InsertUser(name, email, password, age)

		http.Redirect(w, r, "/list", 301)

	}
	front.ExecuteTemplate(w, "create", nil)
}


func Delete(w http.ResponseWriter, r *http.Request) {
	SessionHealthCheck(w,r)
	userId := r.URL.Query().Get("id")
	models.DeleteUser(userId)
	
	http.Redirect(w, r, "/list", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	SessionHealthCheck(w,r)
	if r.Method == "POST" {
		v := validator.New()
		id := r.FormValue("id")
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		age := r.FormValue("age")

		user := models.User{
			Name:     name,
			Age: 	  age,
			Email:    email,
			Password: password,
		}
		fmt.Println(user)
		validationError := v.Struct(user)
		if validationError != nil {
			fmt.Println(validationError)
			http.Redirect(w, r, "/list", 301)
		}

		if len(password) == 0 {
			models.UpdateRecordWithoutPassword(id, name, email,age)
		}else{
			models.UpdateRecordWithPassword(id, name, email, password, age)
		}

		http.Redirect(w, r, "/list", 301)
	}
}


func Edit(w http.ResponseWriter, r *http.Request) {
	SessionHealthCheck(w,r)
	userId := r.URL.Query().Get("id")
	DB := conectionDB()
	selectRecord, err := DB.Query("SELECT * FROM `users` WHERE `id` = ?", userId)
	if err != nil {
		panic(err.Error())
	}
	selectRecord.Scan()
	user := models.User{}
	for selectRecord.Next() {
		var id int
		var name, age, email, password string
		err = selectRecord.Scan(&id, &name, &age, &email, &password)
		if err != nil {
			panic(err.Error())
		}

		user.Id = id
		user.Name = name
		user.Email = email
		user.Password = password
		user.Age = age
	}

	front.ExecuteTemplate(w, "edit", user)
}

func List(w http.ResponseWriter, r *http.Request) {
	SessionHealthCheck(w,r)
	DB := conectionDB()
	records, err := DB.Query("SELECT * FROM `users`")
	if err != nil {
		panic(err.Error())
	}
	user := models.User{}
	arrayUser := []models.User{}

	for records.Next() {
		var id int
		var name, age, email, passaword string
		err = records.Scan(&id, &name, &age, &email, &passaword)
		if err != nil {
			panic(err.Error())
		}
		user.Id = id
		user.Name = name
		user.Email = email
		user.Age = age

		arrayUser = append(arrayUser, user)

	}

	front.ExecuteTemplate(w, "list", arrayUser)
}


