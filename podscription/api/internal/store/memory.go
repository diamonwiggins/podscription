package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"podscription-api/types"
)

// MemoryStore implements an in-memory store with optional file persistence
type MemoryStore struct {
	sessions map[uuid.UUID]*types.Session
	mu       sync.RWMutex
	filePath string
}

// NewMemoryStore creates a new memory store
func NewMemoryStore(filePath string) *MemoryStore {
	store := &MemoryStore{
		sessions: make(map[uuid.UUID]*types.Session),
		filePath: filePath,
	}

	// Try to load existing data
	if filePath != "" {
		store.loadFromFile()
	}

	return store
}

// CreateSession creates a new session
func (s *MemoryStore) CreateSession(name string) (*types.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if name == "" {
		name = fmt.Sprintf("Session %d", len(s.sessions)+1)
	}

	session := &types.Session{
		ID:        uuid.New(),
		Name:      name,
		Messages:  make([]types.Message, 0),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.sessions[session.ID] = session
	s.saveToFile()
	return session, nil
}

// GetSession retrieves a session by ID
func (s *MemoryStore) GetSession(id uuid.UUID) (*types.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[id]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", id)
	}

	// Return a copy to prevent external modification
	sessionCopy := *session
	sessionCopy.Messages = make([]types.Message, len(session.Messages))
	copy(sessionCopy.Messages, session.Messages)

	return &sessionCopy, nil
}

// UpdateSession updates an existing session
func (s *MemoryStore) UpdateSession(session *types.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.sessions[session.ID]; !exists {
		return fmt.Errorf("session not found: %s", session.ID)
	}

	session.UpdatedAt = time.Now()
	s.sessions[session.ID] = session
	s.saveToFile()
	return nil
}

// ListSessions returns all sessions
func (s *MemoryStore) ListSessions() ([]*types.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sessions := make([]*types.Session, 0, len(s.sessions))
	for _, session := range s.sessions {
		// Return copies to prevent external modification
		sessionCopy := *session
		sessionCopy.Messages = make([]types.Message, len(session.Messages))
		copy(sessionCopy.Messages, session.Messages)
		sessions = append(sessions, &sessionCopy)
	}

	return sessions, nil
}

// AddMessage adds a message to a session
func (s *MemoryStore) AddMessage(sessionID uuid.UUID, message types.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	message.ID = uuid.New()
	message.Timestamp = time.Now()
	session.Messages = append(session.Messages, message)
	session.UpdatedAt = time.Now()

	s.saveToFile()
	return nil
}

// loadFromFile loads sessions from file if it exists
func (s *MemoryStore) loadFromFile() {
	if s.filePath == "" {
		return
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		// File doesn't exist or can't be read, start fresh
		return
	}

	var sessions map[uuid.UUID]*types.Session
	if err := json.Unmarshal(data, &sessions); err != nil {
		// Invalid data, start fresh
		return
	}

	s.sessions = sessions
}

// saveToFile saves sessions to file
func (s *MemoryStore) saveToFile() {
	if s.filePath == "" {
		return
	}

	// Create directory if it doesn't exist
	if dir := filepath.Dir(s.filePath); dir != "." {
		os.MkdirAll(dir, 0755)
	}

	data, err := json.MarshalIndent(s.sessions, "", "  ")
	if err != nil {
		return
	}

	os.WriteFile(s.filePath, data, 0644)
}