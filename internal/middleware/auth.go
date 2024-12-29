// internal/middleware/auth.go
package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserContextKey    contextKey = "user"
	IsAdminContextKey contextKey = "is_admin"
)

type User struct {
	ID       int64
	Username string
	Role     string
}

// Add this helper function
func IsAdmin(r *http.Request) bool {
	if r == nil {
		return false
	}
	if admin, ok := r.Context().Value(IsAdminContextKey).(bool); ok {
		return admin
	}
	return false
}

// Modify the RequireAuth middleware
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for session cookie
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Validate JWT token
		token, err := validateToken(cookie.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    "",
				Path:     "/",
				Expires:  time.Now().Add(-1 * time.Hour),
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Check if user has admin role
		role, _ := claims["role"].(string)
		isAdmin := role == "admin"

		// Create user context
		user := &User{
			ID:       int64(claims["user_id"].(float64)),
			Username: claims["username"].(string),
			Role:     role,
		}

		// Add user and admin status to context
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		ctx = context.WithValue(ctx, IsAdminContextKey, isAdmin)

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validateToken validates the JWT token
func validateToken(tokenString string) (*jwt.Token, error) {
	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse the token
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		// Get secret key from environment
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "your-secret-key" // Default for development
		}

		return []byte(secret), nil
	})
}

// CreateToken creates a new JWT token for a user
func CreateToken(userID int64, username, role string) (string, error) {
	// Get secret key from environment
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key" // Default for development
	}

	// Create claims
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and return token
	return token.SignedString([]byte(secret))
}

// GetUserFromContext retrieves the user from the context
func GetUserFromContext(ctx context.Context) *User {
	user, ok := ctx.Value(UserContextKey).(*User)
	if !ok {
		return nil
	}
	return user
}
