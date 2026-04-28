package models

import "time"

const (
	AgentGatewayStatusHealthy   = "healthy"
	AgentGatewayStatusUnhealthy = "unhealthy"
	AgentGatewayStatusUnknown   = "unknown"
)

// AgentGateway represents a running agentgateway reverse proxy instance.
type AgentGateway struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateAgentGatewayInput defines inputs for agentgateway creation.
type CreateAgentGatewayInput struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

// UpdateAgentGatewayInput defines inputs for agentgateway updates.
type UpdateAgentGatewayInput struct {
	Name    *string `json:"name,omitempty"`
	Address *string `json:"address,omitempty"`
}
