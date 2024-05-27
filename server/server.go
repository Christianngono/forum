package main

import (
	"forum"
	"log"
	"net/http"
)

func main() {
	// Initialiser les tables de la base de données
	forum.CreateTables()

	http.HandleFunc("/register", forum.RegisterHandler)
	http.HandleFunc("/login", forum.LoginHandler)
	http.HandleFunc("/create-post", forum.CreatePostHandler)
	http.HandleFunc("/get-posts", forum.GetPostsHandler)
	http.HandleFunc("/create-comment", forum.CreateCommentHandler)
	http.HandleFunc("/get-comments", forum.GetCommentsHandler)
	http.HandleFunc("/like-post", forum.LikePostHandler)
    http.HandleFunc("/dislike-post", forum.DislikePostHandler)
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
