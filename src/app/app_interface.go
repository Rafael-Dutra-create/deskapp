package app

import (
	"database/sql"
	"deskapp/src/internal/config"
	"deskapp/src/internal/utils"
	"github.com/gin-gonic/gin"
)

type AppInterface interface {
    GetName() string
    GetVersion() string
    RegisterRoutes(mux *gin.Engine)
    Initialize() error
    GetControllers() []interface{}
	GetConfig() *config.Config
	GetLogger() *utils.Logger
    GetDB() *sql.DB
}