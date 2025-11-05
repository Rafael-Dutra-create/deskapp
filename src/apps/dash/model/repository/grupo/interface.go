package grupo

import (
	"context"
	"deskapp/src/apps/core/model/repository"
	"deskapp/src/apps/dash/model/entities"
)

type IGrupoQueryBuilder = repository.IQueryBuilder[entities.Grupo, *entities.Grupo]

type IGrupoRepository interface {
	// Where retorna o QueryBuilder específico.
	Where(ctx context.Context, queryFragment string, arg any) IGrupoQueryBuilder
}
