package auth

import (
	"deskapp/src/app"
	"deskapp/src/apps/auth/controller"
	"deskapp/src/internal/utils"
)

type AuthApp struct {
    *app.BaseApp
}

func NewAuthApp(logger *utils.Logger, mode utils.MODE) *AuthApp {
    baseApp := app.NewBaseApp("auth", "1.0.0", logger, mode)
    return &AuthApp{
        BaseApp: baseApp,
    }
}

func (a *AuthApp) Initialize() error {
    a.LogInfo("Inicializando app de autenticação")
    return nil
}

func (a *AuthApp) GetControllers() []interface{} {
    return []interface{}{
       controller.NewLoginController(a, a.View),
    }
}
