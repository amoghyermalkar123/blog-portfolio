// internal/middleware/middleware.go
package middleware

import (
	"blog-portfolio/internal/logger"
	"net/http"
	"time"
)

type Middleware struct {
	logger *logger.Logger
}

func New(logger *logger.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}

// RequestLogger logs each request
func (m *Middleware) RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log the request
		m.logger.Info(
			"Method:", r.Method,
			"Path:", r.URL.Path,
			"Duration:", time.Since(start),
		)
	})
}

// Recover middleware handles panics
func (m *Middleware) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.logger.Error("Recovered from panic:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
