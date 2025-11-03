package app

import (
	"deskapp/src/apps/core/view"
	"deskapp/src/internal/utils"
)

type BaseApp struct {
	Name    string
	Version string
	Logger  *utils.Logger
	View    *view.View
	mode    utils.MODE
}

func NewBaseApp(name, version string, logger *utils.Logger, mode utils.MODE) *BaseApp {
	return &BaseApp{
		Name:    name,
		Version: version,
		Logger:  logger,
		View:    view.NewView(mode),
	}
}

func (ba *BaseApp) GetName() string {
	return ba.Name
}

func (ba *BaseApp) GetVersion() string {
	return ba.Version
}

func (ba *BaseApp) GetLogger() *utils.Logger {
	return ba.Logger
}

func (ba *BaseApp) LogInfo(format string, args ...interface{}) {
	ba.Logger.Infof("[%s] "+format, append([]interface{}{ba.Name}, args...)...)
}

func (ba *BaseApp) GetMode() utils.MODE {
	return ba.mode
}
