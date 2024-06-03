package main

import (
	"forum"
	"log"
	"net/http"
)

func server() {
	// Initialiser les tables de la base de données
	forum.CreateTables()

	// Servir les fichiers statiques du répertoire 'static'
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/logout", forum.LogoutHandler)
	http.HandleFunc("/", forum.IndexHandler)
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
