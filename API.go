package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

// Benutzername und Passwort für die Pseudo-Authentifizierung
const (
	fakeUsername = "admin"
	fakePassword = "adminpass"
)

func main() {
	var err error
	db, err = sql.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	createTable()

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/dashboard", dashboardHandler).Methods("GET")
	r.HandleFunc("/adduser", addUserHandler).Methods("POST")

	http.Handle("/", r)

	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}

func createTable() {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL
    )`)
	if err != nil {
		fmt.Println("Error creating table:", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Login</title>
    </head>
    <body>
        <h1>Login Page</h1>
        <form method="POST" action="/login">
            <label for="username">Username:</label>
            <input type="text" id="username" name="username" required><br><br>
            <label for="password">Password:</label>
            <input type="password" id="password" name="password" required><br><br>
            <input type="submit" value="Login">
        </form>

        <!-- Formular zum Hinzufügen eines Benutzers -->
        <h2>Add User</h2>
        <form method="POST" action="/adduser">
            <label for="new_username">New Username:</label>
            <input type="text" id="new_username" name="new_username" required><br><br>
            <label for="new_password">New Password:</label>
            <input type="password" id="new_password" name="new_password" required><br><br>
            <input type="submit" value="Add User">
        </form>
    </body>
    </html>
    `
	t, _ := template.New("index").Parse(tmpl)
	t.Execute(w, nil)
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	newUsername := r.FormValue("new_username")
	newPassword := r.FormValue("new_password")

	if err := addUser(newUsername, newPassword); err != nil {
		fmt.Println("Error adding user:", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func addUser(username, password string) error {
	_, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
	return err
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Überprüfe die Pseudo-Authentifizierung
	if username == fakeUsername && password == fakePassword {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Dashboard!")
}
