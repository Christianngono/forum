package main

import (
	"log"
	"net/http"

	"forum"
)

func main() {
	// Initialiser la base de données et créer les tables nécessaires
	forum.CreateTables()

	// Définir les routes et leurs handlers correspondants
	http.HandleFunc("/register", forum.RegisterHandler)
	http.HandleFunc("/login", forum.LoginHandler)
	http.HandleFunc("/create-post", forum.CreatePostHandler)
	http.HandleFunc("/get-posts", forum.GetPostsHandler)
	http.HandleFunc("/create-comment", forum.CreateCommentHandler)
	http.HandleFunc("/get-comments", forum.GetCommentsHandler)

	// Démarrer le serveur web
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
