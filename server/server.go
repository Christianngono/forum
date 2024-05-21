package main

import (
	"log"
	"net/http"
	"forum"
)

func main() {
	// Initialiser les tables de la base de données
	forum.CreateTables()
    defer forum.Close()

    http.HandleFunc("/register", forum.RegisterHandler)
    http.HandleFunc("/login", forum.LoginHandler)
    http.HandleFunc("/create-post", forum.CreatePostHandler)
    http.HandleFunc("/get-posts", forum.GetPostsHandler)
    http.HandleFunc("/create-comment", forum.CreateCommentHandler)
    http.HandleFunc("/get-comments", forum.GetCommentsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
	log.Println("Server started on port 8080")
}
