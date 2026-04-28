package database

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/agentregistry-dev/agentregistry/pkg/models"
	"github.com/agentregistry-dev/agentregistry/pkg/registry/auth"
	"github.com/agentregistry-dev/agentregistry/pkg/registry/database"
)

type agentGatewayStore struct {
	repositoryBase
}

var _ database.AgentGatewayStore = (*agentGatewayStore)(nil)

func (s *agentGatewayStore) CreateAgentGateway(ctx context.Context, in *models.CreateAgentGatewayInput) (*models.AgentGateway, error) {
	if in == nil {
		return nil, database.ErrInvalidInput
	}
	if strings.TrimSpace(in.ID) == "" || strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Address) == "" {
		return nil, database.ErrInvalidInput
	}
	if err := s.authz.Check(ctx, auth.PermissionActionPublish, auth.Resource{
		Name: in.ID,
		Type: auth.PermissionArtifactTypeAgentGateway,
	}); err != nil {
		return nil, err
	}
	query := `
		INSERT INTO agent_gateways (id, name, address)
		VALUES ($1, $2, $3)
		RETURNING id, name, address, status, created_at, updated_at
	`
	var gw models.AgentGateway
	err := s.executor.QueryRow(ctx, query, in.ID, in.Name, in.Address).Scan(
		&gw.ID, &gw.Name, &gw.Address, &gw.Status, &gw.CreatedAt, &gw.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, database.ErrAlreadyExists
		}
		return nil, fmt.Errorf("failed to create agent gateway: %w", err)
	}
	return &gw, nil
}

func (s *agentGatewayStore) ListAgentGateways(ctx context.Context) ([]*models.AgentGateway, error) {
	query := `SELECT id, name, address, status, created_at, updated_at FROM agent_gateways ORDER BY created_at ASC`
	rows, err := s.executor.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list agent gateways: %w", err)
	}
	defer rows.Close()
	var out []*models.AgentGateway
	for rows.Next() {
		var gw models.AgentGateway
		if err := rows.Scan(&gw.ID, &gw.Name, &gw.Address, &gw.Status, &gw.CreatedAt, &gw.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan agent gateway: %w", err)
		}
		out = append(out, &gw)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate agent gateways: %w", err)
	}
	return out, nil
}

func (s *agentGatewayStore) GetAgentGateway(ctx context.Context, gatewayID string) (*models.AgentGateway, error) {
	if err := s.authz.Check(ctx, auth.PermissionActionRead, auth.Resource{
		Name: gatewayID,
		Type: auth.PermissionArtifactTypeAgentGateway,
	}); err != nil {
		return nil, err
	}
	query := `SELECT id, name, address, status, created_at, updated_at FROM agent_gateways WHERE id = $1`
	var gw models.AgentGateway
	if err := s.executor.QueryRow(ctx, query, gatewayID).Scan(
		&gw.ID, &gw.Name, &gw.Address, &gw.Status, &gw.CreatedAt, &gw.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get agent gateway: %w", err)
	}
	return &gw, nil
}

func (s *agentGatewayStore) UpdateAgentGateway(ctx context.Context, gatewayID string, in *models.UpdateAgentGatewayInput) (*models.AgentGateway, error) {
	if err := s.authz.Check(ctx, auth.PermissionActionEdit, auth.Resource{
		Name: gatewayID,
		Type: auth.PermissionArtifactTypeAgentGateway,
	}); err != nil {
		return nil, err
	}
	if in == nil {
		return s.GetAgentGateway(ctx, gatewayID)
	}
	current, err := s.GetAgentGateway(ctx, gatewayID)
	if err != nil {
		return nil, err
	}
	name := current.Name
	if in.Name != nil {
		name = *in.Name
	}
	address := current.Address
	if in.Address != nil {
		address = *in.Address
	}
	query := `
		UPDATE agent_gateways
		SET name = $2, address = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, address, status, created_at, updated_at
	`
	var gw models.AgentGateway
	if err := s.executor.QueryRow(ctx, query, gatewayID, name, address).Scan(
		&gw.ID, &gw.Name, &gw.Address, &gw.Status, &gw.CreatedAt, &gw.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrNotFound
		}
		return nil, fmt.Errorf("failed to update agent gateway: %w", err)
	}
	return &gw, nil
}

func (s *agentGatewayStore) UpdateAgentGatewayStatus(ctx context.Context, gatewayID string, status string) error {
	result, err := s.executor.Exec(ctx,
		`UPDATE agent_gateways SET status = $2, updated_at = NOW() WHERE id = $1`,
		gatewayID, status,
	)
	if err != nil {
		return fmt.Errorf("failed to update agent gateway status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return database.ErrNotFound
	}
	return nil
}

func (s *agentGatewayStore) DeleteAgentGateway(ctx context.Context, gatewayID string) error {
	if err := s.authz.Check(ctx, auth.PermissionActionDelete, auth.Resource{
		Name: gatewayID,
		Type: auth.PermissionArtifactTypeAgentGateway,
	}); err != nil {
		return err
	}
	result, err := s.executor.Exec(ctx, `DELETE FROM agent_gateways WHERE id = $1`, gatewayID)
	if err != nil {
		return fmt.Errorf("failed to delete agent gateway: %w", err)
	}
	if result.RowsAffected() == 0 {
		return database.ErrNotFound
	}
	return nil
}
