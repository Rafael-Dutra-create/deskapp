package adequacao

import (
	"database/sql"
	"deskapp/src/apps/core/model/repository"
	"deskapp/src/apps/dash/model/entities"
)

// AdequacaoRepository é o repositório para a entidade Adequacao
type AdequacaoRepository struct {
	*repository.BaseRepository[entities.Adequacao, *entities.Adequacao]
}

// NewAdequacaoRepository cria um novo AdequacaoRepository
func NewAdequacaoRepository(db *sql.DB) *AdequacaoRepository {
	base := repository.NewBaseRepository[entities.Adequacao](db, "adequacao", "ep_dw")
	return &AdequacaoRepository{
		BaseRepository: base,
	}
}

