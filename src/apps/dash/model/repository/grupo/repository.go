package grupo

import (
	"database/sql"
	"deskapp/src/apps/core/model/repository"
	"deskapp/src/apps/dash/model/entities"
)

// GrupoRepository é o repositório para a entidade Grupo
type GrupoRepository struct {
	*repository.BaseRepository[entities.Grupo, *entities.Grupo]
}

// NewGrupoRepository cria um novo GrupoRepository
func NewGrupoRepository(db *sql.DB) *GrupoRepository {
	base := repository.NewBaseRepository[entities.Grupo](db, "grupo", "ep_dw")
	return &GrupoRepository{
		BaseRepository: base,
	}
}
