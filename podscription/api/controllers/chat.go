package controllers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"podscription-api/internal/managers"
	"podscription-api/types"
)

// ChatController handles chat-related business logic
type ChatController struct {
	sessionManager *managers.SessionManager
	logger         *logrus.Logger
}

// NewChatController creates a new chat controller
func NewChatController(sessionManager *managers.SessionManager, logger *logrus.Logger) *ChatController {
	return &ChatController{
		sessionManager: sessionManager,
		logger:         logger,
	}
}

// SendMessage processes a chat message and returns the response
func (c *ChatController) SendMessage(ctx context.Context, req types.ChatRequest) (*types.ChatResponse, error) {
	// Validate request
	if req.Content == "" {
		return nil, &types.ErrorResponse{
			ErrorCode: "INVALID_REQUEST",
			Message:   "Message content cannot be empty",
		}
	}

	var sessionID uuid.UUID
	var err error

	// If no session ID provided, create a new session
	if req.SessionID == nil {
		session, err := c.sessionManager.CreateSession("")
		if err != nil {
			c.logger.WithError(err).Error("failed to create new session for chat")
			return nil, &types.ErrorResponse{
				ErrorCode: "SESSION_CREATION_FAILED",
				Message:   "Failed to create chat session",
			}
		}
		sessionID = session.ID
		c.logger.WithField("session_id", sessionID).Info("created new session for chat")
	} else {
		sessionID = *req.SessionID
	}

	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Process the message
	session, message, err := c.sessionManager.ProcessMessage(ctx, sessionID, req.Content)
	if err != nil {
		c.logger.WithFields(logrus.Fields{
			"session_id": sessionID,
			"error":      err,
		}).Error("failed to process chat message")

		// Check if it's a session not found error
		if err.Error() == "session not found" {
			return nil, &types.ErrorResponse{
				ErrorCode: "SESSION_NOT_FOUND",
				Message:   "Chat session not found",
			}
		}

		return nil, &types.ErrorResponse{
			ErrorCode: "PROCESSING_FAILED",
			Message:   "Failed to process message",
		}
	}

	response := &types.ChatResponse{
		Session: *session,
		Message: *message,
	}

	c.logger.WithFields(logrus.Fields{
		"session_id":     sessionID,
		"message_count":  len(session.Messages),
		"intent_category": message.Intent.Category,
	}).Info("successfully processed chat message")

	return response, nil
}

// CreateSession creates a new chat session
func (c *ChatController) CreateSession(req types.CreateSessionRequest) (*types.Session, error) {
	session, err := c.sessionManager.CreateSession(req.Name)
	if err != nil {
		c.logger.WithError(err).Error("failed to create session")
		return nil, &types.ErrorResponse{
			ErrorCode: "SESSION_CREATION_FAILED",
			Message:   "Failed to create session",
		}
	}

	return session, nil
}

// GetSession retrieves a session by ID
func (c *ChatController) GetSession(sessionID uuid.UUID) (*types.Session, error) {
	session, err := c.sessionManager.GetSession(sessionID)
	if err != nil {
		c.logger.WithFields(logrus.Fields{
			"session_id": sessionID,
			"error":      err,
		}).Error("failed to get session")
		
		return nil, &types.ErrorResponse{
			ErrorCode: "SESSION_NOT_FOUND",
			Message:   "Session not found",
		}
	}

	return session, nil
}

// ListSessions returns all sessions
func (c *ChatController) ListSessions() ([]*types.Session, error) {
	sessions, err := c.sessionManager.ListSessions()
	if err != nil {
		c.logger.WithError(err).Error("failed to list sessions")
		return nil, &types.ErrorResponse{
			ErrorCode: "LISTING_FAILED",
			Message:   "Failed to retrieve sessions",
		}
	}

	return sessions, nil
}