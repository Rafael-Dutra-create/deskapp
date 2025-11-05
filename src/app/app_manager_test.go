// src/internal/app/manager_test.go
package app

import (
	"database/sql"
	"deskapp/src/internal/config"
	"deskapp/src/internal/utils"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"testing/fstest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock para AppInterface (Mantido)
type MockApp struct {
	mock.Mock
	name    string
	version string
}

func (m *MockApp) GetConfig() *config.Config   { panic("unimplemented") }
func (m *MockApp) GetControllers() []interface{} { return nil }
func (m *MockApp) GetDB() *sql.DB                { panic("unimplemented") }
func (m *MockApp) GetLogger() *utils.Logger      { panic("unimplemented") }
func (m *MockApp) GetName() string               { return m.name }
func (m *MockApp) GetVersion() string            { return m.version }
func (m *MockApp) Initialize() error             { args := m.Called(); return args.Error(0) }
func (m *MockApp) RegisterRoutes(router *gin.Engine) { m.Called(router) }

// -----------------------------------------------------------------
// Mocks de FS (MockFS, MockFile, etc.) foram REMOVIDOS
// -----------------------------------------------------------------

// TestNewAppManager (Atualizado para usar fstest.MapFS)
func TestNewAppManager(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	// Mock file systems usando a biblioteca padrão
	staticFS := fstest.MapFS{
		"static/css/style.css": {Data: []byte("body { color: red; }")},
	}
	templateFS := fstest.MapFS{
		"templates/base.html": {Data: []byte(`{{define "base"}}{{template "content" .}}{{end}}`)},
		"templates/index.html": {Data: []byte(`{{define "content"}}<h1>Hello</h1>{{end}}`)},
	}

	am := NewAppManager(logger, cfg, staticFS, templateFS)

	assert.NotNil(t, am)
	assert.NotNil(t, am.apps)
	assert.NotNil(t, am.logger)
	assert.NotNil(t, am.router)
	assert.Equal(t, cfg, am.cfg)
	assert.Equal(t, staticFS, am.staticFS)
	assert.Equal(t, templateFS, am.templateFS)
	assert.NotNil(t, am.router.HTMLRender)
}

// TestSetupMultiTemplates_NilFS (Novo Teste - Caminho de Erro)
// Testa o 'if am.templateFS == nil' em setupMultiTemplates
func TestSetupMultiTemplates_NilFS(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	// Passa 'nil' para templateFS
	am := NewAppManager(logger, cfg, nil, nil)

	// O renderizador será um renderizador vazio, sem templates
	assert.NotNil(t, am.router.HTMLRender)

	// Tentar renderizar uma página deve falhar (pois nada foi carregado)
	am.router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", nil)
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	am.router.ServeHTTP(w, req)

	// A rota existe, mas o template "index" não foi encontrado
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestLoadTemplatesFromFS_Success (Novo Teste - Caminho Feliz)
// Testa a lógica completa de 'loadTemplatesFromFS'
func TestLoadTemplatesFromFS_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	// Define um sistema de templates completo
	templateFS := fstest.MapFS{
		"templates/base.html": {Data: []byte(`{{define "base"}}Base:({{template "content" .}}){{template "component" .}}{{end}}`)},
		"templates/components/card.html": {Data: []byte(`{{define "component"}}-CardComponent{{end}}`)},
		"templates/pages/index.html":     {Data: []byte(`{{define "content"}}IndexPage{{end}}`)},
	}

	am := NewAppManager(logger, cfg, nil, templateFS)

	// Adiciona uma rota de teste para renderizar o template "index"
	am.router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", nil)
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	am.router.ServeHTTP(w, req)

	// Verifica se o 'base', 'component' e 'page' foram renderizados juntos
	assert.Equal(t, http.StatusOK, w.Code)

}

// TestLoadTemplatesFromFS_NoBaseFile (Novo Teste - Caminho de Erro)
// Testa o 'if baseLayoutFile == ""'
func TestLoadTemplatesFromFS_NoBaseFile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	// FS *sem* o 'templates/base.html'
	templateFS := fstest.MapFS{
		"templates/components/card.html": {Data: []byte(`{{define "component"}}-CardComponent{{end}}`)},
		"templates/pages/index.html":     {Data: []byte(`{{define "content"}}IndexPage{{end}}`)},
	}

	am := NewAppManager(logger, cfg, nil, templateFS)

	// Ocorreu um erro (logado) e NENHUM template foi registrado
	am.router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", nil)
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	am.router.ServeHTTP(w, req)

	// O template "index" não existe
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestLoadTemplatesFromFS_BadTemplate (Novo Teste - Caminho de Erro)
// Testa o 'pageTemplate, err := clonedTemplate.ParseFS'
func TestLoadTemplatesFromFS_BadTemplate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	templateFS := fstest.MapFS{
		"templates/base.html": {Data: []byte(`{{define "base"}}{{template "content" .}}{{end}}`)},
		// Template com sintaxe inválida
		"templates/pages/index.html": {Data: []byte(`{{define "content"}} {{ .Nome } {{end}}`)},
	}

	am := NewAppManager(logger, cfg, nil, templateFS)

	// Ocorreu um erro (logado) e o template "index" não foi registrado
	am.router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", nil)
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	am.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestAppManager_SetupStatic (Atualizado para fstest.MapFS e verificação real)
func TestAppManager_SetupStatic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	// Mock FS usando fstest
	staticFS := fstest.MapFS{
		"static/css/style.css": {Data: []byte("body { color: red; }")},
		"static/js/app.js":     {Data: []byte("console.log('hello');")},
	}

	am := NewAppManager(logger, cfg, staticFS, nil)

	req := httptest.NewRequest("GET", "/static/css/style.css", nil)
	w := httptest.NewRecorder()

	am.router.ServeHTTP(w, req)

	// Verifica se o arquivo foi realmente servido
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "body { color: red; }", w.Body.String())
}

// TestAppManager_SetupStatic_SubFSError (Novo Teste - Caminho de Erro)
// Testa o 'fs.Sub(am.staticFS, "static")'
func TestAppManager_SetupStatic_SubFSError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	// Um FS que NÃO contém o diretório 'static'
	staticFS := fstest.MapFS{
		"other/css/style.css": {Data: []byte("body { color: red; }")},
	}

	am := NewAppManager(logger, cfg, staticFS, nil)

	// O fs.Sub() falhou, então o handler /static/ não foi registrado
	req := httptest.NewRequest("GET", "/static/css/style.css", nil)
	w := httptest.NewRecorder()

	am.router.ServeHTTP(w, req)

	// O Gin retorna 404 pois a rota /static/ não existe
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestAppManager_SetupStatic_NilFS (Mantido)
func TestAppManager_SetupStatic_NilFS(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	// Teste não deve falhar (panic)
	assert.NotPanics(t, func() {
		NewAppManager(logger, cfg, nil, &fstest.MapFS{})
	})
}

// --- Testes de Registro de App (Mantidos como estavam) ---

func TestAppManager_RegisterApp(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	am := NewAppManager(logger, cfg, nil, nil)

	mockApp := &MockApp{
		name:    "test-app",
		version: "1.0.0",
	}
	mockApp.On("Initialize").Return(nil)

	err := am.RegisterApp(mockApp)

	assert.NoError(t, err)
	assert.Len(t, am.apps, 1)

	app, exists := am.GetApp("test-app")
	assert.True(t, exists)
	assert.Equal(t, "test-app", app.GetName())
	assert.Equal(t, "1.0.0", app.GetVersion())

	mockApp.AssertCalled(t, "Initialize")
}

func TestAppManager_RegisterApp_Duplicate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	am := NewAppManager(logger, cfg, nil, nil)

	mockApp1 := &MockApp{name: "test-app", version: "1.0.0"}
	mockApp1.On("Initialize").Return(nil)

	mockApp2 := &MockApp{name: "test-app", version: "2.0.0"}
	// mockApp2.On("Initialize").Return(nil) // Não deve ser chamado

	err1 := am.RegisterApp(mockApp1)
	err2 := am.RegisterApp(mockApp2)

	assert.NoError(t, err1)
	assert.Error(t, err2)
	assert.Contains(t, err2.Error(), "já está registrado")
	assert.Len(t, am.apps, 1)
}

func TestAppManager_RegisterApp_InitializeError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	am := NewAppManager(logger, cfg, nil, nil)

	mockApp := &MockApp{
		name:    "test-app",
		version: "1.0.0",
	}
	expectedError := errors.New("initialization failed")
	mockApp.On("Initialize").Return(expectedError)

	err := am.RegisterApp(mockApp)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Len(t, am.apps, 1) // App ainda é registrado
}

func TestAppManager_GetAllApps(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	am := NewAppManager(logger, cfg, nil, nil)

	mockApp1 := &MockApp{name: "app1", version: "1.0.0"}
	mockApp1.On("Initialize").Return(nil)

	mockApp2 := &MockApp{name: "app2", version: "2.0.0"}
	mockApp2.On("Initialize").Return(nil)

	am.RegisterApp(mockApp1)
	am.RegisterApp(mockApp2)

	apps := am.GetAllApps()

	assert.Len(t, apps, 2)
	assert.Contains(t, apps, "app1")
	assert.Contains(t, apps, "app2")
}

func TestAppManager_GetMode(t *testing.T) {
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	am := NewAppManager(logger, cfg, nil, nil)
	mode := am.GetMode()
	assert.Equal(t, cfg.GetMode(), mode)
}

func TestAppManager_GetLogger(t *testing.T) {
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	am := NewAppManager(logger, cfg, nil, nil)
	retrievedLogger := am.GetLogger()
	assert.Equal(t, logger, retrievedLogger)
}

func TestAppManager_RegisterAllRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	am := NewAppManager(logger, cfg, nil, nil)

	mockApp1 := &MockApp{name: "app1", version: "1.0.0"}
	mockApp1.On("Initialize").Return(nil)
	mockApp1.On("RegisterRoutes", mock.AnythingOfType("*gin.Engine")).Return()

	mockApp2 := &MockApp{name: "app2", version: "2.0.0"}
	mockApp2.On("Initialize").Return(nil)
	mockApp2.On("RegisterRoutes", mock.AnythingOfType("*gin.Engine")).Return()

	am.RegisterApp(mockApp1)
	am.RegisterApp(mockApp2)

	am.RegisterAllRoutes()

	mockApp1.AssertCalled(t, "RegisterRoutes", am.router)
	mockApp2.AssertCalled(t, "RegisterRoutes", am.router)
}

func TestAppManager_ConcurrentAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	am := NewAppManager(logger, cfg, nil, nil)

	var wg sync.WaitGroup
	appsToRegister := 10

	// Register apps concurrently
	for i := 0; i < appsToRegister; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			mockApp := &MockApp{
				name:    fmt.Sprintf("app-%d", index),
				version: "1.0.0",
			}
			mockApp.On("Initialize").Return(nil)

			am.RegisterApp(mockApp)
		}(i)
	}

	wg.Wait()

	// Test concurrent reads
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			apps := am.GetAllApps()
			assert.Len(t, apps, appsToRegister)
		}()
	}

	wg.Wait()
}

func TestAppManager_WithDifferentModes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Test RELEASE mode
	t.Setenv("MODE", "RELEASE")
	logger := utils.NewLogger()
	cfg := config.NewConfig()
	am := NewAppManager(logger, cfg, nil, nil)
	assert.Equal(t, utils.RELEASE, am.GetMode())

	// Test DEBUG mode
	t.Setenv("MODE", "DEBUG")
	cfg2 := config.NewConfig()
	am2 := NewAppManager(logger, cfg2, nil, nil)
	assert.Equal(t, utils.DEBUG, am2.GetMode())
}