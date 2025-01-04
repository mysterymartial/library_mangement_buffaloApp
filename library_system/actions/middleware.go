package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gofrs/uuid"
	"log"
	"net/http"
)

const sessionName = "_library_session"

// SetCurrentUser middleware - used to set current user from session
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		session := c.Session()
		if userIDStr := session.Get("current_user_id"); userIDStr != nil {
			userID, err := uuid.FromString(userIDStr.(string))
			if err == nil {
				c.Set("current_user", userID)
				log.Printf("User authenticated: %s", userID)
			}
		}
		return next(c)
	}
}

// Authorize middleware ensures the user is authenticated
func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		session := c.Session()
		userID := session.Get("current_user_id")

		if userID == nil {
			return c.Render(http.StatusUnauthorized, render.JSON(map[string]string{
				"error": "Authentication required",
			}))
		}

		// Set user context for downstream handlers
		c.Set("current_user", userID)
		return next(c)
	}
}

// SecurityHeaders middleware adds security headers
func SecurityHeaders(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		c.Response().Header().Set("X-Frame-Options", "DENY")
		c.Response().Header().Set("X-Content-Type-Options", "nosniff")
		c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
		return next(c)
	}
}

// CORS middleware handles Cross-Origin Resource Sharing
func CORS(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		// Allow specific origins
		origin := c.Request().Header.Get("Origin")
		if origin == "http://localhost:63342" || origin == "http://localhost:3000" {
			c.Response().Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Response().Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request().Method == http.MethodOptions {
			c.Response().WriteHeader(http.StatusNoContent)
			return nil
		}

		return next(c)
	}
}
