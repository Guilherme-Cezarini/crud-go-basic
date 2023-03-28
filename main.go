package main

import (
	"fmt"
	"net/http"
	"sistema/controllers"
)

func main() {
	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/list", controllers.List)
	http.HandleFunc("/create", controllers.Create)
	http.HandleFunc("/delete", controllers.Delete)
	http.HandleFunc("/edit", controllers.Edit)
	http.HandleFunc("/update", controllers.Update)

	fmt.Println("Server on...")
	http.ListenAndServe(":8000", nil)
}



