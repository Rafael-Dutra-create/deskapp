package dash

import (
	"deskapp/src/app"
	"deskapp/src/apps/dash/controller"
	"deskapp/src/internal/config"
	"deskapp/src/internal/utils"
)

type DashApp struct {
    *app.BaseApp
}

func NewDashApp(logger *utils.Logger, cfg *config.Config) *DashApp {
    baseApp := app.NewBaseApp("dash", "1.0.0", logger, cfg)
    return &DashApp{
        BaseApp: baseApp,
    }
}

func (a *DashApp) Initialize() error {
    a.LogInfo("Inicializando app dash")
    return nil
}

func (a *DashApp) GetControllers() []interface{} {
    return []interface{}{
       controller.NewDashController(a),
    }
}
