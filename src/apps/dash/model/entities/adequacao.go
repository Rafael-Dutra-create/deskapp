package entities

import (
	"database/sql"
	"deskapp/src/apps/core/model/entities"
)

// Placeholder para importações dinâmicas

// Adequacao representa a tabela adequacao do banco de dados
type Adequacao struct {
	ID               int             `json:"id"`
	Ano              sql.NullInt16   `json:"ano"`
	IdGrupo          int16           `json:"id_grupo"`
	Part             sql.NullInt16   `json:"part"`
	IdModulo         sql.NullInt16   `json:"id_modulo"`
	IdExame          sql.NullInt16   `json:"id_exame"`
	Adequados        sql.NullInt16   `json:"adequados"`
	Possibilidades   sql.NullInt16   `json:"possibilidades"`
	Total            sql.NullInt16   `json:"total"`
	Esp              sql.NullInt16   `json:"esp"`
	Na21             sql.NullInt16   `json:"na21"`
	Subst            sql.NullInt16   `json:"subst"`
	SumNr            sql.NullInt16   `json:"sum_nr"`
	ContEdu          sql.NullInt16   `json:"cont_edu"`
	ContAnalito      sql.NullInt16   `json:"cont_analito"`
	Edu              sql.NullInt16   `json:"edu"`
	ExameCertificado sql.NullString  `json:"exame_certificado"`
	Adeq             sql.NullFloat64 `json:"adeq"`
	PctMinCert       sql.NullFloat64 `json:"pct_min_cert"`
	IdExameGrupo     sql.NullInt16   `json:"id_exame_grupo"`
	Rodadas          sql.NullInt16   `json:"rodadas"`
	MenorEnvio       sql.NullTime    `json:"menor_envio"`
	MaiorEnvio       sql.NullTime    `json:"maior_envio"`
	Elegivel         sql.NullInt16   `json:"elegivel"`
	Certificado      sql.NullInt32   `json:"certificado"`
	IdModExame       int             `json:"id_mod_exame"`
	MaxObs           sql.NullInt16   `json:"max_obs"`
	ObsTotalValidas  sql.NullInt16   `json:"obs_total_validas"`
	ACum             sql.NullInt16   `json:"a_cum"`
	TCum             sql.NullInt16   `json:"t_cum"`
	AEsp             sql.NullInt16   `json:"a_esp"`
	TEsp             sql.NullInt16   `json:"t_esp"`
	AdeqCum          sql.NullInt16   `json:"adeq_cum"`
	AdeqEsp          sql.NullInt16   `json:"adeq_esp"`
}

// Columns retorna a lista de colunas na ordem exata do ScanRow.
func (m *Adequacao) Columns() []string {
	return []string{
		"id",
		"ano",
		"id_grupo",
		"part",
		"id_modulo",
		"id_exame",
		"adequados",
		"possibilidades",
		"total",
		"esp",
		"na21",
		"subst",
		"sum_nr",
		"cont_edu",
		"cont_analito",
		"edu",
		"exame_certificado",
		"adeq",
		"pct_min_cert",
		"id_exame_grupo",
		"rodadas",
		"menor_envio",
		"maior_envio",
		"elegivel",
		"certificado",
		"id_mod_exame",
		"max_obs",
		"obs_total_validas",
		"a_cum",
		"t_cum",
		"a_esp",
		"t_esp",
		"adeq_cum",
		"adeq_esp",
	}
}

// ScanRow implementa a lógica de scan para um DBScanner (*sql.Row ou *sql.Rows).
func (m *Adequacao) ScanRow(row entities.DBScanner) error {
	return row.Scan(
		&m.ID,
		&m.Ano,
		&m.IdGrupo,
		&m.Part,
		&m.IdModulo,
		&m.IdExame,
		&m.Adequados,
		&m.Possibilidades,
		&m.Total,
		&m.Esp,
		&m.Na21,
		&m.Subst,
		&m.SumNr,
		&m.ContEdu,
		&m.ContAnalito,
		&m.Edu,
		&m.ExameCertificado,
		&m.Adeq,
		&m.PctMinCert,
		&m.IdExameGrupo,
		&m.Rodadas,
		&m.MenorEnvio,
		&m.MaiorEnvio,
		&m.Elegivel,
		&m.Certificado,
		&m.IdModExame,
		&m.MaxObs,
		&m.ObsTotalValidas,
		&m.ACum,
		&m.TCum,
		&m.AEsp,
		&m.TEsp,
		&m.AdeqCum,
		&m.AdeqEsp,
	)
}
