package models

import(
	"database/sql"
	"crypto/md5"
	"encoding/hex"
	"os"
	"fmt"
)

type User struct {
	Id       int
	Name     string 		`validate:"required"`
	Age 	 string 		`validate:"required,min=1"`
	Email    string 		`validate:"required,email"`
	Password string 		`validate:"required"`
}

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

func UpdateRecordWithoutPassword(id string, name string, email string, age string) {
	DB := conectionDB()
	updateRecord, err := DB.Prepare("UPDATE `users` SET `name` = ?,`age` = ?, `email` = ? WHERE `id` = ?")
	updateRecord.Exec(name, age, email, id)
	if err != nil {
		panic(err.Error())
	}
}

func UpdateRecordWithPassword(id string, name string, email string, password string, age string){
	DB := conectionDB()
	updateRecord, err := DB.Prepare("UPDATE `users` SET `name` = ?, `age` = ?, `email` = ?, `passaword` = ? WHERE `id` = ?")
	if err != nil {
		panic(err.Error())
	}
	hash := md5.New()
	defer hash.Reset()
	hash.Write([]byte(password))
	password = hex.EncodeToString(hash.Sum(nil))

	updateRecord.Exec(name, age, email, password, id)
}

func InsertUser(name string, email string, password string, age string){
	DB := conectionDB()
	insertRegistros, err := DB.Prepare("INSERT IGNORE INTO `users` (name,age,email,passaword) VALUES (?, ?, ? ,?)")
	if err != nil {
		panic(err.Error())
	}
	insertRegistros.Exec(name, age, email, password)
}

func DeleteUser(userId string){
	DB := conectionDB()
	deleteRecord, err := DB.Prepare("DELETE FROM `users` WHERE `id` = ?")
	if err != nil {
		panic(err.Error())
	}
	deleteRecord.Exec(userId)
}