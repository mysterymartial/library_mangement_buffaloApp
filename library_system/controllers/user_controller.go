package controllers

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gorilla/sessions"
	"library-system/Dto"
	"library-system/services"
	"log"
	"net/http"
	"strings"
)

const (
	sessionName = "_library_session"
	userIDKey   = "current_user_id"
)

type UserController struct {
	UserService  *services.UserServices
	SessionStore sessions.Store
}

func NewUserController(userService *services.UserServices, sessionStore sessions.Store) *UserController {
	return &UserController{
		UserService:  userService,
		SessionStore: sessionStore,
	}
}

func (uc *UserController) RegisterUser(c buffalo.Context) error {
	var request Dto.UserRequest
	if err := c.Bind(&request); err != nil {
		return c.Render(http.StatusBadRequest, render.JSON(map[string]string{
			"error": "Invalid request format",
		}))
	}

	request.Email = normalizeEmail(request.Email)
	user, err := uc.UserService.RegisterUser(request)
	if err != nil {
		if err.Error() == "email already registered" {
			return c.Render(http.StatusConflict, render.JSON(map[string]string{
				"error": "Email already registered",
			}))
		}
		return c.Render(http.StatusBadRequest, render.JSON(map[string]string{
			"error": err.Error(),
		}))
	}

	session := c.Session()
	session.Set(userIDKey, user.ID.String())
	session.Save()

	return c.Render(http.StatusOK, render.JSON(map[string]interface{}{
		"status": "success",
		"user":   user,
	}))
}

func (uc *UserController) CheckoutBook(c buffalo.Context) error {
	var request Dto.BookActionRequest
	if err := c.Bind(&request); err != nil {
		return c.Render(http.StatusBadRequest, render.JSON(map[string]string{
			"error": "Invalid request format",
		}))
	}

	request.Email = normalizeEmail(request.Email)
	log.Printf("Processing checkout - Book ID: %s, Email: %s", request.BookID, request.Email)

	response, err := uc.UserService.CheckOutBook(request)
	if err != nil {
		log.Printf("Checkout failed: %v", err)
		statusCode := http.StatusBadRequest
		if strings.Contains(err.Error(), "not available") {
			statusCode = http.StatusConflict
		}
		return c.Render(statusCode, render.JSON(map[string]string{
			"error": err.Error(),
		}))
	}

	return c.Render(http.StatusOK, render.JSON(map[string]interface{}{
		"status":   "success",
		"checkout": response,
	}))
}

func (uc *UserController) ReturnBook(c buffalo.Context) error {
	var request Dto.BookActionRequest
	if err := c.Bind(&request); err != nil {
		return c.Render(http.StatusBadRequest, render.JSON(map[string]string{
			"error": "Invalid request format",
		}))
	}

	request.Email = normalizeEmail(request.Email)
	log.Printf("Processing return - Book ID: %s, Email: %s", request.BookID, request.Email)

	response, err := uc.UserService.ReturnBook(request)
	if err != nil {
		log.Printf("Return failed: %v", err)
		statusCode := http.StatusBadRequest
		if strings.Contains(err.Error(), "No active loan found") {
			statusCode = http.StatusNotFound
		}
		return c.Render(statusCode, render.JSON(map[string]string{
			"error": err.Error(),
		}))
	}

	return c.Render(http.StatusOK, render.JSON(map[string]interface{}{
		"status": "success",
		"return": response,
	}))
}

func (uc *UserController) ReserveBook(c buffalo.Context) error {
	var request Dto.BookActionRequest
	if err := c.Bind(&request); err != nil {
		return c.Render(http.StatusBadRequest, render.JSON(map[string]string{
			"error": "Invalid request format",
		}))
	}

	request.Email = normalizeEmail(request.Email)
	log.Printf("Processing reservation - Book ID: %s, Email: %s", request.BookID, request.Email)

	response, err := uc.UserService.ReserveBook(request)
	if err != nil {
		log.Printf("Reservation failed: %v", err)
		statusCode := http.StatusBadRequest
		if strings.Contains(err.Error(), "not available") {
			statusCode = http.StatusConflict
		}
		return c.Render(statusCode, render.JSON(map[string]string{
			"error": err.Error(),
		}))
	}

	return c.Render(http.StatusOK, render.JSON(map[string]interface{}{
		"status":      "success",
		"reservation": response,
	}))
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
