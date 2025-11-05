package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestQueryBuilder_First_Success(t *testing.T) {
	repo, mock, ctx := setup(t)
	defer repo.GetDB().Close()

	expectedSQL := `SELECT id, name, age FROM "public"."users" WHERE id = $1 LIMIT 1`

	// Define as linhas que o mock deve retornar
	rows := sqlmock.NewRows([]string{"id", "name", "age"}).
		AddRow(123, "Test User", 30)

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(123).
		WillReturnRows(rows)

	// Executa a query
	result, err := repo.Where(ctx, "id = $1", 123).First()

	if err != nil {
		t.Errorf("QueryBuilder.First falhou: %s", err)
	}

	if result.ID != 123 || result.Name.String != "Test User" {
		t.Errorf("Resultado inesperado: %+v", result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectativas do SQLMock não atendidas: %s", err)
	}
}

func TestQueryBuilder_First_NotFound(t *testing.T) {
	repo, mock, ctx := setup(t)
	defer repo.GetDB().Close()

	expectedSQL := `SELECT id, name, age FROM "public"."users" WHERE id = $1 LIMIT 1`

	// Simula um erro de "não encontrado"
	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows) // O erro padrão do driver

	// Executa a query
	_, err := repo.Where(ctx, "id = $1", 999).First()

	// Verifica se o erro foi o erro customizado do seu QueryBuilder
	if err == nil {
		t.Fatal("Esperava um erro, mas obteve nil")
	}
	if err.Error() != "registro não encontrado" {
		t.Errorf("Esperava 'registro não encontrado', mas obteve: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectativas do SQLMock não atendidas: %s", err)
	}
}

func TestQueryBuilder_Query_Multiple(t *testing.T) {
	repo, mock, ctx := setup(t)
	defer repo.GetDB().Close()

	// Teste de query mais complexa
	expectedSQL := `SELECT id, name, age FROM "public"."users" WHERE age > $1 AND name = $2 ORDER BY name DESC LIMIT 10`

	rows := sqlmock.NewRows([]string{"id", "name", "age"}).
		AddRow(1, "Alice", 30).
		AddRow(2, "Bob", 40)

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(20, "Ativo").
		WillReturnRows(rows)

	// Executa a query
	results, err := repo.Where(ctx, "age > $1", 20).
		And("name = $2", "Ativo").
		OrderBy("name DESC").
		Limit(10).
		Query()

	if err != nil {
		t.Errorf("QueryBuilder.Query falhou: %s", err)
	}

	if len(results) != 2 {
		t.Fatalf("Esperava 2 resultados, obteve %d", len(results))
	}
	if results[0].Name.String != "Alice" {
		t.Errorf("Resultado 0 inesperado: %+v", results[0])
	}
	if results[1].ID != 2 {
		t.Errorf("Resultado 1 inesperado: %+v", results[1])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectativas do SQLMock não atendidas: %s", err)
	}
}