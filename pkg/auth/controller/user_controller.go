package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	. "github.com/suryaadi44/LoginTest/pkg/auth/entity"
	. "github.com/suryaadi44/LoginTest/pkg/auth/service"
	. "github.com/suryaadi44/LoginTest/pkg/dto"
	. "github.com/suryaadi44/LoginTest/pkg/middleware"
	. "github.com/suryaadi44/LoginTest/pkg/password"
	. "github.com/suryaadi44/LoginTest/pkg/session/entity"
	. "github.com/suryaadi44/LoginTest/pkg/session/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserController struct {
	router            *mux.Router
	userService       *UserService
	sessionService    *SessionService
	middlewareService *Middleware
}

func (u *UserController) InitializeController() {
	u.router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("web/static/"))))

	u.router.HandleFunc("/login", u.loginHandler)
	u.router.HandleFunc("/logout", u.logoutHandler)
	u.router.HandleFunc("/signup", u.signupHandler)
	u.router.HandleFunc("/blocked", u.blockedHandler)

	getContent := u.router.PathPrefix("/succes").Subrouter()
	getContent.Use(u.middlewareService.AuthMiddleware())
	getContent.HandleFunc("", u.succesHandler)
}

func (u *UserController) loginHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/login/index.html"))

	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		payload := Form{}

		if err := decoder.Decode(&payload); err != nil {
			log.Println("[DECODE] Error decoding JSON")
			NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
			return
		}

		// TODO : Add error handling
		saved, _ := u.userService.FindUser(payload.Username)
		if saved.Username == payload.Username && CheckPasswordHash(payload.Password, saved.Password) {
			log.Println("[Login Succes]", payload.Username, "login approved")

			session := Session{
				SessionToken: uuid.NewString(),
				Username:     payload.Username,
				Expire:       time.Now().Add(time.Duration(SESSION_EXPIRE_IN_SECOND) * time.Second),
			}

			u.sessionService.NewSession(session)

			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   session.SessionToken,
				Expires: session.Expire,
			})

			NewBaseResponse(http.StatusSeeOther, false, "/succes").SendResponse(&w)
			return
		}

		log.Println("[Login Failed]", payload.Username, "login failed")
		NewBaseResponse(http.StatusUnauthorized, true, "Inccorect Username or Password").SendResponse(&w)
		return
	}

	var err = tmpl.Execute(w, nil)

	if err != nil {
		NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}
}

func (u *UserController) signupHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/signup/index.html"))

	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		payload := Form{}

		if err := decoder.Decode(&payload); err != nil {
			log.Println("[DECODE] Error decoding JSON")
			NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
			return
		}

		hash, _ := HashPassword(payload.Password)
		data := User{
			Username: payload.Username,
			Password: hash,
		}

		err := u.userService.NewUser(data)
		if err == nil {
			log.Println("[Sign Up Succes]", payload.Username, "created")
			NewBaseResponse(http.StatusSeeOther, false, "/login").SendResponse(&w)
			return
		}

		if !strings.Contains(err.Error(), "duplicate") {
			NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
			return
		}

		log.Println("[Sign Up Failed]", payload.Username, "Already exist")
		NewBaseResponse(http.StatusOK, true, fmt.Sprintf("Accout with username %s already exist", payload.Username)).SendResponse(&w)
		return
	}

	var err = tmpl.Execute(w, nil)

	if err != nil {
		NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
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
		NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}
}

func (u *UserController) blockedHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/template/blocked.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		NewBaseResponse(http.StatusInternalServerError, true, err.Error()).SendResponse(&w)
		return
	}
}

func NewController(router *mux.Router, userService *UserService, sessionService *SessionService, middlewareService *Middleware) *UserController {
	return &UserController{router: router, userService: userService, sessionService: sessionService, middlewareService: middlewareService}
}
