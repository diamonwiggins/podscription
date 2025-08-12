package managers

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
	"podscription-api/pkg/config"
	"podscription-api/types"
)

// OpenAIManager handles OpenAI API interactions
type OpenAIManager struct {
	client            *openai.Client
	config            config.OpenAI
	specializedPrompts *SpecializedPrompts
}

// NewOpenAIManager creates a new OpenAI manager
func NewOpenAIManager(cfg config.OpenAI) *OpenAIManager {
	client := openai.NewClient(cfg.APIKey)
	
	return &OpenAIManager{
		client:             client,
		config:             cfg,
		specializedPrompts: &SpecializedPrompts{},
	}
}

// ClassifyIntent analyzes a user message to determine the Kubernetes troubleshooting category
func (m *OpenAIManager) ClassifyIntent(ctx context.Context, message string) (*types.PodIntent, error) {
	prompt := m.buildIntentClassificationPrompt(message)
	
	resp, err := m.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       m.config.Model,
		Temperature: 0.3, // Lower temperature for more consistent classification
		MaxTokens:   200,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: prompt.System,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt.User,
			},
		},
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to classify intent: %w", err)
	}
	
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no classification response received")
	}
	
	return m.parseIntentResponse(resp.Choices[0].Message.Content)
}

// GenerateDiagnosis creates a medical-themed Kubernetes troubleshooting response
func (m *OpenAIManager) GenerateDiagnosis(ctx context.Context, message string, intent *types.PodIntent, history []types.Message) (*types.Prescription, string, error) {
	prompt := m.buildDiagnosisPrompt(message, intent, history)
	
	resp, err := m.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       m.config.Model,
		Temperature: m.config.Temperature,
		MaxTokens:   m.config.MaxTokens,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: prompt.System,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt.User,
			},
		},
	})
	
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate diagnosis: %w", err)
	}
	
	if len(resp.Choices) == 0 {
		return nil, "", fmt.Errorf("no diagnosis response received")
	}
	
	response := resp.Choices[0].Message.Content
	prescription := m.parseDiagnosisResponse(response, intent)
	
	return prescription, response, nil
}

type promptPair struct {
	System string
	User   string
}

// buildIntentClassificationPrompt creates prompts for intent classification
func (m *OpenAIManager) buildIntentClassificationPrompt(message string) promptPair {
	system := `You are an expert Kubernetes troubleshooting assistant. Your job is to classify user messages into specific Kubernetes problem categories.

Analyze the user's message and classify it into one of these categories:
- networking: Service discovery, ingress, connectivity, DNS issues
- storage: PVC, PV, volume mounts, disk space, storage classes
- pod-issues: Pod startup, container crashes, image pulls, resource constraints
- rbac: Permissions, service accounts, cluster roles, security
- performance: CPU, memory, scaling, resource optimization
- general: General questions, cluster info, basic troubleshooting

Respond with ONLY this format:
CATEGORY: [category name]
CONFIDENCE: [0.0-1.0]
SYMPTOMS: [comma-separated list of 2-3 key symptoms detected]

Be concise and accurate.`

	user := fmt.Sprintf("Classify this Kubernetes issue: %s", message)
	
	return promptPair{System: system, User: user}
}

// buildDiagnosisPrompt creates prompts for medical-themed diagnosis
func (m *OpenAIManager) buildDiagnosisPrompt(message string, intent *types.PodIntent, history []types.Message) promptPair {
	// Use specialized prompts for networking and storage
	switch intent.Category {
	case types.IntentCategoryNetworking:
		return m.specializedPrompts.GetNetworkingPrompt(message, history)
	case types.IntentCategoryStorage:
		return m.specializedPrompts.GetStoragePrompt(message, history)
	default:
		// Fall back to generic Pod Doctor prompt for other categories
		return m.buildGenericDiagnosisPrompt(message, intent, history)
	}
}

// buildGenericDiagnosisPrompt creates generic prompts for non-specialized categories
func (m *OpenAIManager) buildGenericDiagnosisPrompt(message string, intent *types.PodIntent, history []types.Message) promptPair {
	categoryContext := m.getCategoryContext(intent.Category)
	
	system := fmt.Sprintf(`You are the "Pod Doctor" - a Kubernetes troubleshooting assistant with a medical personality. You diagnose and treat "sick" Kubernetes pods and clusters.

Your specialty: %s

PERSONALITY:
- Speak like a doctor treating patients
- Use medical metaphors and terminology
- Be professional but friendly
- Provide clear "prescriptions" (solutions)
- Reference "symptoms" (error conditions) and "treatments" (fixes)

RESPONSE FORMAT - Use exactly this structure:
## Diagnosis: [Medical-style diagnosis name]

[Brief explanation of the issue using medical metaphors]

### Prescribed Treatment:
1. **[Step name]**: ` + "`" + `[command or action]` + "`" + `
2. **[Step name]**: ` + "`" + `[command or action]` + "`" + `
[Continue with numbered steps]

### Follow-up Care:
[Additional guidance or next steps]

*[End with a medical-themed joke or memorable phrase]*

CONTEXT: %s`, intent.Category, categoryContext)

	// Include conversation history for context
	historyContext := ""
	if len(history) > 0 {
		historyContext = "\n\nPrevious consultation history:\n"
		for _, msg := range history {
			if len(historyContext) > 500 { // Limit context length
				break
			}
			historyContext += fmt.Sprintf("%s: %s\n", msg.Role, msg.Content[:min(100, len(msg.Content))])
		}
	}

	user := fmt.Sprintf("Patient symptoms: %s%s", message, historyContext)
	
	return promptPair{System: system, User: user}
}

// getCategoryContext provides specialized context for each category
func (m *OpenAIManager) getCategoryContext(category types.IntentCategory) string {
	contexts := map[types.IntentCategory]string{
		types.IntentCategoryNetworking: "Focus on service discovery, ingress configuration, DNS resolution, network policies, and connectivity issues. Common treatments include checking service selectors, endpoints, and network policies.",
		types.IntentCategoryStorage: "Focus on persistent volumes, volume claims, storage classes, and mount issues. Common treatments include checking PVC status, storage class availability, and mount permissions.",
		types.IntentCategoryPodIssues: "Focus on pod lifecycle, container startup, image pulls, and resource constraints. Common treatments include checking pod events, logs, and resource limits.",
		types.IntentCategoryRBAC: "Focus on permissions, service accounts, roles, and security policies. Common treatments include checking RBAC rules, service account permissions, and security contexts.",
		types.IntentCategoryPerformance: "Focus on resource utilization, scaling, and optimization. Common treatments include adjusting resource requests/limits, HPA configuration, and performance tuning.",
		types.IntentCategoryGeneral: "Provide general Kubernetes guidance and best practices. Focus on cluster health, basic troubleshooting, and educational responses.",
	}
	
	return contexts[category]
}

// parseIntentResponse extracts intent information from the classification response
func (m *OpenAIManager) parseIntentResponse(response string) (*types.PodIntent, error) {
	lines := strings.Split(strings.TrimSpace(response), "\n")
	
	intent := &types.PodIntent{
		Category:   types.IntentCategoryGeneral,
		Confidence: 0.7,
		Symptoms:   []string{},
	}
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		if strings.HasPrefix(line, "CATEGORY:") {
			categoryStr := strings.TrimSpace(strings.TrimPrefix(line, "CATEGORY:"))
			intent.Category = types.IntentCategory(categoryStr)
		} else if strings.HasPrefix(line, "CONFIDENCE:") {
			// Parse confidence (basic parsing, could be improved)
			confStr := strings.TrimSpace(strings.TrimPrefix(line, "CONFIDENCE:"))
			if strings.Contains(confStr, "0.") || strings.Contains(confStr, "1.") {
				// Simple confidence extraction
				if strings.Contains(confStr, "0.9") || strings.Contains(confStr, "0.8") {
					intent.Confidence = 0.9
				} else if strings.Contains(confStr, "0.7") || strings.Contains(confStr, "0.6") {
					intent.Confidence = 0.8
				} else {
					intent.Confidence = 0.7
				}
			}
		} else if strings.HasPrefix(line, "SYMPTOMS:") {
			symptomsStr := strings.TrimSpace(strings.TrimPrefix(line, "SYMPTOMS:"))
			if symptomsStr != "" {
				symptoms := strings.Split(symptomsStr, ",")
				for _, symptom := range symptoms {
					intent.Symptoms = append(intent.Symptoms, strings.TrimSpace(symptom))
				}
			}
		}
	}
	
	return intent, nil
}

// parseDiagnosisResponse extracts prescription information from the diagnosis response
func (m *OpenAIManager) parseDiagnosisResponse(response string, intent *types.PodIntent) *types.Prescription {
	// Extract diagnosis from the response (simple parsing)
	diagnosis := "Kubernetes Issue Diagnosis"
	if strings.Contains(response, "Diagnosis:") {
		parts := strings.Split(response, "Diagnosis:")
		if len(parts) > 1 {
			diagnosisPart := strings.Split(parts[1], "\n")[0]
			diagnosis = strings.TrimSpace(strings.TrimPrefix(diagnosisPart, "🩺"))
		}
	}
	
	// Extract commands (look for backticked commands)
	commands := []string{}
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		if strings.Contains(line, "`") {
			// Extract commands from backticks
			start := strings.Index(line, "`")
			end := strings.LastIndex(line, "`")
			if start != -1 && end != -1 && start != end {
				command := strings.TrimSpace(line[start+1 : end])
				if command != "" && strings.HasPrefix(command, "kubectl") {
					commands = append(commands, command)
				}
			}
		}
	}
	
	// Extract follow-up (look for Follow-up section)
	followUp := ""
	if strings.Contains(response, "Follow-up Care:") {
		parts := strings.Split(response, "Follow-up Care:")
		if len(parts) > 1 {
			followUpPart := strings.Split(parts[1], "*")[0] // Stop at the joke
			followUp = strings.TrimSpace(followUpPart)
		}
	}
	
	return &types.Prescription{
		Diagnosis: diagnosis,
		Treatment: "Refer to the detailed diagnosis above for treatment recommendations.",
		Commands:  commands,
		FollowUp:  followUp,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}