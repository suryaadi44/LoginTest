package main

import (
	"log"
	. "login/pkg/database"
	. "login/pkg/password"
	"net/http"
	"os"
	"text/template"

	"github.com/joho/godotenv"
)

var db = Database{}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("[Start]", os.Getenv("APP_NAME"))
}

func main() {
	db.DBConnect()

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/blocked", blockedHandler)
	http.HandleFunc("/succes", succesHandler)
	host := os.Getenv("HOST")

	log.Printf("[Start] Server started at http://%s", host)
	http.ListenAndServe(host, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/login.html"))

	if r.Method == "POST" {
		uname := r.FormValue("uname")
		pass := r.FormValue("pass")

		saved := db.FindUser(uname)
		if saved.Username == uname && CheckPasswordHash(pass, saved.Password) {
			log.Println("[Login Succes]", uname, "login approved")
			http.Redirect(w, r, "/succes", http.StatusSeeOther)
			return
		}

		log.Println("[Login Failed]", uname, "login failed")
		http.Redirect(w, r, "/blocked", http.StatusSeeOther)
		return
	}

	var err = tmpl.Execute(w, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/signup.html"))

	if r.Method == "POST" {
		uname := r.FormValue("uname")
		hash, _ := HashPassword(r.FormValue("pass"))

		data := User{
			Username: uname,
			Password: hash,
		}

		db.NewUser(data)
		log.Println("[Sign Up Succes]", uname, "created")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	var err = tmpl.Execute(w, nil)

	if err != nil {
		log.Println("[Sign Up Failed] login failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func succesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/template/succes.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func blockedHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/template/blocked.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
