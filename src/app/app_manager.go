// src/internal/app/manager.go (atualizado)
package app

import (
    "deskapp/src/internal/utils"
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "os/exec"
    "runtime"
    "sync"
)

type AppManager struct {
    apps   map[string]AppInterface
    logger *utils.Logger
    mu     sync.RWMutex
    mode   utils.MODE
    mux    *http.ServeMux // Adicione esta linha
}

// Construtor atualizado para criar o mux internamente
func NewAppManager(logger *utils.Logger, mode utils.MODE) *AppManager {
    return &AppManager{
        apps:   make(map[string]AppInterface),
        logger: logger,
        mode:   mode,
        mux:    http.NewServeMux(), // Cria o mux aqui
    }
}

// GetMux para acessar o mux internamente se necessário
func (am *AppManager) GetMux() *http.ServeMux {
    return am.mux
}

// RegisterAllRoutes atualizado para usar o mux interno
func (am *AppManager) RegisterAllRoutes() { // Remove o parâmetro mux
    am.mu.RLock()
    defer am.mu.RUnlock()

    for name, app := range am.apps {
        am.logger.Infof("Registrando rotas para: %s", name)
        app.RegisterRoutes(am.mux) // Usa o mux interno
    }
}

// SetupStatic atualizado para usar o mux interno
func (am *AppManager) SetupStatic() {
    am.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("src/static"))))
}

// Init atualizado para usar o mux interno
func (am *AppManager) Init() {
    port := fmt.Sprintf(":%d", rand.Intn(8000)+1000)
    host := fmt.Sprintf("http://localhost%s", port)

    // Abrir o navegador automaticamente
    go func() { openBrowser(host) }()

    am.logger.Infof("Servidor rodando em %s", host)
    
    // Use o mux interno aqui
    log.Fatal(http.ListenAndServe(port, am.mux))
}

// ... outros métodos permanecem iguais
func (am *AppManager) GetMode() utils.MODE {
    return am.mode
}

func (am *AppManager) GetLogger() *utils.Logger {
    return am.logger
}

func (am *AppManager) RegisterApp(app AppInterface) error {
    am.mu.Lock()
    defer am.mu.Unlock()

    appName := app.GetName()
    if _, exists := am.apps[appName]; exists {
        return fmt.Errorf("app %s já está registrado", appName)
    }

    am.apps[appName] = app
    am.logger.Infof("App registrado: %s v%s", appName, app.GetVersion())

    return app.Initialize()
}

func (am *AppManager) GetApp(name string) (AppInterface, bool) {
    am.mu.RLock()
    defer am.mu.RUnlock()

    app, exists := am.apps[name]
    return app, exists
}

func (am *AppManager) GetAllApps() map[string]AppInterface {
    am.mu.RLock()
    defer am.mu.RUnlock()

    apps := make(map[string]AppInterface)
    for k, v := range am.apps {
        apps[k] = v
    }
    return apps
}

func openBrowser(url string) {
    var cmd string
    var args []string
    switch runtime.GOOS {
    case "windows":
        cmd = "cmd"
        args = []string{"/c", "start"}
    case "darwin":
        cmd = "open"
    default:
        cmd = "xdg-open"
    }
    args = append(args, url)
    exec.Command(cmd, args...).Start()
}