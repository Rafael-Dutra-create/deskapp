package usuario

import (
	"context"
	"deskapp/src/apps/core/model/repository"
)


type IUserRepository interface {
    // Note que Where e Select agora retornam a interface específica.
    Where(ctx context.Context, queryFragment string, arg any) repository.IQueryBuider
    Select(ctx context.Context, columns ...string) repository.IQueryBuider
    
    // Aqui você adicionaria outros métodos, como Create, Update, Delete
    // Ex: Create(ctx context.Context, user *entity.User) error
}