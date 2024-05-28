package main

import (
	"forum"
	"log"
	"net/http"
)

func main() {
	// Initialiser les tables de la base de données
	forum.CreateTables()

	http.HandleFunc("templates/register", forum.RegisterHandler)
	http.HandleFunc("templates/login", forum.LoginHandler)
	http.HandleFunc("templates/create-post", forum.CreatePostHandler)
	http.HandleFunc("templates/get-posts", forum.GetPostsHandler)
	http.HandleFunc("templates/create-comment", forum.CreateCommentHandler)
	http.HandleFunc("templates/get-comments", forum.GetCommentsHandler)
	http.HandleFunc("templates/like-post", forum.LikePostHandler)
	http.HandleFunc("templates/dislike-post", forum.DislikePostHandler)
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
