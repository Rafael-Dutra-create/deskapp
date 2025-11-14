package main

import (
	"database/sql"
	"time"

	"github.com/golang-migrate/migrate/v4"
)


// logMigrationStart registra o início de uma migração
func (mm *MigrationManager) logMigrationStart(tx *sql.Tx, fromVersion, toVersion uint, name, appliedBy string) error {
	if err := mm.ensureMigrationLogsTable(); err != nil {
		return err
	}
	query := `
        INSERT INTO migration_logs 
        (from_version, to_version, migration_name, applied_by, started_at, success) 
        VALUES ($1, $2, $3, $4, $5, false)
        RETURNING id
    `
	var logID int
	return tx.QueryRow(query, fromVersion, toVersion, name, appliedBy, time.Now()).Scan(&logID)
}

// logMigrationResult atualiza o log com o resultado
func (mm *MigrationManager) logMigrationResult(tx *sql.Tx, fromVersion, toVersion uint, name, appliedBy string, executionTime time.Duration, migrationErr error) error {
	success := migrationErr == nil || migrationErr == migrate.ErrNoChange
	errorMsg := ""
	if migrationErr != nil && migrationErr != migrate.ErrNoChange {
		errorMsg = migrationErr.Error()
	}

	query := `
        UPDATE migration_logs 
        SET completed_at = $1, execution_time = $2, success = $3, error_message = $4
        WHERE from_version = $5 AND to_version = $6 AND migration_name = $7
    `
	_, err := tx.Exec(query, time.Now(), executionTime, success, errorMsg, fromVersion, toVersion, name)
	return err
}
