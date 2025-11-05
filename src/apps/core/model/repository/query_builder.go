package repository

import (
	"context"
	"database/sql"
	"deskapp/src/apps/core/model/entities"
	"errors"
	"fmt"
	"strings"
)

// DBScanner define a interface para Scan, implementada por *sql.Row e *sql.Rows.
type DBScanner interface {
	Scan(dest ...any) error
}

type IQueryBuilder[T any, P interface { *T; entities.Entity }] interface {
	And(queryFragment string, arg any) IQueryBuilder[T, P]
	OrderBy(orderBy string) IQueryBuilder[T, P]
	Limit(limit uint64) IQueryBuilder[T, P]
	Offset(offset uint64) IQueryBuilder[T, P]
	First() (*T, error)
	Query() ([]*T, error)
	PrintQuery()
}

type QueryBuilder[T any, P interface { *T; entities.Entity }] struct {
	repo    *BaseRepository[T, P] // Também precisa ser genérico
	ctx     context.Context
	columns []string
	wheres  []string
	args    []any
	orderBy string
	limit   uint64
	offset  uint64
	query string
}

// And adiciona uma condição "AND" à consulta.
func (qb *QueryBuilder[T, P]) And(queryFragment string, arg any) IQueryBuilder[T, P] {
	qb.wheres = append(qb.wheres, queryFragment)
	qb.args = append(qb.args, arg)
	return qb
}


// OrderBy define a cláusula ORDER BY.
func (qb *QueryBuilder[T, P]) OrderBy(orderBy string) IQueryBuilder[T, P] {
	qb.orderBy = orderBy
	return qb
}

// Limit define o LIMIT.
func (qb *QueryBuilder[T, P]) Limit(limit uint64) IQueryBuilder[T, P] {
	qb.limit = limit
	return qb
}

// Offset define o OFFSET.
func (qb *QueryBuilder[T, P]) Offset(offset uint64) IQueryBuilder[T, P] {
	qb.offset = offset
	return qb
}

// buildSelectSQL é um helper interno para montar a string SQL final.
func (qb *QueryBuilder[T, P]) buildSelectSQL() (string, []any) {
	var query strings.Builder
	
	// 1. SELECT
	query.WriteString("SELECT ")
	var model T
	pModel := P(&model)
	query.WriteString(strings.Join(pModel.Columns(), ", "))

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

	qb.query = query.String()

	return query.String(), qb.args
}

func (qb *QueryBuilder[T, P]) PrintQuery() {
	fmt.Println(qb.query)
}

// Query executa a consulta e retorna *sql.Rows (para múltiplos resultados).
func (qb *QueryBuilder[T, P]) Query() ([]*T, error) {
	sql, args := qb.buildSelectSQL()
	rows, err := qb.repo.db.QueryContext(qb.ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*T
	for rows.Next() {
		// 1. Cria um novo destino para cada linha
		dest := new(T)
		
		// 2. Pede para ele se escanear a partir do *sql.Rows
		err := P(dest).ScanRow(rows)
		if err != nil {
			return nil, err // Erro durante o scan da linha
		}
		results = append(results, dest)
	}

	if err = rows.Err(); err != nil {
		return nil, err // Erro pós-iteração (ex: conexão perdida)
	}

	return results, nil
}

// First executa a consulta, adiciona "LIMIT 1" e retorna *sql.Row (para um resultado).
func (qb *QueryBuilder[T, P]) First() (*T, error) {
	if qb.limit == 0 || qb.limit > 1 {
		qb.limit = 1
	}

	query, args := qb.buildSelectSQL()
	row := qb.repo.db.QueryRowContext(qb.ctx, query, args...)

	// 1. Cria o destino (ex: new(Usuario))
	dest := new(T)
	
	// 2. Pede ao destino para "se escanear" a partir do *sql.Row
	//    P(dest) converte &T para P (ex: *Usuario)
	err := P(dest).ScanRow(row)

	if err != nil {
		if err == sql.ErrNoRows {
			// Você pode ter um erro customizado aqui
			return nil, errors.New("registro não encontrado") 
		}
		return nil, err
	}

	return dest, nil
}