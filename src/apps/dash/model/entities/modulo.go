package model

import (
	"database/sql"
)

// Placeholder para importações dinâmicas

// DBScanner define a interface para Scan, implementada por *sql.Row e *sql.Rows.
type DBScanner interface {
	Scan(dest ...any) error
}

// Modulo representa a tabela modulo do banco de dados
type Modulo struct {
	IdSegmento sql.NullInt16  `json:"id_segmento"`
	IdModulo   int16          `json:"id_modulo"`
	IdArea     sql.NullInt16  `json:"id_area"`
	Area       sql.NullString `json:"area"`
	Segmento   sql.NullString `json:"segmento"`
	Modulo     sql.NullString `json:"modulo"`
	ModuloCert sql.NullString `json:"modulo_cert"`
}

// Columns retorna a lista de colunas na ordem exata do ScanRow.
func (m *Modulo) Columns() []string {
	return []string{
		"id_segmento",
		"id_modulo",
		"id_area",
		"area",
		"segmento",
		"modulo",
		"modulo_cert",
	}
}

// ScanRow implementa a lógica de scan para um DBScanner (*sql.Row ou *sql.Rows).
func (m *Modulo) ScanRow(row DBScanner) error {
	return row.Scan(
		&m.IdSegmento,
		&m.IdModulo,
		&m.IdArea,
		&m.Area,
		&m.Segmento,
		&m.Modulo,
		&m.ModuloCert,
	)
}
