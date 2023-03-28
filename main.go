package main

import (
	"fmt"
	"net/http"
	"sistema/controllers"
	"sistema/database/seed"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
    if err != nil {
        fmt.Println("err loading env")
    }

	seed.CreateUserAdmin()


	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/list", controllers.List)
	http.HandleFunc("/create", controllers.Create)
	http.HandleFunc("/delete", controllers.Delete)
	http.HandleFunc("/edit", controllers.Edit)
	http.HandleFunc("/update", controllers.Update)
	http.HandleFunc("/logout", controllers.Logout)


	fmt.Println("Server on...")
	http.ListenAndServe(":8000", nil)
}



