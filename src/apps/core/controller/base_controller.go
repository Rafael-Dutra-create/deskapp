package controller

import (
	"deskapp/src/app"
	"deskapp/src/internal/utils"
)

// BaseController fornece implementações padrão para todos os controllers
type BaseController struct {
	app  app.AppInterface
	name string
}

// NewBaseController cria uma nova instância do controller base
func NewBaseController(app app.AppInterface,  name string) *BaseController {
	return &BaseController{
		app:  app,
		name: name,
	}
}

// GetName implementa IController
func (bc *BaseController) GetName() string {
	return bc.name
}



// GetLogger retorna o logger para uso nos controllers filhos
func (bc *BaseController) GetLogger() *utils.Logger {
	return bc.app.GetLogger()
}

// GetApp retorna a instância do app para uso nos controllers filhos
func (bc *BaseController) GetApp() app.AppInterface {
	return bc.app
}


// LogInfo método helper para logging consistente
func (bc *BaseController) LogInfo(format string, args ...interface{}) {
	bc.app.GetLogger().Infof("[%s] "+format, append([]interface{}{bc.name}, args...)...)
}

// LogError método helper para logging de erros consistente
func (bc *BaseController) LogError(format string, args ...interface{}) {
	bc.app.GetLogger().Errorf("[%s] "+format, append([]interface{}{bc.name}, args...)...)
}

func (bc *BaseController) LogWarning(format string, args ...interface{}) {
	bc.app.GetLogger().Warningf("[%s] "+format, append([]interface{}{bc.name}, args...)...)
}