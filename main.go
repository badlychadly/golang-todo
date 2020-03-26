package main

import (
	// "fmt"
	"log"
	"net/http"
)



func main() {
	DB = &LDB{}
	DB.Initialize("todo.db")


	router := NewRouter()
	log.Fatal(http.ListenAndServe(":3001", router))
}


