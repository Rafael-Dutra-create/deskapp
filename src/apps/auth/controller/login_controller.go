package controller

import (
	"deskapp/src/app"
	"deskapp/src/apps/core/controller"
	"deskapp/src/apps/core/view"
	"net/http"
)

type LoginController struct {
	*controller.BaseController
}

func NewLoginController(app app.AppInterface, view *view.View) *LoginController {
	base := controller.NewBaseController(app, view, "login controller")
	return &LoginController{
		BaseController: base,
	}
}

func (c *LoginController) Login(w http.ResponseWriter, r *http.Request) {
	c.JSON(w, http.StatusOK, nil)
}

func (c *LoginController) Logout(w http.ResponseWriter, r *http.Request) {
	c.JSON(w, http.StatusBadRequest, map[string]string{
		"teste": "aqui",
	})
}