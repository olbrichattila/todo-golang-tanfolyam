package crud

import (
	"net/http"
)

func (s *service) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData, err := s.session.GetAll(s.sessionID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, ok := sessionData["user"]; ok {
			next(w, r)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
