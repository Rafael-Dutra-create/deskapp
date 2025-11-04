package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type BaseRepository[T any] struct {
	db     *sql.DB
	table  string
	schema string
}

// getFullTableName é um helper interno para formatar "schema"."table".
func (r *BaseRepository[T]) getFullTableName() string {
	if r.schema != "" {
		return fmt.Sprintf(`"%s"."%s"`, r.schema, r.table)
	}
	return fmt.Sprintf(`"%s"`, r.table)
}

// GetDB expõe o executor para uso direto, se necessário.
func (r *BaseRepository[T]) GetDB() *sql.DB {
	return r.db
}

func (r *BaseRepository[T]) Where(ctx context.Context, queryFragment string, arg any) IQueryBuider[T] {
	return &QueryBuilder[T]{
		repo:    r,
		ctx:     ctx,
		columns: []string{"*"},
		wheres:  []string{queryFragment},
		args:    []any{arg},
	}
}

// Select inicia uma nova consulta selecionando colunas específicas.
func (r *BaseRepository[T]) Select(ctx context.Context, columns ...string) IQueryBuider[T] {
	if len(columns) == 0 {
		columns = []string{"*"}
	}
	return &QueryBuilder{
		repo:    r,
		ctx:     ctx,
		columns: columns,
	}
}

func NewBaseRepository[T any](db *sql.DB, table string, schema string) *BaseRepository[T] {
	return &BaseRepository[T]{
		db:     db,
		table:  table,
		schema: schema,
	}
}
