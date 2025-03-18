package store

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sanskarm98/auth-service/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// UserStore defines the interface for user data operations
type UserStore interface {
	Create(email, password string) (models.User, error)
	GetByID(id string) (models.User, bool)
	GetByEmail(email string) (models.User, bool)
	Authenticate(email, password string) (models.User, bool)
}

// InMemoryUserStore implements UserStore with in-memory storage
type InMemoryUserStore struct {
	users      map[string]models.User
	usersMutex sync.RWMutex
}

// NewInMemoryUserStore creates a new instance of InMemoryUserStore
func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: make(map[string]models.User),
	}
}

// Create adds a new user to the store
func (s *InMemoryUserStore) Create(email, password string) (models.User, error) {
	// Check if email already exists
	s.usersMutex.RLock()
	for _, user := range s.users {
		if user.Email == email {
			s.usersMutex.RUnlock()
			return models.User{}, errors.New(models.ErrEmailAlreadyExists)
		}
	}
	s.usersMutex.RUnlock()

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, errors.New(models.ErrInternalServerError)
	}

	// Create user
	user := models.User{
		ID:        uuid.New().String(),
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	// Store user
	s.usersMutex.Lock()
	s.users[user.ID] = user
	s.usersMutex.Unlock()

	return user, nil
}

// GetByID retrieves a user by ID
func (s *InMemoryUserStore) GetByID(id string) (models.User, bool) {
	s.usersMutex.RLock()
	defer s.usersMutex.RUnlock()

	user, exists := s.users[id]
	return user, exists
}

// GetByEmail retrieves a user by email
func (s *InMemoryUserStore) GetByEmail(email string) (models.User, bool) {
	s.usersMutex.RLock()
	defer s.usersMutex.RUnlock()

	for _, user := range s.users {
		if user.Email == email {
			return user, true
		}
	}

	return models.User{}, false
}

// Authenticate verifies user credentials and returns the user if valid
func (s *InMemoryUserStore) Authenticate(email, password string) (models.User, bool) {
	user, found := s.GetByEmail(email)
	if !found {
		return models.User{}, false
	}

	// Validate password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return models.User{}, false
	}

	return user, true
}
