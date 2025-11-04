package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)
// DBScanner define a interface para Scan, implementada por *sql.Row e *sql.Rows.
type DBScanner interface {
	Scan(dest ...any) error
}

type IQueryBuider[T any] interface {
	And(queryFragment string, arg any) IQueryBuider[T]
	OrderBy(orderBy string) IQueryBuider[T]
	Limit(limit uint64) IQueryBuider[T]
	Offset(offset uint64) IQueryBuider[T]
	First() (*T, error)
	Query() ([]*T, error)
}


// QueryBuilder armazena o estado da consulta sendo construída.
type QueryBuilder[T any] struct {
	repo    *BaseRepository[T]
	ctx     context.Context
	columns []string // Colunas para o SELECT
	wheres  []string // Condições (ex: "ano = ?")
	args    []any    // Argumentos (ex: 2025)
	orderBy string
	limit   uint64
	offset  uint64
}

// And adiciona uma condição "AND" à consulta.
func (qb *QueryBuilder[T]) And(queryFragment string, arg any) IQueryBuider[T] {
	qb.wheres = append(qb.wheres, queryFragment)
	qb.args = append(qb.args, arg)
	return qb
}

// OrderBy define a cláusula ORDER BY.
func (qb *QueryBuilder[T]) OrderBy(orderBy string) IQueryBuider[T] {
	qb.orderBy = orderBy
	return qb
}

// Limit define o LIMIT.
func (qb *QueryBuilder[T]) Limit(limit uint64) IQueryBuider[T] {
	qb.limit = limit
	return qb
}

// Offset define o OFFSET.
func (qb *QueryBuilder[T]) Offset(offset uint64) IQueryBuider[T] {
	qb.offset = offset
	return qb
}

// buildSelectSQL é um helper interno para montar a string SQL final.
func (qb *QueryBuilder[T]) buildSelectSQL() (string, []any) {
	var query strings.Builder

	// 1. SELECT
	query.WriteString("SELECT ")
	query.WriteString(strings.Join(qb.columns, ", "))

	// 2. FROM
	query.WriteString(" FROM ")
	query.WriteString(qb.repo.getFullTableName())

	// 3. WHERE
	if len(qb.wheres) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(qb.wheres, " AND "))
	}

	// 4. ORDER BY
	if qb.orderBy != "" {
		query.WriteString(" ORDER BY ")
		query.WriteString(qb.orderBy) // CUIDADO: ORDER BY não é parametrizável
	}

	// 5. LIMIT
	if qb.limit > 0 {
		query.WriteString(fmt.Sprintf(" LIMIT %d", qb.limit))
	}

	// 6. OFFSET
	if qb.offset > 0 {
		query.WriteString(fmt.Sprintf(" OFFSET %d", qb.offset))
	}

	return query.String(), qb.args
}

// Query executa a consulta e retorna *sql.Rows (para múltiplos resultados).
func (qb *QueryBuilder[T]) Query() ([]*T, error) {
	sql, args := qb.buildSelectSQL()
	return qb.repo.db.QueryContext(qb.ctx, sql, args...)
}

// First executa a consulta, adiciona "LIMIT 1" e retorna *sql.Row (para um resultado).
func (qb *QueryBuilder[T]) First() (*T, error) {
	if qb.limit == 0 || qb.limit > 1 {
		qb.limit = 1
	}

	sql, args := qb.buildSelectSQL()
	
	// Cria um ponteiro para um novo T (ex: new(entities.Usuario))
	dest := new(T) 

	// sqlx.Get faz a mágica do Scan
	err := qb.repo.db.Query(qb.ctx, dest, sql, args...)
	if err != nil {
		return nil, err 
	}
	
	return dest, nil
}