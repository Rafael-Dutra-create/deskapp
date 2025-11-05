package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"deskapp/src/apps/core/model/entities"

	"github.com/DATA-DOG/go-sqlmock"
)


// MockUser é o nosso 'T'
type MockUser struct {
	ID   int64          `json:"id"`
	Name sql.NullString `json:"name"`
	Age  int            `json:"age"`
}

// MockUser implementa entities.Entity (via *MockUser)

// Columns retorna a lista de colunas na ordem exata do ScanRow.
func (m *MockUser) Columns() []string {
	return []string{
		"id",
		"name",
		"age",
	}
}

// ScanRow implementa a lógica de scan
func (m *MockUser) ScanRow(row entities.DBScanner) error {
	return row.Scan(
		&m.ID,
		&m.Name,
		&m.Age,
	)
}

// --- 2. FUNÇÃO DE SETUP ---

// setup cria um repo mockado para cada teste
func setup(t *testing.T) (
	*BaseRepository[MockUser, *MockUser], // O Repositório concreto
	sqlmock.Sqlmock, // O mock do DB
	context.Context, // Contexto
) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("falha ao criar sqlmock: %s", err)
	}

	// Criamos o repo concreto passando o DB mock e os nomes
	repo := NewBaseRepository[MockUser](db, "users", "public")
	ctx := context.Background()

	return repo, mock, ctx
}

// --- 3. TESTES DO REPOSITORY.GO (Insert, Update, Delete) ---

func TestInsert(t *testing.T) {
	repo, mock, ctx := setup(t)
	defer repo.GetDB().Close()

	user := &MockUser{
		Name: sql.NullString{String: "Test User", Valid: true},
		Age:  30,
	}

	// Query esperada (note que 'id' é pulado, como na sua lógica)
	// Usamos MustCompile para tratar a query como Regexp
	expectedSQL := `INSERT INTO "public"."users" (name, age) VALUES ($1, $2)`
	mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
		WithArgs("Test User", 30).
		WillReturnResult(sqlmock.NewResult(1, 1)) 

	err := repo.Insert(ctx, user)
	if err != nil {
		t.Errorf("Insert falhou: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectativas do SQLMock não atendidas: %s", err)
	}
}

func TestUpdate(t *testing.T) {
	repo, mock, ctx := setup(t)
	defer repo.GetDB().Close()

	user := &MockUser{
		ID:   123,
		Name: sql.NullString{String: "Updated User", Valid: true},
		Age:  31,
	}

	// Query esperada (note a ordem dos placeholders)
	expectedSQL := `UPDATE "public"."users" SET name = $1, age = $2 WHERE id = $3`
	mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
		WithArgs("Updated User", 31, int64(123)).
		WillReturnResult(sqlmock.NewResult(0, 1)) 

	err := repo.Update(ctx, user)
	if err != nil {
		t.Errorf("Update falhou: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectativas do SQLMock não atendidas: %s", err)
	}
}

func TestDelete(t *testing.T) {
	repo, mock, ctx := setup(t)
	defer repo.GetDB().Close()

	user := &MockUser{ID: 123} // Só precisamos do ID

	expectedSQL := `DELETE FROM "public"."users" WHERE id = $1`
	mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
		WithArgs(int64(123)).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 0 = ID (não relevante), 1 = linha afetada

	err := repo.Delete(ctx, user)
	if err != nil {
		t.Errorf("Delete falhou: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectativas do SQLMock não atendidas: %s", err)
	}
}
