package main 


import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
)

func dMeth(myslice interface{}) {
	arr, ok := myslice.([]int)
	if ok {
		fmt.Printf("my arr %v\n", arr )
		// return
	}
	newArr := append(arr, 10)
	fmt.Printf("out of Block %v\n", newArr )
}


func HomePage(w http.ResponseWriter, r *http.Request) {
	lists := DB.ViewLists()
	for i, list := range lists {
		fmt.Fprintf(w, "list: %v, i: %v\n", list, i)
	}
		// dMeth([]int{1,2,3,4,5})
}



func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", HomePage).Methods("GET")
	router.HandleFunc("/", HandleNewList).Methods("POST")
	router.HandleFunc("/{id}", HandleListView).Methods("GET")
	router.HandleFunc("/{id}", HandleDeleteList).Methods("DELETE")
	router.HandleFunc("/{id}/items", HandleNewItem).Methods("POST")
	return router
}

func HandleNewList(w http.ResponseWriter, r *http.Request) {
	var list List
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&list)
	if err != nil {
		fmt.Fprintf(w, "decoder failed")
	}
	err = DB.CreateList(&list)
	if err != nil {
		fmt.Fprintf(w, "DB error")
	}
	// listBytes, _ := json.Marshal(list)
	// fmt.Fprintf(w, string(listBytes))
	json.NewEncoder(w).Encode(list)
}


func HandleListView(w http.ResponseWriter, r *http.Request) {
	listId := mux.Vars(r)["id"]
	list, err := DB.ViewList(listId)
	// fmt.Printf("id Type: %v\n", listId)
	// if list == (List{}) {
	// 	fmt.Fprint(w, "No results found")
	// 	return
	// }

	if err != nil {
		fmt.Fprintf(w, "No results found %v\n", err)
		return
	}
	fmt.Printf("returned list: %v\n", list)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	json.NewEncoder(w).Encode(list)
}


func HandleDeleteList(w http.ResponseWriter, r *http.Request) {
	listId := mux.Vars(r)["id"]
	err := DB.DeleteList(listId)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, "Successfully Deleted")
}

func HandleNewItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	listId := mux.Vars(r)["id"]
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&item)
	// item.ListId = listId
	err = DB.CreateItem(&item, listId)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	// fmt.Fprintf(w, "Item Added")
	json.NewEncoder(w).Encode(item)
}