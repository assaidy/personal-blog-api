package main

import (
	"log"
	"net/http"

	"github.com/assaidy/personal-blog-api/handlers"
	"github.com/gorilla/mux"
)

const Port = ":8080"

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/posts", handlers.Make(handlers.HandleCreatePost)).
		Methods("POST")
	router.HandleFunc("/posts", handlers.Make(handlers.HandleGetAllPosts)).
		Methods("GET")
	router.HandleFunc("/posts/{id:[0-9]+}", handlers.Make(handlers.HandleGetPostById)).
		Methods("GET")
	router.HandleFunc("/posts/{id:[0-9]+}", handlers.Make(handlers.HandleUpdatePostById)).
		Methods("PUT")
	router.HandleFunc("/posts/{id:[0-9]+}", handlers.Make(handlers.HandleDeletePostById)).
		Methods("DELETE")
	// TODO: test the api endpoints
	// TODO: add search by term

	log.Printf("starting server at port %s", Port)
	log.Fatal(http.ListenAndServe(Port, router))
}
