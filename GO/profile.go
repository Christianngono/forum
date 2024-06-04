// handlers/profile.go
package forum

import (
	"database/sql"
	forum "forum/GO"
	"html/template"
	"net/http"
	"strconv"
)

func ProfileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Pour cet exemple, nous utilisons un ID d'utilisateur statique.
		// Dans une application réelle, vous devrez récupérer l'ID de l'utilisateur connecté.
		userID, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		user, err := forum.GetUserByID(db, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("templates/profile.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, user)
	}
}
