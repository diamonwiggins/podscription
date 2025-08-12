package managers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"podscription-api/internal/store"
	"podscription-api/types"
)

// SessionManager handles session-related operations
type SessionManager struct {
	store       store.Store
	openAI      *OpenAIManager
	logger      *logrus.Logger
}

// NewSessionManager creates a new session manager
func NewSessionManager(store store.Store, openAI *OpenAIManager, logger *logrus.Logger) *SessionManager {
	return &SessionManager{
		store:  store,
		openAI: openAI,
		logger: logger,
	}
}

// CreateSession creates a new chat session
func (m *SessionManager) CreateSession(name string) (*types.Session, error) {
	session, err := m.store.CreateSession(name)
	if err != nil {
		m.logger.WithError(err).Error("failed to create session")
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	m.logger.WithFields(logrus.Fields{
		"session_id": session.ID,
		"name":       session.Name,
	}).Info("created new session")

	return session, nil
}

// GetSession retrieves a session by ID
func (m *SessionManager) GetSession(id uuid.UUID) (*types.Session, error) {
	session, err := m.store.GetSession(id)
	if err != nil {
		m.logger.WithFields(logrus.Fields{
			"session_id": id,
			"error":      err,
		}).Error("failed to get session")
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

// ListSessions returns all sessions
func (m *SessionManager) ListSessions() ([]*types.Session, error) {
	sessions, err := m.store.ListSessions()
	if err != nil {
		m.logger.WithError(err).Error("failed to list sessions")
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	return sessions, nil
}

// ProcessMessage processes a user message and generates an AI response
func (m *SessionManager) ProcessMessage(ctx context.Context, sessionID uuid.UUID, content string) (*types.Session, *types.Message, error) {
	// Get the session
	session, err := m.store.GetSession(sessionID)
	if err != nil {
		return nil, nil, fmt.Errorf("session not found: %w", err)
	}

	// Add the user message
	userMessage := types.Message{
		Role:    types.MessageRoleUser,
		Content: content,
	}

	if err := m.store.AddMessage(sessionID, userMessage); err != nil {
		return nil, nil, fmt.Errorf("failed to add user message: %w", err)
	}

	m.logger.WithFields(logrus.Fields{
		"session_id": sessionID,
		"content_length": len(content),
	}).Info("processing user message")

	// Classify the intent
	intent, err := m.openAI.ClassifyIntent(ctx, content)
	if err != nil {
		m.logger.WithError(err).Error("failed to classify intent")
		// Continue with a default intent rather than failing
		intent = &types.PodIntent{
			Category:   types.IntentCategoryGeneral,
			Confidence: 0.5,
			Symptoms:   []string{"unknown issue"},
		}
	}

	m.logger.WithFields(logrus.Fields{
		"session_id": sessionID,
		"intent_category": intent.Category,
		"confidence": intent.Confidence,
	}).Info("classified user intent")

	// Get recent message history for context
	recentHistory := m.getRecentHistory(session.Messages, 5)

	// Generate the diagnosis
	prescription, treatment, err := m.openAI.GenerateDiagnosis(ctx, content, intent, recentHistory)
	if err != nil {
		m.logger.WithError(err).Error("failed to generate diagnosis")
		return nil, nil, fmt.Errorf("failed to generate diagnosis: %w", err)
	}

	// Create the assistant message
	assistantMessage := types.Message{
		Role:         types.MessageRoleAssistant,
		Content:      treatment,
		Intent:       intent,
		Prescription: prescription,
	}

	// Add the assistant message
	if err := m.store.AddMessage(sessionID, assistantMessage); err != nil {
		return nil, nil, fmt.Errorf("failed to add assistant message: %w", err)
	}

	m.logger.WithFields(logrus.Fields{
		"session_id": sessionID,
		"diagnosis": prescription.Diagnosis,
		"commands_count": len(prescription.Commands),
	}).Info("generated diagnosis and response")

	// Get the updated session
	updatedSession, err := m.store.GetSession(sessionID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get updated session: %w", err)
	}

	// Return the last message (the assistant's response)
	lastMessage := updatedSession.Messages[len(updatedSession.Messages)-1]
	return updatedSession, &lastMessage, nil
}

// getRecentHistory returns the most recent messages for context
func (m *SessionManager) getRecentHistory(messages []types.Message, count int) []types.Message {
	if len(messages) <= count {
		return messages
	}
	return messages[len(messages)-count:]
}