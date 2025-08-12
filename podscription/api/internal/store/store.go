package store

import (
	"github.com/google/uuid"
	"podscription-api/types"
)

// Store defines the interface for session storage
type Store interface {
	CreateSession(name string) (*types.Session, error)
	GetSession(id uuid.UUID) (*types.Session, error)
	UpdateSession(session *types.Session) error
	ListSessions() ([]*types.Session, error)
	AddMessage(sessionID uuid.UUID, message types.Message) error
}