package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gofrs/uuid"
	"log"
	"net/http"
)

const sessionName = "_library_session"

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

func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		session := c.Session()
		userID := session.Get("current_user_id")

		if userID == nil {
			return c.Render(http.StatusUnauthorized, render.JSON(map[string]string{
				"error": "Authentication required",
			}))
		}

		c.Set("current_user", userID)
		return next(c)
	}
}

func SecurityHeaders(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		c.Response().Header().Set("X-Frame-Options", "DENY")
		c.Response().Header().Set("X-Content-Type-Options", "nosniff")
		c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
		return next(c)
	}
}

func CORS(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		// Set CORS headers for all responses
		c.Response().Header().Set("Access-Control-Allow-Origin", "http://localhost:63342")
		c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-CSRF-Token")
		c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
		c.Response().Header().Set("Access-Control-Max-Age", "300")

		// Handle preflight requests
		if c.Request().Method == http.MethodOptions {
			return c.Render(http.StatusOK, render.String(""))
		}

		return next(c)
	}
}
