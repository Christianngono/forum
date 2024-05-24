package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type UserData struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func main() {
	http.HandleFunc("/", renderForm)
	http.HandleFunc("/CréerUnCompte", handleRegister)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func renderForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("CréerUnCompte.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var userData UserData
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	userData.FirstName = r.FormValue("firstName")
	userData.LastName = r.FormValue("lastName")
	userData.Email = r.FormValue("email")
	userData.Phone = r.FormValue("phone")
	userData.Password = r.FormValue("password")
	userData.ConfirmPassword = r.FormValue("confirmPassword")

	if userData.Password != userData.ConfirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	// Simulating local storage by just printing the data.
	jsonData, err := json.Marshal(userData)
	if err != nil {
		http.Error(w, "Unable to marshal JSON", http.StatusInternalServerError)
		return
	}
	log.Printf("User Data: %s\n", jsonData)

	// Redirect to login page
	http.Redirect(w, r, "/SeConnecter", http.StatusSeeOther)
}
