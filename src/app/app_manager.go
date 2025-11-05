// src/internal/app/manager.go
package app

import (
	"deskapp/src/internal/config"
	"deskapp/src/internal/utils"
	"fmt"
	"html/template"
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
	// Cria o renderizador multitemplate
	render := multitemplate.NewRenderer()
	am.router.HTMLRender = render

	

	if am.templateFS == nil {
		am.logger.Warning("TemplateFS é nil - templates não serão carregados")
		return
	}

	am.logger.Info("Iniciando carregamento de templates com multitemplate...")
	// Carrega templates usando o padrão de layouts e includes
	am.loadTemplatesFromFS(render)

	

	// Define o renderizador no router
	am.router.HTMLRender = render

	am.logger.Info("✅ Sistema multitemplate configurado com sucesso!")
}

func (am *AppManager) loadTemplatesFromFS(render multitemplate.Renderer) {
	var templateFiles []string
	err := fs.WalkDir(am.templateFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && (strings.HasSuffix(path, ".html") || strings.HasSuffix(path, ".tmpl")) {
			templateFiles = append(templateFiles, path)
		}
		return nil
	})

	if err != nil {
		am.logger.Errorf("❌ Erro ao buscar templates: %v", err)
		return
	}

	// 1. Separa layouts, componentes e páginas
	layouts := []string{}
	components := []string{}
	pages := []string{}
	baseLayoutFile := "" // O arquivo 'base.html'

	for _, file := range templateFiles {
		// Encontra o 'base.html'
		if file == "templates/base.html" || strings.HasSuffix(file, "/base.html") {
			baseLayoutFile = file
			continue // Não o adicione a nenhuma outra lista
		}

		if strings.Contains(file, "layouts/") {
			layouts = append(layouts, file)
		} else if strings.Contains(file, "components/") {
			components = append(components, file)
		} else {
			pages = append(pages, file)
		}
	}

	if baseLayoutFile == "" {
		am.logger.Error("❌ Erro crítico: 'base.html' não encontrado nos templates.")
		return
	}

	am.logger.Infof("📊 Layout base: %s", baseLayoutFile)
	am.logger.Infof("📊 Layouts adicionais: %d", len(layouts))
	am.logger.Infof("📊 Componentes: %d", len(components))
	am.logger.Infof("📊 Páginas: %d", len(pages))
	am.logger.Infof("📊 Layouts encontrados: %d", len(layouts))
	am.logger.Infof("📊 Componentes encontrados: %d", len(components))
	am.logger.Infof("📊 Páginas encontradas: %d", len(pages))

	// 2. Carrega todos os Layouts e Componentes UMA ÚNICA VEZ
	baseTemplate, err := template.New("base").ParseFS(am.templateFS, baseLayoutFile)
	if err != nil {
		am.logger.Errorf("❌ Erro ao parsear template base '%s': %v", baseLayoutFile, err)
		return
	}

	commonFiles := append(layouts, components...)
	if len(commonFiles) > 0 {
		_, err = baseTemplate.ParseFS(am.templateFS, commonFiles...)
		if err != nil {
			am.logger.Errorf("❌ Erro ao parsear templates comuns: %v", err)
			return
		}
	}
	am.logger.Info("✅ Templates comuns (layouts + componentes) carregados.")

	// 3. Para cada página, CLONA o template base e adiciona o arquivo da página
	templateCount := 0
	for _, page := range pages {
		name := strings.TrimSuffix(filepath.Base(page), ".html")
		name = strings.TrimSuffix(name, ".tmpl")

		clonedTemplate, err := baseTemplate.Clone()
		if err != nil {
			am.logger.Errorf("❌ Erro ao clonar template base para %s: %v", name, err)
			continue
		}

		// Faz o parse APENAS do arquivo da página no clone
		pageTemplate, err := clonedTemplate.ParseFS(am.templateFS, page)
		if err != nil {
			am.logger.Errorf("❌ Erro ao parsear template de página %s (%s): %v", name, page, err)
			continue
		}

		render.Add(name, pageTemplate)
		templateCount++
	}

	am.logger.Infof("🎉 Total de %d páginas registradas com sucesso!", templateCount)
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
	go func() { 
		if am.GetMode() == utils.RELEASE {
			openBrowser(host)
		}
		 }()

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
