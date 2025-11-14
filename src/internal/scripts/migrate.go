package main

import (
	"database/sql"
	"deskapp/src/internal/config"
	"deskapp/src/internal/database"
	"deskapp/src/internal/utils"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MigrationManager struct {
	ScriptBase
	migrationPath string
	m             *migrate.Migrate
	db            *sql.DB
	dsn           string
}

type MigrationRecord struct {
	Version     uint
	Dirty       bool
	Name        *string
	AppliedAt   time.Time
	AppliedBy   string
	Environment string
	Notes       string
}

func (s *MigrationManager) Name() string {
	return "migrate"
}

func (s *MigrationManager) Description() string {
	return "Executa operações de migração do banco de dados"
}

func (s *MigrationManager) Execute(args []string) error {
	cfg := config.NewConfig()
	db, err := database.InitDB(cfg.DBDSN)
	logger := utils.NewLogger()
	if err != nil {
		logger.Errorf("Database URL: %s", cfg.DBDSN)
		return fmt.Errorf("falha ao abrir conexão com DB: %v", err)
	}
	defer db.Close()

	// Definir o caminho das migrações
	migrationsPath := "src/migrations"
	s.migrationPath = migrationsPath

	if err := createMigrationsDir(migrationsPath); err != nil {
		return fmt.Errorf("falha ao criar diretório de migrações: %v", err)
	}

	// Criar gerenciador de migrações
	mm, err := NewMigrationManager(db, migrationsPath, cfg.DBDSN)
	if err != nil {
		return fmt.Errorf("falha ao criar gerenciador de migrações: %v", err)
	}

	// Se não houver argumentos, mostrar menu interativo
	if len(args) == 0 {
		return fmt.Errorf("Nenhum parametro fornecido")
	}

	// Processar comandos via argumentos
	return s.processCommand(mm, args)
}

func NewMigrationManager(db *sql.DB, migrationsPath, dsn string) (*MigrationManager, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %v", err)
	}

	return &MigrationManager{m: m, db: db, migrationPath: migrationsPath, dsn: dsn}, nil
}

// createMigrationsDir cria o diretório de migrações se não existir
func createMigrationsDir(migrationsPath string) error {
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		fmt.Printf("📁 Criando diretório de migrações: %s\n", migrationsPath)
		if err := os.MkdirAll(migrationsPath, 0755); err != nil {
			return err
		}
		fmt.Printf("✅ Diretório criado: %s\n", migrationsPath)
	}
	return nil
}


func (s *MigrationManager) processCommand(mm *MigrationManager, args []string) error {
	command := args[0]

	switch command {
	case "up":
		return mm.Up()
	case "down":
		return mm.Down()
	case "force":
		if len(args) < 2 {
			return fmt.Errorf("número de steps não especificado. Uso: migrate steps <número>")
		}
		steps, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("número de steps inválido: %v", err)
		}
		return mm.Force(steps)
	case "create":
		if len(args) < 2 {
			return fmt.Errorf("nome da migração não especificado. Uso: migrate create <nome>")
		}
		return CreateMigration(s.migrationPath, args[1])
	case "status":
		return mm.PrintStatus()
	default:
		return fmt.Errorf("comando desconhecido: %s. Comandos disponíveis: up, down, force, create, status", command)
	}
}

func (mm *MigrationManager) getMigrationName(version uint) string {
	// Implementar lógica para extrair o nome do arquivo de migração
	// baseado no version number
	files, err := os.ReadDir(mm.migrationPath)
	if err != nil {
		return fmt.Sprintf("migration_%d", version)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), fmt.Sprintf("_%d.", version)) {
			// Extrair nome descritivo do arquivo
			name := strings.TrimSuffix(file.Name(), ".up.sql")
			name = strings.TrimSuffix(name, ".down.sql")
			parts := strings.SplitN(name, "_", 2)
			if len(parts) > 1 {
				return parts[1]
			}
		}
	}

	return fmt.Sprintf("migration_%d", version)
}

func (mm *MigrationManager) Up() error {
    // Obter versão atual
    currentVersion, dirty, err := mm.m.Version()
    if err != nil && err != migrate.ErrNilVersion {
        return fmt.Errorf("failed to get current version: %v", err)
    }

    // Verificar se está dirty
    if dirty {
        return fmt.Errorf("database is in dirty state (version %d). Please clean it first", currentVersion)
    }

    // Obter migrações pendentes
    pendingMigrations, err := mm.GetPendingMigrations()
	 if err != nil {
        return fmt.Errorf("failed to get pending migrations: %v", err)
    }
	appliedsMigrations, err := mm.GetAppliedMigrations()
	var lastMigration uint
	if size := len(appliedsMigrations); size > 0 {
		lastMigration = appliedsMigrations[len(appliedsMigrations)-1]
	} 
	

    if err != nil {
        return fmt.Errorf("failed to get applieds migrations: %v", err)
    }

    if len(pendingMigrations) == 0 {
        fmt.Println("✅ Nenhuma migração pendente")
        return nil
    }

    fmt.Printf("📋 Migrações pendentes encontradas: %d\n", len(pendingMigrations))
    
    // Executar em transação única
    if err := mm.executeMigrationsInTransaction(pendingMigrations,lastMigration, "up"); err != nil {
        return fmt.Errorf("❌ Migration failed: %v", err)
    }

    fmt.Println("✅ Todas as migrações foram aplicadas com sucesso")
    return nil
}


func (mm *MigrationManager) Down() error {
    currentVersion, dirty, err := mm.m.Version()
    if err != nil && err != migrate.ErrNilVersion {
        return fmt.Errorf("failed to get current version: %v", err)
    }

    if err == migrate.ErrNilVersion {
        fmt.Println("✅ Nenhuma migração aplicada")
        return nil
    }

    if dirty {
        return fmt.Errorf("database is in dirty state (version %d). Please clean it first", currentVersion)
    }

    // Obter migrações aplicadas
    appliedMigrations, err := mm.GetAppliedMigrations()
    if err != nil {
        return err
    }

    if len(appliedMigrations) == 0 {
        fmt.Println("✅ Nenhuma migração para reverter")
        return nil
    }

    // A última migração aplicada é a que deve ser revertida
    lastApplied := appliedMigrations[len(appliedMigrations)-1]
    
    fmt.Printf("📋 Revertendo migração: v%d (%s)\n", lastApplied, mm.GetMigrationName(lastApplied))

	var previousVersion uint
    if len(appliedMigrations) > 1 {
        // A versão anterior é a penúltima da lista
        previousVersion = appliedMigrations[len(appliedMigrations)-2]
    } else {
        // Estamos revertendo a única migração, então voltamos para a versão 0
        previousVersion = 0
    }
    
    // Executar apenas UMA migração down em transação
    if err := mm.executeSingleMigrationInTransaction(previousVersion, "down"); err != nil {
        return fmt.Errorf("❌ Migration down failed: %v", err)
    }

    fmt.Println("✅ Migração revertida com sucesso")
    return nil
}

func (mm *MigrationManager) executeSingleMigrationInTransaction(version uint, direction string) error {
    tx, err := mm.db.Begin()
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %v", err)
    }

    var success bool
    defer func() {
        if !success {
            if rbErr := tx.Rollback(); rbErr != nil {
                log.Printf("⚠️ Error during rollback: %v", rbErr)
            }
            fmt.Println("🔻 Rollback executado - alterações revertidas")
        }
    }()

    migrationName := mm.GetMigrationName(version)
    fmt.Printf("🔻 Executando migração %s: v%d (%s)...\n", direction, version, migrationName)

    // Executar arquivo de migração
    if err := mm.executeMigrationFile(tx, version, direction); err != nil {
        return fmt.Errorf("failed to execute migration file v%d: %v", version, err)
    }
	
	

    // Atualizar schema_migrations
    if err := mm.updateSchemaVersionInTx(tx, version, false); err != nil {
        return fmt.Errorf("failed to update schema version v%d: %v", version, err)
    }

    // Commit
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %v", err)
    }

    success = true
    fmt.Printf("✅ Migração %s v%d concluída\n", direction, version)
    return nil
}

// GetAppliedMigrations retorna as migrações já aplicadas em ordem crescente
func (mm *MigrationManager) GetAppliedMigrations() ([]uint, error) {
    currentVersion, _, err := mm.m.Version()
    if err != nil && err != migrate.ErrNilVersion {
        return nil, err
    }

    // Se não há migrações aplicadas
    if err == migrate.ErrNilVersion {
        return []uint{}, nil
    }

    // Obter todas as migrações disponíveis
    allVersions, err := mm.GetMigrationSequence()
    if err != nil {
        return nil, err
    }

    // Filtrar apenas as que foram aplicadas (version <= currentVersion)
    var applied []uint
    for _, version := range allVersions {
        if version <= currentVersion {
            applied = append(applied, version)
        }
    }

    return applied, nil
}

func (mm *MigrationManager) Force(n int) error {

	return mm.m.Force(n)
}

func MaskDSN(dsn string) string {
    // Para PostgreSQL DSN no formato: "postgres://user:password@host:port/database?sslmode=disable"
    if strings.HasPrefix(dsn, "postgresql://") {
        // Remover a senha para exibição segura
        re := regexp.MustCompile(`postgresql://([^:]+):[^@]+@`)
        masked := re.ReplaceAllString(dsn, "postgres://$1:****@")
        return masked
    }
    
    // Para outros formatos, retornar uma versão genérica
    return dsn
}

// Método para obter DSN mascarado
func (mm *MigrationManager) GetMaskedDSN() string {

	return MaskDSN(mm.dsn)
}


func (mm *MigrationManager) PrintStatus() error {
	versions, err := mm.GetMigrationSequence()
    if err != nil {
        return err
    }

    currentVersion, dirty, _ := mm.m.Version()

    fmt.Println("\n📋 Sequência de Migrações:")
    fmt.Println("┌────┬──────────────────┬──────────────────────┬────────────────────┐")
    fmt.Println("│ #  │ Versão           │ Nome                 │ Status             │")
    fmt.Println("├────┼──────────────────┼──────────────────────┼────────────────────┤")

    for i, version := range versions {
        status := "Pendente"
        if version == currentVersion {
            if dirty {
                status = "Dirty"
            } else {
                status = "Atual"
            }
        } else if version < currentVersion {
            status = "Aplicada"
        }

        name := mm.GetMigrationName(version)
        // Truncar nome se for muito longo
        if len(name) > 20 {
            name = name[:20] + "..."
        }

        fmt.Printf("│ %-2d │ %-16d │ %-20s │ %-18s │\n", 
            i+1, version, name, status)
    }

    fmt.Println("└────┴──────────────────┴──────────────────────┴────────────────────┘")

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

func (mm *MigrationManager) SaveMigrationMetadata(version uint, name string, notes string) error {
	query := `
		UPDATE schema_migrations 
		SET migration_name = $1, applied_at = CURRENT_TIMESTAMP, applied_by = $2
		WHERE version = $3
	`

	// Obter informações do ambiente
	appliedBy := os.Getenv("USER")
	if appliedBy == "" {
		appliedBy = "unknown"
	}


	_, err := mm.db.Exec(query, name, appliedBy, version)
	return err
}

func (mm *MigrationManager) GetMigrationHistory() ([]MigrationRecord, error) {
	query := `
		SELECT version, dirty, migration_name, applied_at, applied_by
		FROM migration_logs 
		ORDER BY version DESC
		LIMIT 10
	`

	rows, err := mm.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []MigrationRecord
	for rows.Next() {
		var m MigrationRecord
		err := rows.Scan(&m.Version, &m.Dirty, &m.Name, &m.AppliedAt, &m.AppliedBy,)
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, m)
	}

	return migrations, nil
}
