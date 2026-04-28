-- Add agent_gateways table for tracking running agentgateway reverse proxy instances.

CREATE TABLE IF NOT EXISTS agent_gateways (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'unknown',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT check_agent_gateway_status_valid CHECK (status IN ('healthy', 'unhealthy', 'unknown'))
);

CREATE INDEX IF NOT EXISTS idx_agent_gateways_status ON agent_gateways (status);
