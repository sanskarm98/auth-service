package store

import (
	"sync"
)

// TokenStore defines the interface for token operations
type TokenStore interface {
	StoreRefreshToken(token, userID string)
	GetUserIDByRefreshToken(token string) (string, bool)
	DeleteRefreshToken(token string)
	IsTokenRevoked(token string) bool
	RevokeToken(token string)
}

// InMemoryTokenStore implements TokenStore with in-memory storage
type InMemoryTokenStore struct {
	refreshTokens     map[string]string // token -> userID
	revokedTokens     map[string]bool   // Blacklist for revoked tokens
	refreshTokenMutex sync.RWMutex
	revokedTokenMutex sync.RWMutex
}

// NewInMemoryTokenStore creates a new instance of InMemoryTokenStore
func NewInMemoryTokenStore() *InMemoryTokenStore {
	return &InMemoryTokenStore{
		refreshTokens: make(map[string]string),
		revokedTokens: make(map[string]bool),
	}
}

// StoreRefreshToken stores a refresh token with associated userID
func (s *InMemoryTokenStore) StoreRefreshToken(token, userID string) {
	s.refreshTokenMutex.Lock()
	defer s.refreshTokenMutex.Unlock()
	s.refreshTokens[token] = userID
}

// GetUserIDByRefreshToken retrieves the userID associated with a refresh token
func (s *InMemoryTokenStore) GetUserIDByRefreshToken(token string) (string, bool) {
	s.refreshTokenMutex.RLock()
	defer s.refreshTokenMutex.RUnlock()
	userID, exists := s.refreshTokens[token]
	return userID, exists
}

// DeleteRefreshToken removes a refresh token from the store
func (s *InMemoryTokenStore) DeleteRefreshToken(token string) {
	s.refreshTokenMutex.Lock()
	defer s.refreshTokenMutex.Unlock()
	delete(s.refreshTokens, token)
}

// IsTokenRevoked checks if a token has been revoked
func (s *InMemoryTokenStore) IsTokenRevoked(token string) bool {
	s.revokedTokenMutex.RLock()
	defer s.revokedTokenMutex.RUnlock()
	_, revoked := s.revokedTokens[token]
	return revoked
}

// RevokeToken adds a token to the revoked list
func (s *InMemoryTokenStore) RevokeToken(token string) {
	s.revokedTokenMutex.Lock()
	defer s.revokedTokenMutex.Unlock()
	s.revokedTokens[token] = true
}
