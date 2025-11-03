// src/internal/app/manager.go
package app

import (
	"deskapp/src/internal/config"
	"deskapp/src/internal/utils"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

type AppManager struct {
	apps       map[string]AppInterface
	logger     *utils.Logger
	mu         sync.RWMutex
	router     *gin.Engine
	cfg        *config.Config
	staticFS   fs.FS
	templateFS fs.FS
}

func NewAppManager(logger *utils.Logger, cfg *config.Config, staticFS fs.FS, templateFS fs.FS) *AppManager {
	if cfg.GetMode() == utils.RELEASE {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	am := &AppManager{
		apps:       make(map[string]AppInterface),
		logger:     logger,
		cfg:        cfg,
		router:     router,
		staticFS:   staticFS,
		templateFS: templateFS,
	}

	// Configure templates primeiro
	am.setupMultiTemplates()
	// Depois configure arquivos estáticos
	am.SetupStatic()

	return am
}

func (am *AppManager) setupMultiTemplates() {
	if am.templateFS == nil {
		am.logger.Warning("TemplateFS é nil - templates não serão carregados")
		return
	}

	am.logger.Info("Iniciando carregamento de templates com multitemplate...")

	// Cria o renderizador multitemplate
	render := multitemplate.NewRenderer()

	// Carrega templates usando o padrão de layouts e includes
	am.loadTemplatesFromFS(render)

	// Define o renderizador no router
	am.router.HTMLRender = render

	am.logger.Info("✅ Sistema multitemplate configurado com sucesso!")
}

func (am *AppManager) loadTemplatesFromFS(render multitemplate.Renderer) {
	// Encontra todos os arquivos de template
	var templateFiles []string
	err := fs.WalkDir(am.templateFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && (strings.HasSuffix(path, ".html") || strings.HasSuffix(path, ".tmpl")) {
			templateFiles = append(templateFiles, path)
			am.logger.Infof("📄 Encontrado template: %s", path)
		}
		return nil
	})

	if err != nil {
		am.logger.Errorf("❌ Erro ao buscar templates: %v", err)
		return
	}

	// Separa layouts e páginas
	layouts := []string{}
	pages := []string{}

	for _, file := range templateFiles {
		if strings.Contains(file, "layouts/") || file == "templates/base.html" {
			layouts = append(layouts, file)
		} else {
			pages = append(pages, file)
		}
	}

	am.logger.Infof("📊 Layouts encontrados: %d", len(layouts))
	am.logger.Infof("📊 Páginas encontradas: %d", len(pages))

	// Se não encontrou layouts específicos, usa base.html como layout padrão
	if len(layouts) == 0 {
		for _, file := range templateFiles {
			if file == "templates/base.html" {
				layouts = append(layouts, file)
				break
			}
		}
	}

	// Adiciona templates combinando layouts com páginas usando AddFromFS
	templateCount := 0

	for _, page := range pages {
		// Para cada página, combina com todos os layouts
		for _, layout := range layouts {
			// Nome do template é o nome da página sem extensão
			name := strings.TrimSuffix(filepath.Base(page), ".html")
			name = strings.TrimSuffix(name, ".tmpl")

			// Combina layout + página
			files := []string{layout, page}

			// Usa AddFromFS para adicionar do embed.FS
			render.AddFromFS(name, am.templateFS, files...)

			am.logger.Infof("✅ Template registrado: %s → [%s, %s]", name, filepath.Base(layout), filepath.Base(page))
			templateCount++

			// Para cada página, só usa um layout (evita duplicação)
			break
		}
	}

	am.logger.Infof("🎉 Total de templates registrados: %d", templateCount)

	// Debug: verifica os templates registrados
	am.debugRegisteredTemplates(render)
}

func (am *AppManager) debugRegisteredTemplates(render multitemplate.Renderer) {
	am.logger.Info("🔍 Verificando templates registrados no multitemplate...")

	// Tenta acessar os templates via type assertion
	if renderMap, ok := render.(multitemplate.Render); ok {
		am.logger.Infof("📋 Total de templates registrados: %d", len(renderMap))
		for name := range renderMap {
			am.logger.Infof("   - '%s'", name)
		}
	} else {
		am.logger.Error("❌ Não foi possível acessar a lista de templates")
	}
}

func (am *AppManager) SetupStatic() {
	if am.staticFS == nil {
		am.logger.Warning("StaticFS é nil - arquivos estáticos não serão servidos")
		return
	}

	// Cria um sub-filesystem para a pasta static
	staticSubFS, err := fs.Sub(am.staticFS, "static")
	if err != nil {
		am.logger.Errorf("❌ Erro ao criar sub-filesystem para static: %v", err)
		return
	}

	// Debug: listar arquivos estáticos disponíveis
	am.logger.Info("📁 Conteúdo do StaticFS:")
	err = fs.WalkDir(staticSubFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			am.logger.Infof("   📄 %s", path)
		}
		return nil
	})
	if err != nil {
		am.logger.Errorf("❌ Erro ao listar arquivos estáticos: %v", err)
	}

	// Configura o StaticFS com o sub-filesystem
	am.router.StaticFS("/static", http.FS(staticSubFS))
	am.logger.Info("✅ Sistema de arquivos estáticos configurado em /static")

	// Middleware para log de requisições estáticas
	am.router.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/static/") {
			am.logger.Infof("📦 Requisição estática: %s", c.Request.URL.Path)
		}
		c.Next()
	})
}

// ... resto dos métodos permanece igual
func (am *AppManager) RegisterAllRoutes() {
	am.mu.RLock()
	defer am.mu.RUnlock()

	for name, app := range am.apps {
		am.logger.Infof("Registrando rotas para: %s", name)
		app.RegisterRoutes(am.router)
	}
}

func (am *AppManager) Init() {
	host := fmt.Sprintf("http://localhost:%s", am.cfg.Port)

	// Abrir o navegador automaticamente
	go func() { openBrowser(host) }()

	am.logger.Infof("Servidor rodando em %s", host)

	log.Fatal(am.router.Run(fmt.Sprintf(":%s", am.cfg.Port)))
}

func (am *AppManager) GetMode() utils.MODE {
	return am.cfg.GetMode()
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
