package store

// Store defines the interface for all storage operations
// This can be extended in the future to include other data stores
type Store interface {
	Users() UserStore
	Tokens() TokenStore
}

// InMemoryStore implements Store with in-memory storage
type InMemoryStore struct {
	userStore  UserStore
	tokenStore TokenStore
}

// NewInMemoryStore creates a new instance of InMemoryStore
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		userStore:  NewInMemoryUserStore(),
		tokenStore: NewInMemoryTokenStore(),
	}
}

// Users returns the user store
func (s *InMemoryStore) Users() UserStore {
	return s.userStore
}

// Tokens returns the token store
func (s *InMemoryStore) Tokens() TokenStore {
	return s.tokenStore
}
