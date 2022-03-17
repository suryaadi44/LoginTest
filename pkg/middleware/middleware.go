package middleware

import (
	"net/http"

	. "login/pkg/session/service"

	"github.com/gorilla/mux"
)

type Middleware struct {
	session *SessionService
}

func (m Middleware) AuthMiddleware() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := &http.Cookie{}

			if storedCookie, _ := r.Cookie("session_token"); storedCookie != nil {
				c = storedCookie
			}

			if c.Value == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return

			}

			_, err := m.session.FindSession(c.Value)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

func NewMiddlewareService(session *SessionService) *Middleware {
	return &Middleware{session: session}
}
