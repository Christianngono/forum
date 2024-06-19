package main

import (
	"forum"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading.env file")
	} else {
		log.Println("Loaded.env file")
	}
}

func main() {
	// Initialiser les tables de la base de données
	forum.InitDB()
	// Fermer la base de données à la fin de l'exécution
	defer forum.DB.Close()

	// Servir les fichiers statiques du répertoire 'static'
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))

	http.HandleFunc("/logout", forum.LogoutHandler)
	http.HandleFunc("/", forum.IndexHandler)
	http.HandleFunc("/register", forum.RegisterHandler)
	http.HandleFunc("/login", forum.LoginHandler)
	http.HandleFunc("/create-post", forum.CreatePostHandler)
	http.HandleFunc("/edit-post", forum.EditPostHandler)
	http.HandleFunc("/update-post", forum.UpdatePostHandler)
	http.HandleFunc("/filter-post", forum.FilterPostHandler)
	http.HandleFunc("/delete-post", forum.DeletePostHandler)
	http.HandleFunc("/get-dislike-post", forum.DislikePostHandler)
	http.HandleFunc("/create-comment", forum.CreateCommentHandler)
	http.HandleFunc("/edit-comment", forum.EditCommentHandler)
	http.HandleFunc("/filter-comment", forum.FilterCommentHandler)
	http.HandleFunc("/update-comment", forum.UpdateCommentHandler)
	http.HandleFunc("/delete-comment", forum.DeleteCommentHandler)
	http.HandleFunc("/get-like-post", forum.LikePostHandler)
	http.HandleFunc("/get-dislike-post", forum.DislikePostHandler)
	http.HandleFunc("/like-post", forum.LikePostHandler)
	http.HandleFunc("/get-posts", forum.GetPostsHandler)
	http.HandleFunc("/get-post", forum.GetPostHandler)
    http.HandleFunc("/get-comments", forum.GetCommentsHandler)
	http.HandleFunc("/get-comment", forum.GetCommentHandler)
	


	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
