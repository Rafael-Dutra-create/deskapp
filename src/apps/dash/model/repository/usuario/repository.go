package usuario

import (
	"database/sql"
	"deskapp/src/apps/core/model/repository"
	"deskapp/src/apps/dash/model/entities"
)

// UsuarioRepository é o repositório para a entidade Usuario
type UsuarioRepository struct {
	*repository.BaseRepository[entities.Usuario, *entities.Usuario]
}




// NewUsuarioRepository cria um novo UsuarioRepository
func NewUsuarioRepository(db *sql.DB) *UsuarioRepository {
	base := repository.NewBaseRepository[entities.Usuario](db, "usuario", "public")
	return &UsuarioRepository{
		BaseRepository: base,
	}
}
