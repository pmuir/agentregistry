package agentgateways

import (
	"context"
	"errors"
	"net/http"

	agentgatewaysvc "github.com/agentregistry-dev/agentregistry/internal/registry/service/agentgateway"
	"github.com/agentregistry-dev/agentregistry/pkg/models"
	"github.com/agentregistry-dev/agentregistry/pkg/registry/database"
	"github.com/danielgtaylor/huma/v2"
)

type AgentGatewayListInput struct{}

type AgentGatewayByIDInput struct {
	GatewayID string `path:"gatewayId" json:"gatewayId" doc:"Agent gateway ID"`
}

type CreateAgentGatewayRequest struct {
	Body models.CreateAgentGatewayInput
}

type AgentGatewaysListResponse struct {
	Body struct {
		AgentGateways []models.AgentGateway `json:"agentGateways"`
		Count         int                   `json:"count"`
	}
}

type AgentGatewayResponse struct {
	Body models.AgentGateway
}

func agentGatewayReadHTTPError(action string, err error) error {
	switch {
	case errors.Is(err, database.ErrInvalidInput):
		return huma.Error400BadRequest(err.Error())
	case errors.Is(err, database.ErrNotFound):
		return huma.Error404NotFound("Agent gateway not found")
	default:
		return huma.Error500InternalServerError(action, err)
	}
}

func agentGatewayWriteHTTPError(action string, err error) error {
	switch {
	case errors.Is(err, database.ErrInvalidInput):
		return huma.Error400BadRequest(err.Error())
	case errors.Is(err, database.ErrAlreadyExists):
		return huma.Error409Conflict("Agent gateway already exists")
	case errors.Is(err, database.ErrNotFound):
		return huma.Error404NotFound("Agent gateway not found")
	default:
		return huma.Error500InternalServerError(action, err)
	}
}

func RegisterAgentGatewaysEndpoints(api huma.API, basePath string, svc agentgatewaysvc.Registry) {
	huma.Register(api, huma.Operation{
		OperationID: "list-agent-gateways",
		Method:      http.MethodGet,
		Path:        basePath + "/agentgateways",
		Summary:     "List agent gateways",
		Description: "List registered agentgateway instances.",
		Tags:        []string{"agentgateways"},
	}, func(ctx context.Context, _ *AgentGatewayListInput) (*AgentGatewaysListResponse, error) {
		resp := &AgentGatewaysListResponse{}
		resp.Body.AgentGateways = []models.AgentGateway{}
		gateways, err := svc.ListAgentGateways(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to list agent gateways", err)
		}
		for _, gw := range gateways {
			if gw == nil {
				continue
			}
			resp.Body.AgentGateways = append(resp.Body.AgentGateways, *gw)
		}
		resp.Body.Count = len(resp.Body.AgentGateways)
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "create-agent-gateway",
		Method:      http.MethodPost,
		Path:        basePath + "/agentgateways",
		Summary:     "Create agent gateway",
		Description: "Register a new agentgateway instance.",
		Tags:        []string{"agentgateways"},
	}, func(ctx context.Context, input *CreateAgentGatewayRequest) (*AgentGatewayResponse, error) {
		gw, err := svc.RegisterAgentGateway(ctx, &input.Body)
		if err != nil {
			return nil, agentGatewayWriteHTTPError("Failed to create agent gateway", err)
		}
		return &AgentGatewayResponse{Body: *gw}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-agent-gateway",
		Method:      http.MethodGet,
		Path:        basePath + "/agentgateways/{gatewayId}",
		Summary:     "Get agent gateway",
		Description: "Get an agentgateway instance by ID.",
		Tags:        []string{"agentgateways"},
	}, func(ctx context.Context, input *AgentGatewayByIDInput) (*AgentGatewayResponse, error) {
		gw, err := svc.GetAgentGateway(ctx, input.GatewayID)
		if err != nil {
			return nil, agentGatewayReadHTTPError("Failed to get agent gateway", err)
		}
		return &AgentGatewayResponse{Body: *gw}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-agent-gateway",
		Method:      http.MethodDelete,
		Path:        basePath + "/agentgateways/{gatewayId}",
		Summary:     "Delete agent gateway",
		Description: "Delete an agentgateway instance by ID.",
		Tags:        []string{"agentgateways"},
	}, func(ctx context.Context, input *AgentGatewayByIDInput) (*struct{}, error) {
		err := svc.DeleteAgentGateway(ctx, input.GatewayID)
		if err != nil {
			return nil, agentGatewayWriteHTTPError("Failed to delete agent gateway", err)
		}
		return &struct{}{}, nil
	})
}
