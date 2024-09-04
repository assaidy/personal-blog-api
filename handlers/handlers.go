package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/assaidy/personal-blog-api/db"
	"github.com/assaidy/personal-blog-api/types"
	"github.com/gorilla/mux"
)

func HandleCreatePost(w http.ResponseWriter, r *http.Request) error {
	post := &types.Post{}

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		return types.InvalidJSONError()
	}

	post.CreatedAt = time.Now().UTC()
	post.UpdatedAt = time.Now().UTC()

	postId, err := db.CreatePost(post)
	if err != nil {
		return err
	}
	post.Id = postId

	return WriteJSON(w, http.StatusOK, post)
}

// func HandleGetAllPosts(w http.ResponseWriter, r *http.Request) error {
// 	posts, err := db.GetAllPosts()
// 	if err != nil {
// 		return err
// 	}
//
// 	return WriteJSON(w, http.StatusOK, posts)
// }

func HandleGetAllPostsByTerm(w http.ResponseWriter, r *http.Request) error {
	term := strings.TrimSpace(r.URL.Query().Get("term"))

	posts, err := db.GetAllPostsByTerm(term)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, posts)
}

func HandleGetPostById(w http.ResponseWriter, r *http.Request) error {
	// NOTE: i didn't handle errors as it's garanteed to get a valid integer
	// with the regex validation by mux
	// see routes in main file to understand
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	post, err := db.GetPost(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, post)
}

func HandleUpdatePostById(w http.ResponseWriter, r *http.Request) error {
	postReq := &types.Post{}

	err := json.NewDecoder(r.Body).Decode(&postReq)
	if err != nil {
		return types.InvalidJSONError()
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	postReq.Id = id

	post, err := db.GetPost(postReq.Id)
	if err != nil {
		return err
	}

	post.Title = postReq.Title
	post.Content = postReq.Content
	post.Category = postReq.Category
	post.Tags = postReq.Tags
	post.UpdatedAt = time.Now().UTC()

	err = db.UpdatePost(post)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, post)
}

func HandleDeletePostById(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err := db.DeletePost(id)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
