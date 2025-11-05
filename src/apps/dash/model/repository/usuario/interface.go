package usuario

import (
	"context"
	"deskapp/src/apps/core/model/repository"
	"deskapp/src/apps/dash/model/entities"
)

// 1. Definimos um alias para o QueryBuilder genérico.
//    Isso é só para deixar o código mais limpo e legível.
//    (IUsuarioQueryBuilder == base.IQueryBuilder[entities.Usuario, *entities.Usuario])
type IUsuarioQueryBuilder = repository.IQueryBuilder[entities.Usuario, *entities.Usuario]

// 2. Definimos a interface do repositório específico.
type IUsuarioRepository interface {
	// --- Métodos "Herdados" do BaseRepository Genérico ---	
	// Where retorna o QueryBuilder específico do usuário.
	Where(ctx context.Context, queryFragment string, arg any) IUsuarioQueryBuilder


	// --- Novos Métodos Específicos do Usuário ---
	GetByEmail(ctx context.Context, email string) (*entities.Usuario, error)
}