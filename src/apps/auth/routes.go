// src/apps/auth/routes.go
package auth

import (
	"deskapp/src/apps/auth/controller"
	"net/http"
)

func (a *AuthApp) RegisterRoutes(mux *http.ServeMux) {
	controllers := a.GetControllers()

	for _, controllerInterface := range controllers {
		switch ctl := controllerInterface.(type) {
		case *controller.LoginController:
			mux.HandleFunc("/auth/login", ctl.Login)
			mux.HandleFunc("/auth/logout", ctl.Logout)
		}
	
	}
}
