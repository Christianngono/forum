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
	http.HandleFunc("HTML/CréerUnCompte", forum.RegisterHandler)
	http.HandleFunc("HTML/SeConnecter", forum.LoginHandler)
	http.HandleFunc("HTML/create-post", forum.CreatePostHandler)
	http.HandleFunc("HTML/get-posts", forum.GetPostsHandler)
	http.HandleFunc("HTML/create-comment", forum.CreateCommentHandler)
	http.HandleFunc("HTML/get-comments", forum.GetCommentsHandler)

	// Démarrer le serveur web
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
