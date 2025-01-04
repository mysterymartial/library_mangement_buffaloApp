package mock

import (
	"errors"
	"github.com/gofrs/uuid"
	"library-system/models"
)

type MockUserRepo struct {
	MockUser            []models.User
	AddUserError        error
	GetUserByIDError    error
	GetUserByEmailError error
	UpdateUserError     error
	DeleteUserError     error
}

func (r *MockUserRepo) AddUser(user *models.User) error {
	if r.AddUserError != nil {
		return r.AddUserError
	}

	for _, existingUser := range r.MockUser {
		if existingUser.Email == user.Email {
			return errors.New("email already exists")
		}
	}
	r.MockUser = append(r.MockUser, *user)
	return nil
}

func (r *MockUserRepo) GetUserByID(userID uuid.UUID) (*models.User, error) {
	if r.GetUserByIDError != nil {
		return nil, r.GetUserByIDError
	}
	for _, user := range r.MockUser {
		if user.ID == userID {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *MockUserRepo) GetUserByEmail(email string) (*models.User, error) {
	if r.GetUserByEmailError != nil {
		return nil, r.GetUserByEmailError
	}
	for _, user := range r.MockUser {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, nil
}
