package dash

import (
	"deskapp/src/apps/dash/controller"

	"github.com/gin-gonic/gin"
)

func (a *DashApp) RegisterRoutes(router *gin.Engine) {
    controllers := a.GetControllers()
    
    for _, controllerInterface := range controllers {
        switch ctrl := controllerInterface.(type) {
        case *controller.DashController:
            // Rotas básicas
            router.GET("/", ctrl.Index)
            
            // Rotas CRUD
            // mux.HandleFunc("/dash/get/", ctrl.GetDash)
            // mux.HandleFunc("/dash/create/", ctrl.CreateDash)
            // mux.HandleFunc("/dash/update/", ctrl.UpdateDash)
            // mux.HandleFunc("/dash/delete/", ctrl.DeleteDash)
            
            // // Rotas de API
            // mux.HandleFunc("/api/dash/", ctrl.APIHandler)
            // mux.HandleFunc("/api/dash/status/", ctrl.APIHandler)
        }
    }
}
