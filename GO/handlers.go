package forum

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID              int    `json:"id"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Genre           string `json:"genre"`
	Nom             string `json:"nom"`
	Prenom          string `json:"prenom"`
	DateNaissance   string `json:"dateNaissance"`
	Telephone       string `json:"telephone"`
}

var templates = template.Must(template.ParseFiles(
	"../templates/home.html",
	"../templates/register.html",
	"../templates/login.html",
	"../templates/create-post.html",
	"../templates/posts.html",
	"../templates/comments.html",
))

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "../templates/home.html", nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Initialiser le store de session
	var store = sessions.NewCookieStore([]byte("secret-key"))
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Save(r, w)

	session.Options.MaxAge = -1
	renderTemplate(w, "home.html", nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		renderTemplate(w, "register.html", nil)
		return
	}

	var user User
	user.Email = r.FormValue("email")
	user.Password = r.FormValue("password")
	user.ConfirmPassword = r.FormValue("confirmPassword")
	user.Genre = r.FormValue("genre")
	user.Nom = r.FormValue("nom")
	user.Prenom = r.FormValue("prenom")
	user.DateNaissance = r.FormValue("dateNaissance")
	user.Telephone = r.FormValue("telephone")

	// Vérifier que les mots de passe correspondent
	if user.Password != user.ConfirmPassword {
		http.Error(w, "Les mots de passe ne correspondent pas", http.StatusBadRequest)
		return
	}

	// Vérification de l'âge
	birthDate, err := time.Parse("2006-01-02", user.DateNaissance)
	if err != nil {
		http.Error(w, "Date de naissance invalide", http.StatusBadRequest)
		return
	}

	age := time.Now().Sub(birthDate).Hours() / 24 / 365
	if age < 16 {
		http.Error(w, "Vous devez avoir au moins 16 ans pour vous inscrire", http.StatusBadRequest)
		return
	}

	// Vérification du numéro de téléphone
	phoneRegex := regexp.MustCompile(`^\d{10}$`)
	if !phoneRegex.MatchString(user.Telephone) {
		http.Error(w, "Numéro de téléphone invalide", http.StatusBadRequest)
		return
	}

	// Validation de la force du mot de passe
	if !isPasswordStrong(user.Password) {
		http.Error(w, "Le mot de passe doit contenir au moins 8 caractères, une lettre majuscule, une lettre minuscule, un chiffre et un caractère spécial", http.StatusBadRequest)
		return
	}

	// Hachage du mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	stmt, err := db.Prepare("INSERT INTO users (email, password, genre, nom, prenom, dateNaissance, telephone) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Email, user.Password, user.Genre, user.Nom, user.Prenom, user.DateNaissance, user.Telephone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "login.html", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Code pour gérer la connexion
	if r.Method != http.MethodPost {
		renderTemplate(w, "login.html", nil)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var storedUser User
	err = db.QueryRow("SELECT id, password FROM users WHERE email = ?", user.Email).Scan(&storedUser.ID, &storedUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate a new session UUID
	sessionID, err := uuid.NewRandom()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store the session UUID in the cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   sessionID.String(),
		Expires: time.Now().Add(24 * time.Hour),
	})

	w.WriteHeader(http.StatusOK)
	renderTemplate(w, "home.html", nil)
}

func isPasswordStrong(password string) bool {
	// Vérifie si le mot de passe a au moins 8 caractères
	if len(password) < 8 {
		return false
	}

	// Vérifie s'il contient au moins une lettre majuscule
	hasUppercase := false
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUppercase = true
			break
		}
	}
	if !hasUppercase {
		return false
	}

	// Vérifie s'il contient au moins une lettre minuscule
	hasLowercase := false
	for _, char := range password {
		if unicode.IsLower(char) {
			hasLowercase = true
			break
		}
	}
	if !hasLowercase {
		return false
	}

	// Vérifie s'il contient au moins un chiffre
	hasDigit := false
	for _, char := range password {
		if unicode.IsDigit(char) {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		return false
	}

	// Vérifie s'il contient au moins un caractère spécial
	hasSpecialChar := false
	specialChars := "!@#$%^&*()-_=+[]{}|;:',<.>/?"
	for _, char := range password {
		if strings.ContainsRune(specialChars, char) {
			hasSpecialChar = true
			break
		}
	}
	if !hasSpecialChar {
		return false
	}

	return true
}
