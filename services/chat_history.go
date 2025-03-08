package services

import (
	"sync"
)

// Message represents a single message in the chat history
type Message struct {
	Role    string `json:"role"`    // "user" or "assistant"
	Content string `json:"content"` // The message content
}

// ChatHistory stores conversation history for each agent
type ChatHistory struct {
	conversations map[string][]Message
	mutex         sync.RWMutex
	maxLength     int
}

// NewChatHistory creates a new chat history manager
func NewChatHistory(maxLength int) *ChatHistory {
	return &ChatHistory{
		conversations: make(map[string][]Message),
		maxLength:     maxLength,
	}
}

// AddMessage adds a message to the conversation history for a specific agent
func (ch *ChatHistory) AddMessage(agentID, role, content string) {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()

	// Create conversation if it doesn't exist
	if _, exists := ch.conversations[agentID]; !exists {
		ch.conversations[agentID] = make([]Message, 0)
	}

	// Add the message
	ch.conversations[agentID] = append(ch.conversations[agentID], Message{
		Role:    role,
		Content: content,
	})

	// Trim history if it exceeds maximum length
	if len(ch.conversations[agentID]) > ch.maxLength {
		ch.conversations[agentID] = ch.conversations[agentID][len(ch.conversations[agentID])-ch.maxLength:]
	}
}

// GetHistory returns the conversation history for a specific agent
func (ch *ChatHistory) GetHistory(agentID string) []Message {
	ch.mutex.RLock()
	defer ch.mutex.RUnlock()

	if history, exists := ch.conversations[agentID]; exists {
		// Return a copy to prevent concurrent modification issues
		result := make([]Message, len(history))
		copy(result, history)
		return result
	}

	return []Message{}
}

// ClearHistory clears the conversation history for a specific agent
func (ch *ChatHistory) ClearHistory(agentID string) {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()

	delete(ch.conversations, agentID)
}
