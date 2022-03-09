package user_controller

import (
	"html/template"
	"log"
	. "login/pkg/password"
	. "login/pkg/user/entity"
	. "login/pkg/user/service"
	"net/http"
	"os"
)

type UserController struct {
	userService *UserService
}

func (u *UserController) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", u.loginHandler)
	mux.HandleFunc("/signup", u.signupHandler)
	mux.HandleFunc("/blocked", u.blockedHandler)
	mux.HandleFunc("/succes", u.succesHandler)

	return mux
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
	var tmpl = template.Must(template.ParseFiles("web/template/login.html"))

	if r.Method == "POST" {
		uname := r.FormValue("uname")
		pass := r.FormValue("pass")

		//saved := db.FindUser(uname)
		saved := u.userService.FindUser(uname)
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

func (u *UserController) signupHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("web/template/signup.html"))

	if r.Method == "POST" {
		uname := r.FormValue("uname")
		hash, _ := HashPassword(r.FormValue("pass"))

		data := User{
			Username: uname,
			Password: hash,
		}

		u.userService.NewUser(data)
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

func NewController(service *UserService) *UserController {
	return &UserController{userService: service}
}
