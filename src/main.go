package main

import (
	"deskapp/src/app"
	"deskapp/src/apps/auth"
	"deskapp/src/apps/core"
	"deskapp/src/apps/dash"
	"deskapp/src/internal/utils"
	"flag"
)

var modeArg string

func init() {
	flag.StringVar(&modeArg, "mode", "RELEASE", "DEBUG or RELEASE")
	flag.Parse()
}

func main() {
	mode := utils.RELEASE

	if modeArg == utils.DEBUG.String() {
		mode = utils.DEBUG
	}

	logger := utils.NewLogger() 
	app := app.NewAppManager(logger, mode)
	app.SetupStatic()
	app.RegisterApp(core.NewCoreApp(logger, mode))
	app.RegisterApp(auth.NewAuthApp(logger, mode))
	app.RegisterApp(dash.NewDashApp(logger, mode))
	

	app.RegisterAllRoutes()

	app.Init()
}
