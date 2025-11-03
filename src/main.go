package main

import (
	"deskapp/src/app"
	"deskapp/src/apps/core"
	"deskapp/src/apps/dash"
	"deskapp/src/internal/config"
	"deskapp/src/internal/database"
	"deskapp/src/internal/utils"

	"github.com/joho/godotenv"
)

var logger *utils.Logger

func init() {
	err := godotenv.Load()
	logger = utils.NewLogger() 
	if err != nil {
		logger.Warning("Não foi possível ler o .env")
	}
}

func main() {
	cfg := config.NewConfig()
	dbConn, err := database.InitDB(cfg.DBDSN)
	if err != nil {
		logger.Warningf("Falha ao conectar com banco: %v", err)
	}
	defer dbConn.Close()


	
	app := app.NewAppManager(logger, cfg, StaticFS, TemplateFS)
	app.RegisterApp(core.NewCoreApp(logger, cfg))
	app.RegisterApp(dash.NewDashApp(logger, cfg))
	

	app.RegisterAllRoutes()

	app.Init()
}
