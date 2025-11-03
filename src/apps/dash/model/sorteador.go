package model

import (
	"database/sql"
)

// Sorteador representa a tabela sorteador do banco de dados
type Sorteador struct {
	IdSorteador int `json:"id_sorteador"`
	Nome string `json:"nome"`
	Host string `json:"host"`
	Port int `json:"port"`
	IdEquipamento sql.NullInt32 `json:"id_equipamento"`
}


// Validate (exemplo)
func (m *Sorteador) Validate() error {
    // TODO: Adicionar regras de validação
    // Ex: if m.Name == "" {
    //     return fmt.Errorf("nome não pode estar vazio")
    // }
    return nil
}
