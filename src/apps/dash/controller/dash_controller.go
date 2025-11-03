package controller

import (
	"deskapp/src/app"
	"deskapp/src/apps/core/controller"
	"deskapp/src/apps/core/view"
	"net/http"
)

type DashController struct {
	*controller.BaseController
}

func NewDashController(app app.AppInterface, view *view.View) *DashController {
	base := controller.NewBaseController(app, view, "dash controller")
	return &DashController{
		BaseController: base,
	}
}

func (c *DashController) Index(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":       "dash",
		"Page":        "dash",
		"ActiveMenu":  "dash",
		"Message":     "Bem-vindo ao app dash",
	}
	c.Render(w, "dash/index.html", data)
}

func (c *DashController) GetDash(w http.ResponseWriter, r *http.Request) {
	c.JSON(w, http.StatusOK, map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"id":   "1",
				"name": "Exemplo dash 1",
			},
			{
				"id":   "2", 
				"name": "Exemplo dash 2",
			},
		},
		"total": 2,
		"app":   "dash",
	})
}

func (c *DashController) CreateDash(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		c.JSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Método não permitido",
		})
		return
	}

	c.JSON(w, http.StatusCreated, map[string]interface{}{
		"message": "dash criado com sucesso",
		"app":     "dash",
		"id":      "12345",
	})
}

func (c *DashController) UpdateDash(w http.ResponseWriter, r *http.Request) {
	c.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "dash atualizado com sucesso",
		"app":     "dash",
	})
}

func (c *DashController) DeleteDash(w http.ResponseWriter, r *http.Request) {
	c.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "dash deletado com sucesso", 
		"app":     "dash",
	})
}

// APIHandler - Exemplo de endpoint de API
func (c *DashController) APIHandler(w http.ResponseWriter, r *http.Request) {
	c.JSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"app":     "dash",
		"version": "1.0.0",
		"data":    "Dados da API do app dash",
	})
}
