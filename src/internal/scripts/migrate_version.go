package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
)

func (mm *MigrationManager) ensureMigrationLogsTable() error {
	checkQuery := `SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'migration_logs')`
	var exists bool
	if err := mm.db.QueryRow(checkQuery).Scan(&exists); err != nil {
		return err
	}

	if !exists {
		fmt.Println("📋 Criando tabela migration_logs...")
		createQuery := `
            CREATE TABLE migration_logs (
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
            )
        `
		if _, err := mm.db.Exec(createQuery); err != nil {
			return fmt.Errorf("failed to create migration_logs table: %v", err)
		}
		fmt.Println("✅ Tabela migration_logs criada com sucesso")
	}
	return nil
}

// GetMigrationSequence retorna a sequência ordenada de migrações disponíveis
func (mm *MigrationManager) GetMigrationSequence() ([]uint, error) {
	files, err := os.ReadDir(mm.migrationPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %v", err)
	}

	var versions []uint

	for _, file := range files {
		filename := file.Name()

		// Verificar se é um arquivo de migração up
		if strings.HasSuffix(filename, ".up.sql") {
			// Extrair a versão do nome do arquivo (parte antes do primeiro _)
			parts := strings.Split(filename, "_")
			if len(parts) > 0 {
				// A versão é a primeira parte do nome do arquivo
				versionStr := parts[0]

				// Converter para uint
				version, err := strconv.ParseUint(versionStr, 10, 64)
				if err != nil {
					// Se não for número, pular (pode ser um arquivo inválido)
					continue
				}

				versions = append(versions, uint(version))
			}
		}
	}

	// Ordenar as versões em ordem crescente
	sort.Slice(versions, func(i, j int) bool {
		return versions[i] < versions[j]
	})

	return versions, nil
}

// GetPreviousVersion encontra a versão anterior na sequência
func (mm *MigrationManager) GetPreviousVersion(currentVersion uint) (uint, error) {
	versions, err := mm.GetMigrationSequence()
	if err != nil {
		return 0, err
	}

	if len(versions) == 0 {
		return 0, fmt.Errorf("nenhuma migração encontrada")
	}

	// Encontrar a posição da versão atual
	currentIndex := -1
	for i, v := range versions {
		if v == currentVersion {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		return 0, fmt.Errorf("versão atual %d não encontrada na sequência", currentVersion)
	}

	// Se for a primeira migração, retornar 0 (versão base)
	if currentIndex == 0 {
		return 0, nil
	}

	// Retornar a versão anterior
	return versions[currentIndex-1], nil
}

// executeMigrationsInTransaction executa múltiplas migrações em uma transação
func (mm *MigrationManager) executeMigrationsInTransaction(versions []uint, olderVersion uint, direction string) error {
	// Iniciar transação
	tx, err := mm.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Flag para controlar rollback
	var success bool
	defer func() {
		if !success {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("⚠️ Error during rollback: %v", rbErr)
			}
			fmt.Println("🔻 Rollback executado - todas as alterações foram revertidas")
		}
	}()

	fmt.Printf("🚀 Iniciando transação com %d migrações...\n", len(versions))

	for i, version := range versions {
		migrationName := mm.GetMigrationName(version)
		start := time.Now()
		if err := mm.logMigrationStart(tx, version, olderVersion, migrationName, os.Getenv("USER")); err != nil {
			return err
		}

		fmt.Printf("📦 Aplicando migração %d/%d: %s (v%d)...\n",
			i+1, len(versions), migrationName, version)

		// Ler e executar o arquivo SQL manualmente na transação
		if err := mm.executeMigrationFile(tx, version, direction); err != nil {
			return fmt.Errorf("failed to execute migration %d (%s): %v",
				version, migrationName, err)
		}

		// Atualizar schema_migrations dentro da transação
		if err := mm.updateSchemaVersionInTx(tx, version, false); err != nil {
			return fmt.Errorf("failed to update schema version for %d: %v", version, err)
		}
		if err := mm.logMigrationResult(tx, version, olderVersion, migrationName, os.Getenv("USER"), time.Since(start), err); err != nil {
			return err
		}

		fmt.Printf("✅ Migração %d aplicada com sucesso\n", version)
	}

	// Commit da transação
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	success = true
	fmt.Printf("✅ Transação commitada com sucesso - %d migrações aplicadas\n", len(versions))
	return nil
}

// GetNextVersion encontra a próxima versão na sequência
func (mm *MigrationManager) GetNextVersion(currentVersion uint) (uint, error) {
	versions, err := mm.GetMigrationSequence()
	if err != nil {
		return 0, err
	}

	if len(versions) == 0 {
		return 0, fmt.Errorf("nenhuma migração encontrada")
	}

	// Se não há versão atual, retornar a primeira
	if currentVersion == 0 {
		return versions[0], nil
	}

	// Encontrar a posição da versão atual
	currentIndex := -1
	for i, v := range versions {
		if v == currentVersion {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		return 0, fmt.Errorf("versão atual %d não encontrada na sequência", currentVersion)
	}

	// Se for a última migração, retornar erro
	if currentIndex == len(versions)-1 {
		return 0, fmt.Errorf("já está na última versão")
	}

	// Retornar a próxima versão
	return versions[currentIndex+1], nil
}

// executeMigrationFile lê e executa um arquivo de migração na transação
func (mm *MigrationManager) executeMigrationFile(tx *sql.Tx, version uint, direction string) error {
	// Encontrar o arquivo de migração
	filename := fmt.Sprintf("%d_", version)
	files, err := os.ReadDir(mm.migrationPath)
	if err != nil {
		return err
	}

	var migrationFile string
	for _, file := range files {
		if strings.Contains(file.Name(), filename) &&
			strings.HasSuffix(file.Name(), fmt.Sprintf(".%s.sql", direction)) {
			migrationFile = file.Name()
			break
		}
	}

	if migrationFile == "" {
		return fmt.Errorf("migration file not found for version %d direction %s", version, direction)
	}

	// Ler conteúdo do arquivo
	content, err := os.ReadFile(filepath.Join(mm.migrationPath, migrationFile))
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %v", migrationFile, err)
	}

	// Executar SQL na transação
	if _, err := tx.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute SQL from %s: %v", migrationFile, err)
	}

	return nil
}

// updateSchemaVersionInTx atualiza a tabela schema_migrations na transação
func (mm *MigrationManager) updateSchemaVersionInTx(tx *sql.Tx, version uint, dirty bool) error {
	// A sintaxe correta define explicitamente coluna = valor
	query := `
        UPDATE schema_migrations 
        SET version = $1, dirty = $2
    `

	result, err := tx.Exec(query, version, dirty)
	if err != nil {
		return err
	}

	// ⚠️ IMPORTANTE: Se a tabela estiver vazia (primeira vez rodando),
	// o UPDATE retorna sucesso mas não grava nada (0 linhas afetadas).
	// É boa prática verificar isso:
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		// Opcional: Retornar erro ou fazer um INSERT aqui se for a primeira execução
		query := `
        INSERT INTO schema_migrations (version, dirty)
        VALUES($1, false)`
		_, err = tx.Exec(query, version)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetPendingMigrations retorna migrações pendentes em ordem
func (mm *MigrationManager) GetPendingMigrations() ([]uint, error) {
	currentVersion, _, err := mm.m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return nil, err
	}

	allVersions, err := mm.GetMigrationSequence()
	if err != nil {
		return nil, err
	}

	var pending []uint
	for _, version := range allVersions {
		if version > currentVersion {
			pending = append(pending, version)
		}
	}

	return pending, nil
}

// GetMigrationName obtém o nome da migração baseado na versão
func (mm *MigrationManager) GetMigrationName(version uint) string {
	files, err := os.ReadDir(mm.migrationPath)
	if err != nil {
		return fmt.Sprintf("migration_%d", version)
	}

	versionStr := fmt.Sprintf("%d", version)

	for _, file := range files {
		filename := file.Name()
		if strings.Contains(filename, versionStr) && strings.HasSuffix(filename, ".up.sql") {
			// Extrair o nome descritivo (tudo após a versão e o primeiro _)
			name := strings.TrimSuffix(filename, ".up.sql")
			parts := strings.SplitN(name, "_", 2)
			if len(parts) > 1 {
				return parts[1]
			}
			return name
		}
	}

	return fmt.Sprintf("migration_%d", version)
}

// PrintMigrationSequence mostra todas as migrações em ordem
func (mm *MigrationManager) PrintMigrationSequence() error {
	versions, err := mm.GetMigrationSequence()
	if err != nil {
		return err
	}

	currentVersion, dirty, _ := mm.m.Version()

	fmt.Println("\n📋 Sequência de Migrações:")
	fmt.Println("┌────┬──────────────┬──────────────────────┬──────────┐")
	fmt.Println("│ #  │ Versão       │ Nome                 │ Status   │")
	fmt.Println("├────┼──────────────┼──────────────────────┼──────────┤")

	for i, version := range versions {
		status := "Pendente"
		if version == currentVersion {
			if dirty {
				status = "⚠️ Dirty"
			} else {
				status = "✅ Atual"
			}
		} else if version < currentVersion {
			status = "✅ Aplicada"
		}

		name := mm.GetMigrationName(version)
		// Truncar nome se for muito longo
		if len(name) > 20 {
			name = name[:20] + "..."
		}

		fmt.Printf("│ %-2d │ %-12d │ %-20s │ %-8s │\n",
			i+1, version, name, status)
	}

	fmt.Println("└────┴──────────────┴──────────────────────┴──────────┘")

	if currentVersion == 0 {
		fmt.Printf("📊 Status: Nenhuma migração aplicada\n")
	} else {
		fmt.Printf("📊 Status atual: Versão %d", currentVersion)
		if dirty {
			fmt.Printf(" (dirty)\n")
		} else {
			fmt.Printf(" (clean)\n")
		}
	}

	return nil
}
