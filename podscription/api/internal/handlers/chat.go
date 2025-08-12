package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"podscription-api/controllers"
	"podscription-api/types"
)

// ChatHandler handles HTTP requests for chat operations
type ChatHandler struct {
	controller *controllers.ChatController
	logger     *logrus.Logger
}

// NewChatHandler creates a new chat handler
func NewChatHandler(controller *controllers.ChatController, logger *logrus.Logger) *ChatHandler {
	return &ChatHandler{
		controller: controller,
		logger:     logger,
	}
}

// SendMessage handles POST /api/chat
func (h *ChatHandler) SendMessage(c *gin.Context) {
	var req types.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("invalid chat request payload")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			ErrorCode: "INVALID_PAYLOAD",
			Message:   "Invalid request payload",
		})
		return
	}

	response, err := h.controller.SendMessage(c.Request.Context(), req)
	if err != nil {
		// Check if it's our custom error type
		if errorResp, ok := err.(*types.ErrorResponse); ok {
			h.logErrorResponse(errorResp, c)
			
			switch errorResp.ErrorCode {
			case "SESSION_NOT_FOUND":
				c.JSON(http.StatusNotFound, errorResp)
			case "INVALID_REQUEST":
				c.JSON(http.StatusBadRequest, errorResp)
			default:
				c.JSON(http.StatusInternalServerError, errorResp)
			}
			return
		}

		// Generic error
		h.logger.WithError(err).Error("internal error processing chat message")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			ErrorCode: "INTERNAL_ERROR",
			Message:   "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateSession handles POST /api/sessions
func (h *ChatHandler) CreateSession(c *gin.Context) {
	var req types.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("invalid create session request payload")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			ErrorCode: "INVALID_PAYLOAD",
			Message:   "Invalid request payload",
		})
		return
	}

	session, err := h.controller.CreateSession(req)
	if err != nil {
		if errorResp, ok := err.(*types.ErrorResponse); ok {
			h.logErrorResponse(errorResp, c)
			c.JSON(http.StatusInternalServerError, errorResp)
			return
		}

		h.logger.WithError(err).Error("internal error creating session")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			ErrorCode: "INTERNAL_ERROR",
			Message:   "Internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, session)
}

// GetSession handles GET /api/sessions/:id
func (h *ChatHandler) GetSession(c *gin.Context) {
	sessionIDStr := c.Param("id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		h.logger.WithField("session_id", sessionIDStr).Error("invalid session ID format")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			ErrorCode: "INVALID_SESSION_ID",
			Message:   "Invalid session ID format",
		})
		return
	}

	session, err := h.controller.GetSession(sessionID)
	if err != nil {
		if errorResp, ok := err.(*types.ErrorResponse); ok {
			h.logErrorResponse(errorResp, c)
			
			if errorResp.ErrorCode == "SESSION_NOT_FOUND" {
				c.JSON(http.StatusNotFound, errorResp)
			} else {
				c.JSON(http.StatusInternalServerError, errorResp)
			}
			return
		}

		h.logger.WithError(err).Error("internal error getting session")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			ErrorCode: "INTERNAL_ERROR",
			Message:   "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, session)
}

// ListSessions handles GET /api/sessions
func (h *ChatHandler) ListSessions(c *gin.Context) {
	sessions, err := h.controller.ListSessions()
	if err != nil {
		if errorResp, ok := err.(*types.ErrorResponse); ok {
			h.logErrorResponse(errorResp, c)
			c.JSON(http.StatusInternalServerError, errorResp)
			return
		}

		h.logger.WithError(err).Error("internal error listing sessions")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			ErrorCode: "INTERNAL_ERROR",
			Message:   "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"sessions": sessions,
	})
}

// HealthCheck handles GET /health
func (h *ChatHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"service": "podscription-api",
		"version": "1.0.0",
	})
}

// logErrorResponse logs error responses with context
func (h *ChatHandler) logErrorResponse(errorResp *types.ErrorResponse, c *gin.Context) {
	h.logger.WithFields(logrus.Fields{
		"error_code":    errorResp.ErrorCode,
		"error_message": errorResp.Message,
		"method":        c.Request.Method,
		"path":          c.Request.URL.Path,
		"remote_addr":   c.ClientIP(),
	}).Error("returning error response")
}