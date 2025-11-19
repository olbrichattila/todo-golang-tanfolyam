package crud

import (
	"net/http"

	"github.com/google/uuid"
)

const (
	sessionCookieName = "todo_session"
)

func (s *service) sessionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			if err != http.ErrNoCookie {
				http.Error(w, "cookie read error", http.StatusInternalServerError)
				return
			}
		}

		if cookie == nil {
			s.sessionID = uuid.New().String()

			http.SetCookie(w, &http.Cookie{
				Name:  sessionCookieName,
				Value: s.sessionID,
			})

			next(w, r)
			return

		}

		s.sessionID = cookie.Value

		next(w, r)
	}
}
