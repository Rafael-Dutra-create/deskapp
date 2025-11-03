package core

import (
	"deskapp/src/app"
	"deskapp/src/internal/config"
	"deskapp/src/internal/utils"
)

type CoreApp struct {
	*app.BaseApp
}

// GetControllers implements app.AppInterface.
func (a *CoreApp) GetControllers() []interface{} {
	return []interface{}{}
}

// Initialize implements app.AppInterface.
func (a *CoreApp) Initialize() error {
	a.LogInfo("Inicializando core")
	return nil
}

func NewCoreApp(logger *utils.Logger, cfg *config.Config) *CoreApp {
	baseApp := app.NewBaseApp("core", "1.0.0", logger, cfg)
	return &CoreApp{
		BaseApp: baseApp,
	}
}
