package seed
import(
	"sistema/database/models"
	"os"
	"crypto/md5"
	"encoding/hex"
	"database/sql"
	"fmt"


)

func conectionDB() (conection *sql.DB) {
	Driver := "mysql"
	User := os.Getenv("DB_USER")
	Password := os.Getenv("DB_PASSWORD")
	Database := os.Getenv("DB_DATABASE")

	con, err := sql.Open(Driver, User+":"+Password+"@tcp(127.0.0.1)/"+Database)
	if err != nil {
		panic(err.Error())
	}

	return con

}

func CreateUserAdmin() {
	email := os.Getenv("ADMIN_EMAIL")
	DB := conectionDB()
	rows, err := DB.Query("SELECT `passaword` FROM `users` WHERE `email` = ? LIMIT 1", email)
	if err != nil {
		
		return
	} 
	count := make([]string, 0)
	for rows.Next() {
		var password string
		if err := rows.Scan(&password); err != nil {
			return
		}
		count = append(count, password)
	}

	
	if len(count) == 0 {
		var passwordAdmin string
		hash := md5.New()
		defer hash.Reset()
		hash.Write([]byte(os.Getenv("ADMIN_PASSWORD")))
		passwordAdmin = hex.EncodeToString(hash.Sum(nil))
		models.InsertUser(os.Getenv("ADMIN_NAME"), os.Getenv("ADMIN_EMAIL"), passwordAdmin, os.Getenv("ADMIN_AGE"))
		fmt.Println("Admin created.")
	}

	
		
}