package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type BaseRepository struct {
	db     *sql.DB
	table  string
	schema string
}

// getFullTableName é um helper interno para formatar "schema"."table".
func (r *BaseRepository) getFullTableName() string {
	// Usamos aspas para proteger contra nomes de tabelas com
	// caracteres especiais ou palavras-chave SQL.
	if r.schema != "" {
		return fmt.Sprintf(`"%s"."%s"`, r.schema, r.table)
	}
	return fmt.Sprintf(`"%s"`, r.table)
}

// GetDB expõe o executor para uso direto, se necessário.
func (r *BaseRepository) GetDB() *sql.DB {
	return r.db
}

func (r *BaseRepository) Where(ctx context.Context, queryFragment string, arg any) *QueryBuilder {
	return &QueryBuilder{
		repo:    r,
		ctx:     ctx,
		columns: []string{"*"},
		wheres:  []string{queryFragment},
		args:    []any{arg},
	}
}

// Select inicia uma nova consulta selecionando colunas específicas.
func (r *BaseRepository) Select(ctx context.Context, columns ...string) *QueryBuilder {
	if len(columns) == 0 {
		columns = []string{"*"}
	}
	return &QueryBuilder{
		repo:    r,
		ctx:     ctx,
		columns: columns,
	}
}

func NewBaseRepository(db *sql.DB, table string, scheme string) *BaseRepository {
	return &BaseRepository{
		db: db,
		table: table,
		schema: scheme,
	}
}
