package healthcheck

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/agentregistry-dev/agentregistry/pkg/models"
	"github.com/agentregistry-dev/agentregistry/pkg/registry/auth"
	"github.com/agentregistry-dev/agentregistry/pkg/registry/database"
)

const maxConcurrentProbes = 10

type GatewayHealthChecker struct {
	store    database.AgentGatewayStore
	interval time.Duration
	client   *http.Client
}

func NewGatewayHealthChecker(store database.AgentGatewayStore, interval time.Duration, client *http.Client) *GatewayHealthChecker {
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	return &GatewayHealthChecker{
		store:    store,
		interval: interval,
		client:   client,
	}
}

func (c *GatewayHealthChecker) Run(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.checkAll(ctx)
		}
	}
}

func (c *GatewayHealthChecker) checkAll(ctx context.Context) {
	systemCtx := auth.WithSystemContext(ctx)

	gateways, err := c.store.ListAgentGateways(systemCtx)
	if err != nil {
		slog.Debug("health checker: failed to list agent gateways", "error", err)
		return
	}

	sem := make(chan struct{}, maxConcurrentProbes)
	var wg sync.WaitGroup

	for _, gw := range gateways {
		if gw == nil {
			continue
		}
		wg.Add(1)
		sem <- struct{}{}
		go func(gw *models.AgentGateway) {
			defer wg.Done()
			defer func() { <-sem }()

			status := c.probe(ctx, gw.Address)
			if status != gw.Status {
				if err := c.store.UpdateAgentGatewayStatus(systemCtx, gw.ID, status); err != nil {
					slog.Debug("health checker: failed to update status", "gateway", gw.ID, "error", err)
				}
			}
		}(gw)
	}

	wg.Wait()
}

func (c *GatewayHealthChecker) probe(ctx context.Context, address string) string {
	url := strings.TrimRight(address, "/") + "/healthz"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return models.AgentGatewayStatusUnhealthy
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return models.AgentGatewayStatusUnhealthy
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return models.AgentGatewayStatusHealthy
	}
	return models.AgentGatewayStatusUnhealthy
}
