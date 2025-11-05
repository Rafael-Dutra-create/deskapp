package controller

import (
	"deskapp/src/app"
	"deskapp/src/internal/utils"
)

// IController define a interface que todos os controllers devem implementar
type IController interface {
    // GetName retorna o nome do controller para logging e identificação
    GetName() string
    GetLogger() *utils.Logger
    GetApp() app.AppInterface 
    LogInfo(format string, args ...interface{})
    LogError(format string, args ...interface{})
    LogWarning(format string, args ...interface{})
    
}