package types

import (
	"time"

	"github.com/google/uuid"
)

// MessageRole represents the role of a message sender
type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
)

// IntentCategory represents the type of Kubernetes issue
type IntentCategory string

const (
	IntentCategoryNetworking  IntentCategory = "networking"
	IntentCategoryStorage     IntentCategory = "storage"
	IntentCategoryPodIssues   IntentCategory = "pod-issues"
	IntentCategoryRBAC        IntentCategory = "rbac"
	IntentCategoryPerformance IntentCategory = "performance"
	IntentCategoryGeneral     IntentCategory = "general"
)

// PodIntent represents the classified intent of a user's message
type PodIntent struct {
	Category   IntentCategory `json:"category"`
	Confidence float64        `json:"confidence"`
	Symptoms   []string       `json:"symptoms"`
}

// Prescription represents the AI's structured response
type Prescription struct {
	Diagnosis string   `json:"diagnosis"`
	Treatment string   `json:"treatment"`
	Commands  []string `json:"commands,omitempty"`
	FollowUp  string   `json:"followUp,omitempty"`
}

// Message represents a single message in a conversation
type Message struct {
	ID           uuid.UUID     `json:"id"`
	Role         MessageRole   `json:"role"`
	Content      string        `json:"content"`
	Timestamp    time.Time     `json:"timestamp"`
	Intent       *PodIntent    `json:"intent,omitempty"`
	Prescription *Prescription `json:"prescription,omitempty"`
}

// Session represents a conversation session
type Session struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Messages  []Message `json:"messages"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ChatRequest represents an incoming chat message request
type ChatRequest struct {
	SessionID *uuid.UUID `json:"sessionId,omitempty"`
	Content   string     `json:"content"`
}

// ChatResponse represents the response to a chat request
type ChatResponse struct {
	Session Session `json:"session"`
	Message Message `json:"message"`
}

// CreateSessionRequest represents a request to create a new session
type CreateSessionRequest struct {
	Name string `json:"name,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	ErrorCode string `json:"error"`
	Code      string `json:"code,omitempty"`
	Message   string `json:"message"`
}

// Error implements the error interface
func (e *ErrorResponse) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.ErrorCode
}