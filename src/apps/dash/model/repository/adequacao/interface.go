package adequacao

import (
	"context"
	"deskapp/src/apps/core/model/repository"
	"deskapp/src/apps/dash/model/entities"
)

type IAdequacaoQueryBuilder = repository.IQueryBuilder[entities.Adequacao, *entities.Adequacao]

type IAdequacaoRepository interface {
	// Where retorna o QueryBuilder específico.
	Where(ctx context.Context, queryFragment string, arg any) IAdequacaoQueryBuilder
	Create( entities.Adequacao ) error 
}
