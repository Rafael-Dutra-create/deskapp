package entities

import (
	"database/sql"
	"time"
)

// Placeholder para importações dinâmicas

// DBScanner define a interface para Scan, implementada por *sql.Row e *sql.Rows.
type DBScanner interface {
	Scan(dest ...any) error
}

// Usuario representa a tabela usuario do banco de dados
type Usuario struct {
	Email         sql.NullString `json:"email"`
	Cpf           sql.NullString `json:"cpf"`
	Cnpj          sql.NullString `json:"cnpj"`
	IsAdmin       bool           `json:"is_admin"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	EmailVerified sql.NullTime   `json:"email_verified"`
	ID            string         `json:"id"`
	Image         sql.NullString `json:"image"`
	Name          sql.NullString `json:"name"`
	Password      sql.NullString `json:"password"`
}

// Columns retorna a lista de colunas na ordem exata do ScanRow.
func (m *Usuario) Columns() []string {
	return []string{
		"email",
		"cpf",
		"cnpj",
		"is_admin",
		"created_at",
		"updated_at",
		"email_verified",
		"id",
		"image",
		"name",
		"password",
	}
}

// ScanRow implementa a lógica de scan para um DBScanner (*sql.Row ou *sql.Rows).
func (m *Usuario) ScanRow(row DBScanner) error {
	return row.Scan(
		&m.Email,
		&m.Cpf,
		&m.Cnpj,
		&m.IsAdmin,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.EmailVerified,
		&m.ID,
		&m.Image,
		&m.Name,
		&m.Password,
	)
}
