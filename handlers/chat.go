// handlers/chat.go
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"ai-agent-app/services" // Import the services package

	"github.com/gorilla/mux"
)

// ChatResponse represents the structure of the chat response
type ChatResponse struct {
	Message string `json:"message"`
}

// Global chat history for web requests
var webChatHistory = services.NewChatHistory(10)

// ChatWithAgent handles chat requests with the agent
func ChatWithAgent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentID := vars["agentID"]

	// Validate agentID
	if agentID == "" {
		http.Error(w, "agentID is required", http.StatusBadRequest)
		return
	}

	// Extract the message from the request body
	var requestBody struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	// Validate the message
	if requestBody.Message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	// Add user message to history
	webChatHistory.AddMessage(agentID, "user", requestBody.Message)

	// Get the conversation history
	history := webChatHistory.GetHistory(agentID)

	var responseMessage string
	var err error

	if agentID[:6] == "openai" {
		// Send message to OpenAI
		responseMessage, err = services.SendMessageToOpenAI(
			os.Getenv("OPENAI_API_KEY"),
			requestBody.Message,
			history, // Pass the message history
		)

		if err != nil {
			http.Error(w, "Error communicating with the agent", http.StatusInternalServerError)
			log.Printf("Error communicating with agent %s: %v", agentID, err)
			return
		}

		// Add the assistant's response to history
		webChatHistory.AddMessage(agentID, "assistant", responseMessage)
	} else if agentID[:4] == "grok" {
		responseMessage, err = services.SendMessageToGrok(requestBody.Message)
	} else {
		http.Error(w, "Unknown agent type", http.StatusBadRequest)
		return
	}

	// Handle any errors from the service call
	if err != nil {
		http.Error(w, "Error communicating with the agent", http.StatusInternalServerError)
		log.Printf("Error communicating with agent %s: %v", agentID, err)
		return
	}

	// Log the chat request
	log.Printf("Chat request for agentID: %s, message: %s", agentID, requestBody.Message)

	// Send response
	response := ChatResponse{
		Message: responseMessage,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ConsoleChatWithAgent handles chat interactions from the console
func ConsoleChatWithAgent(agentID string, message string, chatHistory *services.ChatHistory) (string, error) {
	var responseMessage string
	var err error

	// Get the conversation history
	history := chatHistory.GetHistory(agentID)
	// Use the OpenAI API to generate a response
	responseMessage, err = services.SendMessageToOpenAI(
		os.Getenv("OPENAI_API_KEY"),
		message,
		history,
	)

	if err != nil {
		return "", fmt.Errorf("error communicating with agent %s: %v", agentID, err)
	}

	// Log the console chat request
	log.Printf("Console chat request for agentID: %s, message: %s", agentID, message)

	return responseMessage, nil
}
