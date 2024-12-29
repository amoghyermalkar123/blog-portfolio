// internal/handlers/auth.go
package handlers

import (
	"blog-portfolio/internal/logger"
	"blog-portfolio/internal/middleware"
	"blog-portfolio/web/pages"
	"net/http"
	"time"
)

type AuthHandlers struct {
	logger *logger.Logger
}

func NewAuthHandlers(logger *logger.Logger) *AuthHandlers {
	return &AuthHandlers{
		logger: logger,
	}
}

// ShowLogin handles displaying the login page
func (h *AuthHandlers) ShowLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := pages.Login().Render(r.Context(), w)
		if err != nil {
			h.logger.Error("Error rendering login page:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// HandleLogin processes the login form
func (h *AuthHandlers) HandleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// TODO: Replace with actual user validation from database
		if username == "admin" && password == "admin" {
			// Create token
			token, err := middleware.CreateToken(1, username, "admin")
			if err != nil {
				h.logger.Error("Error creating token:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Set cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    token,
				Path:     "/",
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})

			// Redirect to admin dashboard
			http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
			return
		}

		// Return to login page with error
		err := pages.Login().Render(r.Context(), w)
		if err != nil {
			h.logger.Error("Error rendering login page:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// HandleLogout logs out the user
func (h *AuthHandlers) HandleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Clear the session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    "",
			Path:     "/",
			Expires:  time.Now().Add(-1 * time.Hour),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})

		// Redirect to login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
