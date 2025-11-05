package repository

import (
	"context"
	"database/sql"
	"deskapp/src/apps/core/model/entities"
	"fmt"
)

type BaseRepository[T any, P interface { *T; entities.Entity }] struct {
	db     *sql.DB // Voltamos ao sql.DB padrão!
	table  string
	schema string
}

// getFullTableName é um helper interno para formatar "schema"."table".
func (r *BaseRepository[T, P]) getFullTableName() string {
	if r.schema != "" {
		return fmt.Sprintf(`"%s"."%s"`, r.schema, r.table)
	}
	return fmt.Sprintf(`"%s"`, r.table)
}

// GetDB expõe o executor para uso direto, se necessário.
func (r *BaseRepository[T, P]) GetDB() *sql.DB {
	return r.db
}

func (r *BaseRepository[T, P]) Where(ctx context.Context, queryFragment string, arg any) IQueryBuilder[T, P] {
	return &QueryBuilder[T, P]{
		repo:    r,
		ctx:     ctx,
		wheres:  []string{queryFragment},
		args:    []any{arg},
	}
}

// Select inicia uma nova consulta selecionando colunas específicas.
func (r *BaseRepository[T, P]) Select(ctx context.Context) IQueryBuilder[T, P] {
	return &QueryBuilder[T, P]{
		repo:    r,
		ctx:     ctx,
	}
}

func NewBaseRepository[T any, P interface { *T; entities.Entity }](db *sql.DB, table string, schema string) *BaseRepository[T, P] {
	return &BaseRepository[T, P]{
		db:     db,
		table:  table,
		schema: schema,
	}
}
