package main

import (
	"log"
	"net/http"

	"practice5/database"
	"practice5/handler"
	"practice5/repository"
)

func main() {

	db := database.InitDB()

	repo := repository.Repository{
		DB: db,
	}

	h := handler.Handler{
		Repo: &repo,
	}

	http.HandleFunc("/users", h.GetUsers)

	http.HandleFunc("/common-friends", h.GetCommonFriends)

	log.Println("Server running on :8080")

	http.ListenAndServe(":8080", nil)
}