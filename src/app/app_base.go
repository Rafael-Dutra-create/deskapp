package app

import (
	"database/sql"
	"deskapp/src/apps/core/view"
	"deskapp/src/internal/config" // ⬅ NOVO IMPORT
	"deskapp/src/internal/database"
	"deskapp/src/internal/utils"
)

type BaseApp struct {
    Name    string
    Version string
    Logger  *utils.Logger
    View    *view.View
    Config  *config.Config // ⬅ SUBSTITUIU 'mode'
}

// NewBaseApp agora aceita *config.Config em vez de utils.MODE
func NewBaseApp(name, version string, logger *utils.Logger, cfg *config.Config) *BaseApp {
    // 💡 O View.NewView agora precisa da informação do modo que está dentro do cfg
    return &BaseApp{
        Name:    name,
        Version: version,
        Logger:  logger,
        View:    view.NewView(cfg.GetMode()), // Use o método GetMode do Config
        Config:  cfg, // Armazena a configuração
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


// 🆕 ADICIONE um método para obter o Config inteiro
func (ba *BaseApp) GetConfig() *config.Config {
    return ba.Config
}

func (ba *BaseApp) GetDB() *sql.DB {
	return database.GetDB()
}