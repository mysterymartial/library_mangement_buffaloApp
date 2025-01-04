package controllers

import (
	"fmt"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gofrs/uuid"
	"github.com/gorilla/sessions"
	"library-system/Dto"
	"library-system/services"
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

	user, err := uc.UserService.RegisterUser(request)
	if err != nil {
		if err.Error() == "email already exists" {
			return c.Render(http.StatusConflict, render.JSON(map[string]string{
				"error": "Email already registered",
			}))
		}
		return c.Render(http.StatusBadRequest, render.JSON(map[string]string{
			"error": err.Error(),
		}))
	}

	session := c.Session()
	session.Set("current_user_id", user.ID.String())
	session.Save()

	return c.Render(http.StatusOK, render.JSON(map[string]interface{}{
		"status": "success",
		"user":   user,
	}))
}

func (uc *UserController) CheckoutBook(c buffalo.Context) error {

	session := c.Session()
	userIDStr := session.Get("current_user_id")

	if userIDStr == nil {
		return c.Render(http.StatusUnauthorized, render.JSON(map[string]string{
			"error": "Authentication required",
		}))
	}

	var request Dto.BookActionRequest
	if err := c.Bind(&request); err != nil {
		return c.Render(http.StatusBadRequest, render.JSON(map[string]string{
			"error": "Invalid request format",
		}))
	}

	if request.BookID == uuid.Nil {
		return c.Render(http.StatusBadRequest, render.JSON(map[string]string{
			"error": "Invalid book ID",
		}))
	}

	if strings.TrimSpace(request.UserName) == "" {
		return c.Render(http.StatusBadRequest, render.JSON(map[string]string{
			"error": "User name is required",
		}))
	}

	userID, err := uuid.FromString(userIDStr.(string))
	if err != nil {
		return c.Render(http.StatusUnauthorized, render.JSON(map[string]string{
			"error": "Invalid session",
		}))
	}

	request.UserID = userID
	response, err := uc.UserService.CheckOutBook(request)
	if err != nil {
		return c.Render(http.StatusBadRequest, render.JSON(map[string]string{
			"error": err.Error(),
		}))
	}

	return c.Render(http.StatusOK, render.JSON(map[string]interface{}{
		"status":   "success",
		"checkout": response,
	}))
}

func (uc *UserController) ReserveBook(c buffalo.Context) error {
	currentUserID, err := uc.getUserIDFromSession(c)
	if err != nil {
		return c.Render(http.StatusUnauthorized, render.JSON(map[string]string{
			"error": "Authentication required",
		}))
	}

	var request Dto.BookActionRequest
	if err := c.Bind(&request); err != nil {
		return c.Render(http.StatusBadRequest, render.JSON(map[string]string{
			"error": "Invalid request format",
		}))
	}

	request.UserID = currentUserID
	action, err := uc.UserService.ReserveBook(request)
	if err != nil {
		return c.Render(http.StatusBadRequest, render.JSON(map[string]string{
			"error": err.Error(),
		}))
	}

	return c.Render(http.StatusOK, render.JSON(map[string]interface{}{
		"status":      "success",
		"reservation": action,
	}))
}

func (uc *UserController) ReturnBook(c buffalo.Context) error {
	currentUserID, err := uc.getUserIDFromSession(c)
	if err != nil {
		return c.Render(http.StatusUnauthorized, render.JSON(map[string]string{
			"error": "Authentication required",
		}))
	}

	var request Dto.BookActionRequest
	if err := c.Bind(&request); err != nil {
		return c.Render(http.StatusBadRequest, render.JSON(map[string]string{
			"error": "Invalid request format",
		}))
	}

	if request.BookID == uuid.Nil {
		return c.Render(http.StatusBadRequest, render.JSON(map[string]string{
			"error": "Invalid book ID",
		}))
	}

	request.UserID = currentUserID
	action, err := uc.UserService.ReturnBook(request)
	if err != nil {
		return c.Render(http.StatusBadRequest, render.JSON(map[string]string{
			"error": err.Error(),
		}))
	}

	return c.Render(http.StatusOK, render.JSON(map[string]interface{}{
		"status": "success",
		"return": action,
	}))
}

func (uc *UserController) getUserIDFromSession(c buffalo.Context) (uuid.UUID, error) {
	session, err := uc.SessionStore.Get(c.Request(), sessionName)
	if err != nil {
		return uuid.Nil, fmt.Errorf("session retrieval failed")
	}

	userIDStr, ok := session.Values[userIDKey].(string)
	if !ok || userIDStr == "" {
		return uuid.Nil, fmt.Errorf("user not authenticated")
	}

	userID, err := uuid.FromString(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID format")
	}

	return userID, nil
}
