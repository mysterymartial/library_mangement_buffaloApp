package repository

import (
	"database/sql"
	"errors"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"library-system/models"
	"log"
)

type UserRepository interface {
	AddUser(user *models.User) error
	GetUserByID(ID uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

type UserRepositoryImpl struct {
	DB *pop.Connection
}

func NewUserRepository(db *pop.Connection) *UserRepositoryImpl {
	db, err := pop.Connect("development")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	} else {
		log.Println("Successfully connected to database")
	}
	return &UserRepositoryImpl{DB: db}
}

func (r *UserRepositoryImpl) AddUser(user *models.User) error {
	return r.DB.Transaction(func(tx *pop.Connection) error {
		existingUser := &models.User{}
		err := tx.Where("email = ?", user.Email).First(existingUser)
		if err == nil {
			return errors.New("email already exists")
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return tx.Create(user)
	})
}

func (r *UserRepositoryImpl) GetUserByID(ID uuid.UUID) (*models.User, error) {
	user := &models.User{}
	err := r.DB.Where("id = ?", ID).First(user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryImpl) GetUserByEmail(email string) (*models.User, error) {
	log.Printf("Querying user by email: %s", email)
	user := &models.User{}
	err := r.DB.Where("email = ?", email).First(user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No user found with email: %s", email)
			return nil, nil
		}
		log.Printf("Error querying user by email: %v", err)
		return nil, err
	}
	log.Printf("User found: %v", user)
	return user, nil
}
