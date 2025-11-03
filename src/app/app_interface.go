package app

import (
	"deskapp/src/internal/utils"
	"net/http"
)

type AppInterface interface {
    GetName() string
    GetVersion() string
    RegisterRoutes(mux *http.ServeMux)
    Initialize() error
    GetControllers() []interface{}
	GetMode() utils.MODE
	GetLogger() *utils.Logger
}