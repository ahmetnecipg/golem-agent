// services/agent.go
package services

import (
	"ai-agent-app/database"
	"ai-agent-app/models"
	"log"
)

// CreateAgent saves a new agent to the database
func CreateAgent(agent models.Agent) error {
	// Prepare the SQL statement
	query := `INSERT INTO agents (id, name, type, context) VALUES ($1, $2, $3, $4)`
	_, err := database.Exec(query, agent.ID, agent.Name, agent.Type, agent.Context)
	if err != nil {
		log.Printf("Error saving agent to database: %v", err)
		return err
	}
	return nil
}
