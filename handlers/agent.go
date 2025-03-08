package handlers

import (
	"ai-agent-app/models"
	"ai-agent-app/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// CreateAgentResponse represents the structure of the create agent response
type CreateAgentResponse struct {
	Agent models.Agent `json:"agent"`
}

// CreateAgent handles the creation of a new agent
func CreateAgent(w http.ResponseWriter, r *http.Request) {
	var agent models.Agent
	if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	// Basic validation (you can expand this as needed)
	if agent.Name == "" {
		http.Error(w, "Agent name is required", http.StatusBadRequest)
		return
	}

	// Generate a unique ID for the agent (this can be improved)
	agent.ID = generateUniqueID() // Implement this function to generate a unique ID

	// Call the service to save the agent to the database
	if err := services.CreateAgent(agent); err != nil {
		http.Error(w, "Error saving agent to database", http.StatusInternalServerError)
		log.Printf("Error saving agent: %v", err)
		return
	}

	// Log the creation of the agent
	log.Printf("Agent created: %+v", agent)

	// Send response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateAgentResponse{Agent: agent})
}

// CreateDefaultAgent creates a default agent and returns its ID
func CreateDefaultAgent() (string, error) {
	// You can customize the agent properties here
	agent := models.Agent{
		Name: "Console Agent",
		// Remove the Description field since it doesn't exist in the models.Agent struct
		// Add other properties as needed based on the actual fields in models.Agent
	}

	// Generate a unique ID for the agent
	agent.ID = generateUniqueID()

	// Store the agent (this depends on your existing implementation)
	// For example:
	// database.SaveAgent(agent)

	return agent.ID, nil
}

// generateUniqueID creates a unique identifier for the agent
func generateUniqueID() string {
	// Simple implementation - in a real app, use a proper UUID library
	return fmt.Sprintf("agent-%d", time.Now().UnixNano())
}
