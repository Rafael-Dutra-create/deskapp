package entities

import (
	"database/sql"
	"deskapp/src/apps/core/model/entities"
)

// Placeholder para importações dinâmicas

// Grupo representa a tabela grupo do banco de dados
type Grupo struct {
	Grupo sql.NullString `json:"grupo"`
	ID    int16          `json:"id"`
}

// Columns retorna a lista de colunas na ordem exata do ScanRow.
func (m *Grupo) Columns() []string {
	return []string{
		"grupo",
		"id",
	}
}

// ScanRow implementa a lógica de scan para um DBScanner (*sql.Row ou *sql.Rows).
func (m *Grupo) ScanRow(row entities.DBScanner) error {
	return row.Scan(
		&m.Grupo,
		&m.ID,
	)
}
