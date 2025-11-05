package repository

import (
	"context"
	"database/sql"
	"deskapp/src/apps/core/model/entities"
	"fmt"
	"reflect"
	"strings"
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

/* 
getEntityColumnMap usa reflexão para mapear "nome_da_coluna" -> valor
 Ex: "email" -> "teste@exemplo.com", "id" -> 123
*/
func (r *BaseRepository[T, P]) getEntityColumnMap(entity P) (map[string]any, error) {
	// P é um ponteiro para T (ex: *Usuario), então .Elem() pega a struct (Usuario)
	v := reflect.ValueOf(entity).Elem()
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("entidade não é uma struct")
	}
	t := v.Type()

	colMap := make(map[string]any)
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		// Usa a tag 'json' como o nome da coluna,
		// exatamente como seu 'tablemap' faz
		colName := field.Tag.Get("json")
		if colName != "" {
			colMap[colName] = v.Field(i).Interface()
		}
	}
	return colMap, nil
}

// Insert constrói e executa um INSERT
func (r *BaseRepository[T, P]) Insert(ctx context.Context, entity P) error {
	tableName := r.getFullTableName()

	// 1. Pega os valores da entidade usando reflexão
	colsMap, err := r.getEntityColumnMap(entity)
	if err != nil {
		return err
	}

	// 2. Pega a *ordem* das colunas (do método gerado)
	cols := any(entity).(interface{ Columns() []string }).Columns()

	values := make([]any, 0, len(cols))
	placeholders := make([]string, 0, len(cols))
	into := make([]string, 0, len(cols))

	// 3. Monta os slices de valores e placeholders na ordem correta
	for _, colName := range cols {
		// Ignora 'id' em inserts (assumindo que é auto-increment/default)
		if colName == "id" {
			continue
		}
		into = append(into, colName)
		values = append(values, colsMap[colName])
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(values))) // len(values) é 1-based
	}

	// 4. Constrói a query
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(into, ", "),
		strings.Join(placeholders, ", "),
	)
	fmt.Println(query)
	// 5. Executa
	_, err = r.db.ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("erro ao inserir: %w", err)
	}
	return nil
}

// Update constrói e executa um UPDATE
func (r *BaseRepository[T, P]) Update(ctx context.Context, entity P) error {
	tableName := r.getFullTableName()

	// 1. Pega os valores da entidade
	colsMap, err := r.getEntityColumnMap(entity)
	if err != nil {
		return err
	}

	// 2. Pega a *ordem* das colunas
	cols := any(entity).(interface{ Columns() []string }).Columns()

	setClauses := make([]string, 0)
	values := make([]any, 0)
	var pkValue any
	placeholderIndex := 1 // Placeholders são 1-based

	// 3. Monta a cláusula SET, separando o 'id'
	for _, colName := range cols {
		if colName == "id" {
			pkValue = colsMap[colName]
			continue
		}

		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", colName, placeholderIndex))
		values = append(values, colsMap[colName])
		placeholderIndex++
	}

	// 4. Valida se a PK existe
	if pkValue == nil {
		return fmt.Errorf("entidade sem 'id' para atualização")
	}
	
	// Adiciona o valor do ID no final da lista de argumentos
	values = append(values, pkValue)

	// 5. Constrói a query
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = $%d",
		tableName,
		strings.Join(setClauses, ", "),
		placeholderIndex, // O placeholder final é o do ID
	)

	// 6. Executa
	_, err = r.db.ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("erro ao atualizar: %w", err)
	}
	return nil
}

// Delete constrói e executa um DELETE
func (r *BaseRepository[T, P]) Delete(ctx context.Context, entity P) error {
	tableName := r.getFullTableName()

	// 1. Pega os valores da entidade
	colsMap, err := r.getEntityColumnMap(entity)
	if err != nil {
		return err
	}

	// 2. Valida e pega o valor da PK
	pkValue, ok := colsMap["id"]
	if !ok || pkValue == nil {
		return fmt.Errorf("entidade sem 'id' para exclusão")
	}

	// 3. Constrói a query
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)

	// 4. Executa
	_, err = r.db.ExecContext(ctx, query, pkValue)
	if err != nil {
		return fmt.Errorf("erro ao excluir: %w", err)
	}
	return nil
}


func NewBaseRepository[T any, P interface { *T; entities.Entity }](db *sql.DB, table string, schema string) *BaseRepository[T, P] {
	return &BaseRepository[T, P]{
		db:     db,
		table:  table,
		schema: schema,
	}
}

