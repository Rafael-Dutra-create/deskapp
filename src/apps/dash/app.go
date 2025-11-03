package dash

import (
	"deskapp/src/app"
	"deskapp/src/apps/dash/controller"
	"deskapp/src/internal/utils"
)

type DashApp struct {
    *app.BaseApp
}

func NewDashApp(logger *utils.Logger, mode utils.MODE) *DashApp {
    baseApp := app.NewBaseApp("dash", "1.0.0", logger, mode)
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
       controller.NewDashController(a, a.View),
    }
}
