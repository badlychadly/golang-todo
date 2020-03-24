package main

import (
	// "fmt"
	"log"
	"net/http"
)

// func homeLink(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Welcome home!")
// }

func main() {
	// var err error
	DB = &LDB{}
	DB.Initialize("todo.db")


	// router := mux.NewRouter().StrictSlash(true)
	// router.HandleFunc("/", homeLink)
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":3001", router))
}


