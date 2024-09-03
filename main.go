package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const Port = ":8080"

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/posts", Make(HandleCreatePost)).Methods("POST")
	router.HandleFunc("/posts", Make(HandleGetAllPosts)).Methods("GET")
	router.HandleFunc("/posts/{id:[0-9]+}", Make(HandleGetPostById)).Methods("GET")
	router.HandleFunc("/posts/{id:[0-9]+}", Make(HandleUpdatePostById)).Methods("PUT")
	router.HandleFunc("/posts/{id:[0-9]+}", Make(HandleDeletePostById)).Methods("DELETE")
    // TODO: search by term

	log.Printf("starting server at port %s", Port)
	log.Fatal(http.ListenAndServe(Port, router))
}
