package controller

import (
	"context"
	"deskapp/src/app"
	"deskapp/src/apps/core/controller"
	"deskapp/src/apps/dash/model/repository/usuario"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashController struct {
	*controller.BaseController
}

func NewDashController(app app.AppInterface) *DashController {

	base := controller.NewBaseController(app, "dash controller")
	return &DashController{
		BaseController: base,
		// dashService: dashService,
	}
}

// 1. Index - Substitui w e r por *gin.Context e usa ctx.HTML
func (c *DashController) Index(ctx *gin.Context) {
	data := map[string]interface{}{
        "Title":      "dash",
        "Page":       "dash",
        "ActiveMenu": "dash",
        "Name":       "Dash",
        "LowerName":  "dash",
    }

    resp, err := http.Get("https://estatapi.controllab.com/rodadas/mod/368?total=1&ano[]=2026")

    if err != nil {
        c.GetLogger().Warningf("falha na requisição da API: %v", err)
        data["Message"] = "Erro ao carregar dados da API."
    } else {
        defer resp.Body.Close()

        // 🚨 MUDANÇA CRÍTICA: Decodificar para uma lista (slice) de mapas
        var rodadasData []map[string]interface{} // <-- Novo tipo
        body, err := io.ReadAll(resp.Body)
        
        if err != nil {
            // ... (log de erro de leitura do corpo)
            data["Message"] = "Erro ao ler corpo da resposta."
        } else if err := json.Unmarshal(body, &rodadasData); err != nil { // <-- Decodifica para a lista
            c.GetLogger().Warningf("falha ao decodificar JSON: %v", err)
            data["Message"] = "Erro ao decodificar dados da API."
        } else {
            data["Message"] = fmt.Sprintf("Dados carregados com sucesso! Total: %d", len(rodadasData))
            // 🚨 INJETA A LISTA NO MAPA DE DADOS DO TEMPLATE
            data["Rodadas"] = rodadasData 
        }
    }

    repo := usuario.NewUsuarioRepository(c.GetApp().GetDB())
    usuarios, _ := repo.Select(context.Background()).Query()
    data["User"] = usuarios[0]
	ctx.HTML(http.StatusOK, "dash_index", data)
}

