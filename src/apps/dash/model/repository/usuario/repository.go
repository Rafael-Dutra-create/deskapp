package usuario

import (
	"database/sql"
	"deskapp/src/apps/core/model/repository"
)



// UsuarioRepository é o repositório para a entidade Usuario
type UsuarioRepository struct {
	*repository.BaseRepository
}




// NewUsuarioRepository cria um novo UsuarioRepository
func NewUsuarioRepository(db *sql.DB) *UsuarioRepository {
	base := repository.NewBaseRepository(db, "usuario", "")
	return &UsuarioRepository{
		BaseRepository: base,
	}
}
