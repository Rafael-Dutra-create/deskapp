-- Migration: init
-- Created at: 2025-11-14T11:08:02-03:00
-- Created by: victorgomes


CREATE TABLE IF NOT EXISTS migration_logs (
    id SERIAL PRIMARY KEY,
    from_version BIGINT,
    to_version BIGINT,
    migration_name VARCHAR(255) NOT NULL,
    applied_by VARCHAR(100) NOT NULL,
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    execution_time INTERVAL,
    success BOOLEAN NOT NULL DEFAULT false,
    error_message TEXT,
    environment VARCHAR(50) DEFAULT 'development'
);

CREATE INDEX IF NOT EXISTS idx_migration_logs_timestamps ON migration_logs(started_at, completed_at);
CREATE INDEX IF NOT EXISTS idx_migration_logs_versions ON migration_logs(from_version, to_version);