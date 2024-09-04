package main

import (
	"log"
	"net/http"

	h "github.com/assaidy/personal-blog-api/handlers"
	"github.com/gorilla/mux"
)

const Port = ":8080"

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/posts", h.Make(h.HandleCreatePost)).
		Methods("POST")
	router.HandleFunc("/posts", h.Make(h.HandleGetAllPosts)).
		Methods("GET")
	router.HandleFunc("/posts/{id:[0-9]+}", h.Make(h.HandleGetPostById)).
		Methods("GET")
	router.HandleFunc("/posts/{id:[0-9]+}", h.Make(h.HandleUpdatePostById)).
		Methods("PUT")
	router.HandleFunc("/posts/{id:[0-9]+}", h.Make(h.HandleDeletePostById)).
		Methods("DELETE")
	// TODO: add search by term

	log.Printf("starting server at port %s", Port)
	log.Fatal(http.ListenAndServe(Port, router))
}
