package agentgateway

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/agentregistry-dev/agentregistry/pkg/models"
	"github.com/agentregistry-dev/agentregistry/pkg/registry/database"
)

type Dependencies struct {
	StoreDB       database.Store
	AgentGateways database.AgentGatewayStore
}

type Registry interface {
	ListAgentGateways(ctx context.Context) ([]*models.AgentGateway, error)
	RegisterAgentGateway(ctx context.Context, in *models.CreateAgentGatewayInput) (*models.AgentGateway, error)
	GetAgentGateway(ctx context.Context, gatewayID string) (*models.AgentGateway, error)
	UpdateAgentGateway(ctx context.Context, gatewayID string, in *models.UpdateAgentGatewayInput) (*models.AgentGateway, error)
	ApplyAgentGateway(ctx context.Context, gatewayID string, in *models.UpdateAgentGatewayInput) (*models.AgentGateway, error)
	DeleteAgentGateway(ctx context.Context, gatewayID string) error
}

type registry struct {
	agentGateways database.AgentGatewayStore
}

var _ Registry = (*registry)(nil)

func New(deps Dependencies) Registry {
	if deps.AgentGateways == nil && deps.StoreDB != nil {
		deps.AgentGateways = deps.StoreDB.AgentGateways()
	}
	return &registry{agentGateways: deps.AgentGateways}
}

func (r *registry) ListAgentGateways(ctx context.Context) ([]*models.AgentGateway, error) {
	return r.agentGateways.ListAgentGateways(ctx)
}

func (r *registry) RegisterAgentGateway(ctx context.Context, in *models.CreateAgentGatewayInput) (*models.AgentGateway, error) {
	if in == nil {
		return nil, database.ErrInvalidInput
	}
	if strings.TrimSpace(in.Address) == "" {
		return nil, fmt.Errorf("%w: agent gateway address is required", database.ErrInvalidInput)
	}
	return r.agentGateways.CreateAgentGateway(ctx, in)
}

func (r *registry) GetAgentGateway(ctx context.Context, gatewayID string) (*models.AgentGateway, error) {
	if strings.TrimSpace(gatewayID) == "" {
		return nil, fmt.Errorf("%w: agent gateway id is required", database.ErrInvalidInput)
	}
	return r.agentGateways.GetAgentGateway(ctx, gatewayID)
}

func (r *registry) UpdateAgentGateway(ctx context.Context, gatewayID string, in *models.UpdateAgentGatewayInput) (*models.AgentGateway, error) {
	return r.agentGateways.UpdateAgentGateway(ctx, gatewayID, in)
}

func (r *registry) ApplyAgentGateway(ctx context.Context, gatewayID string, in *models.UpdateAgentGatewayInput) (*models.AgentGateway, error) {
	resolvedID := strings.TrimSpace(gatewayID)
	if resolvedID == "" {
		return nil, fmt.Errorf("%w: agent gateway id is required", database.ErrInvalidInput)
	}
	if in == nil {
		return nil, database.ErrInvalidInput
	}

	existing, err := r.agentGateways.GetAgentGateway(ctx, resolvedID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return nil, err
	}
	if existing != nil {
		return r.agentGateways.UpdateAgentGateway(ctx, existing.ID, in)
	}

	if in.Name == nil || strings.TrimSpace(*in.Name) == "" {
		return nil, fmt.Errorf("%w: agent gateway name is required when creating", database.ErrInvalidInput)
	}
	address := ""
	if in.Address != nil {
		address = *in.Address
	}
	if strings.TrimSpace(address) == "" {
		return nil, fmt.Errorf("%w: agent gateway address is required when creating", database.ErrInvalidInput)
	}

	return r.agentGateways.CreateAgentGateway(ctx, &models.CreateAgentGatewayInput{
		ID:      resolvedID,
		Name:    *in.Name,
		Address: address,
	})
}

func (r *registry) DeleteAgentGateway(ctx context.Context, gatewayID string) error {
	if strings.TrimSpace(gatewayID) == "" {
		return fmt.Errorf("%w: agent gateway id is required", database.ErrInvalidInput)
	}
	return r.agentGateways.DeleteAgentGateway(ctx, gatewayID)
}
