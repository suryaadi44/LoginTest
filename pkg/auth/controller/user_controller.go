package controller

import (
	"html/template"
	"log"
	. "login/pkg/auth/entity"
	. "login/pkg/auth/service"
	. "login/pkg/middleware"
	. "login/pkg/password"
	. "login/pkg/session/entity"
	. "login/pkg/session/service"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserController struct {
	userService       *UserService
	sessionService    *SessionService
	middlewareService *Middleware
}

func (u *UserController) Handler() http.Handler {
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("web/static/"))))

	r.HandleFunc("/login", u.loginHandler)
	r.HandleFunc("/logout", u.logoutHandler)
	r.HandleFunc("/signup", u.signupHandler)
	r.HandleFunc("/blocked", u.blockedHandler)

	getContent := r.PathPrefix("/succes").Subrouter()
	getContent.Use(u.middlewareService.AuthMiddleware())
	getContent.HandleFunc("", u.succesHandler)

	return r
}

func (u *UserController) Run() {
	host := os.Getenv("HOST")

	httpServer := &http.Server{
		Addr:    host,
		Handler: u.Handler(),
	}

	log.Printf("[Start] Server started at http://%s", host)
	httpServer.ListenAndServe()
}

func (u *UserController) loginHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/login/index.html"))

	if r.Method == "POST" {
		uname := r.FormValue("username")
		pass := r.FormValue("password")

		// TODO : Add error handling
		saved, _ := u.userService.FindUser(uname)
		if saved.Username == uname && CheckPasswordHash(pass, saved.Password) {
			log.Println("[Login Succes]", uname, "login approved")

			session := Session{
				SessionToken: uuid.NewString(),
				Username:     uname,
				Expire:       time.Now().Add(time.Duration(SESSION_EXPIRE_IN_SECOND) * time.Second),
			}

			u.sessionService.NewSession(session)

			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   session.SessionToken,
				Expires: session.Expire,
			})

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

func (u *UserController) signupHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/signup/index.html"))

	if r.Method == "POST" {
		uname := r.FormValue("username")
		hash, _ := HashPassword(r.FormValue("password"))

		data := User{
			Username: uname,
			Password: hash,
		}

		err := u.userService.NewUser(data)
		if err == nil {
			log.Println("[Sign Up Succes]", uname, "created")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if !strings.Contains(err.Error(), "duplicate") {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("[Sign Up Failed]", uname, "Already exist")
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return
	}

	var err = tmpl.Execute(w, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (u UserController) logoutHandler(w http.ResponseWriter, r *http.Request) {
	if storedCookie, _ := r.Cookie("session_token"); storedCookie != nil {
		u.sessionService.DeleteSession(storedCookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		MaxAge:  -1,
		Expires: time.Unix(0, 0),
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (u *UserController) succesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/template/succes.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (u *UserController) blockedHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/template/blocked.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func NewController(userService *UserService, sessionService *SessionService, middlewareService *Middleware) *UserController {
	return &UserController{userService: userService, sessionService: sessionService, middlewareService: middlewareService}
}
