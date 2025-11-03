package controller

import (
	"deskapp/src/app"
	"deskapp/src/apps/core/view"
	"deskapp/src/internal/utils"
	"encoding/json"
	"net/http"
)

// BaseController fornece implementações padrão para todos os controllers
type BaseController struct {
	app  app.AppInterface
	view *view.View
	name string
}

// NewBaseController cria uma nova instância do controller base
func NewBaseController(app app.AppInterface, view *view.View, name string) *BaseController {
	return &BaseController{
		app:  app,
		view: view,
		name: name,
	}
}

// GetName implementa IController
func (bc *BaseController) GetName() string {
	return bc.name
}

// Render implementa IController - renderiza templates com dados
func (bc *BaseController) Render(w http.ResponseWriter, template string, data map[string]interface{}) {
	if bc.app.GetMode() == utils.DEBUG {
		bc.app.GetLogger().Infof("[%s] Renderizando template: %s", bc.name, template)
	}

	// Adiciona dados comuns a todos os templates, se necessário
	if data == nil {
		data = make(map[string]interface{})
	}

	// Dados que podem ser úteis em todos os templates
	data["AppName"] = "Meu App"
	data["ControllerName"] = bc.name

	bc.view.Render(w, template, data)
}

// JSON implementa IController - envia resposta JSON
func (bc *BaseController) JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	if bc.app.GetMode() == utils.DEBUG {
		bc.app.GetLogger().Infof("[%s] Enviando resposta JSON - Status: %d", bc.name, statusCode)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		bc.app.GetLogger().Errorf("[%s] Erro ao codificar JSON: %v", bc.name, err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
	}
}

// Error implementa IController - envia resposta de erro padrão
func (bc *BaseController) Error(w http.ResponseWriter, statusCode int, message string) {
	bc.app.GetLogger().Warningf("[%s] Erro %d: %s", bc.name, statusCode, message)

	errorResponse := map[string]interface{}{
		"error": map[string]interface{}{
			"code":       statusCode,
			"message":    message,
			"controller": bc.name,
		},
	}

	bc.JSON(w, statusCode, errorResponse)
}

// GetLogger retorna o logger para uso nos controllers filhos
func (bc *BaseController) GetLogger() *utils.Logger {
	return bc.app.GetLogger()
}

// GetApp retorna a instância do app para uso nos controllers filhos
func (bc *BaseController) GetApp() app.AppInterface {
	return bc.app
}

// GetView retorna a instância da view para uso nos controllers filhos
func (bc *BaseController) GetView() *view.View {
	return bc.view
}

// LogInfo método helper para logging consistente
func (bc *BaseController) LogInfo(format string, args ...interface{}) {
	bc.app.GetLogger().Infof("[%s] "+format, append([]interface{}{bc.name}, args...)...)
}

// LogError método helper para logging de erros consistente
func (bc *BaseController) LogError(format string, args ...interface{}) {
	bc.app.GetLogger().Errorf("[%s] "+format, append([]interface{}{bc.name}, args...)...)
}
