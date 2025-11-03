package model

import (
	"fmt"
	"time"
)

type Dash struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// NewDash cria uma nova instância
func NewDash(name string) *Dash {
    now := time.Now()
    return &Dash{
        Name:      name,
        CreatedAt: now,
        UpdatedAt: now,
    }
}

// Métodos do modelo podem ser adicionados aqui
func (m *Dash) Validate() error {
    if m.Name == "" {
        return fmt.Errorf("nome não pode estar vazio")
    }
    return nil
}
