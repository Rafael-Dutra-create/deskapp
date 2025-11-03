// src/internal/app/manager_test.go
package app

import (
	"database/sql"
	"deskapp/src/internal/config"
	"deskapp/src/internal/utils"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock para AppInterface
type MockApp struct {
	mock.Mock
	name    string
	version string
}

// GetConfig implements AppInterface.
func (m *MockApp) GetConfig() *config.Config {
	panic("unimplemented")
}

// GetControllers implements AppInterface.
func (m *MockApp) GetControllers() []interface{} {
	return nil
}

// GetDB implements AppInterface.
func (m *MockApp) GetDB() *sql.DB {
	panic("unimplemented")
}

// GetLogger implements AppInterface.
func (m *MockApp) GetLogger() *utils.Logger {
	panic("unimplemented")
}

func (m *MockApp) GetName() string {
	return m.name
}

func (m *MockApp) GetVersion() string {
	return m.version
}

func (m *MockApp) Initialize() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockApp) RegisterRoutes(router *gin.Engine) {
	m.Called(router)
}

// Mock para FileSystem
type MockFS struct {
	mock.Mock
	files map[string]string
}

func (m *MockFS) Open(name string) (fs.File, error) {
	if content, exists := m.files[name]; exists {
		return &MockFile{name: name, content: strings.NewReader(content)}, nil
	}
	return nil, os.ErrNotExist
}

func (m *MockFS) ReadDir(name string) ([]fs.DirEntry, error) {
	var entries []fs.DirEntry
	for file := range m.files {
		if strings.HasPrefix(file, name) {
			entries = append(entries, &MockDirEntry{name: file})
		}
	}
	return entries, nil
}

func (m *MockFS) ReadFile(name string) ([]byte, error) {
	if content, exists := m.files[name]; exists {
		return []byte(content), nil
	}
	return nil, os.ErrNotExist
}

func (m *MockFS) Stat(name string) (fs.FileInfo, error) {
	if _, exists := m.files[name]; exists {
		return &MockFileInfo{name: name}, nil
	}
	return nil, os.ErrNotExist
}

// Mock para File
type MockFile struct {
	name    string
	content *strings.Reader
}

func (m *MockFile) Read(p []byte) (int, error) {
	return m.content.Read(p)
}

func (m *MockFile) Close() error {
	return nil
}

func (m *MockFile) Stat() (fs.FileInfo, error) {
	return &MockFileInfo{name: m.name}, nil
}

// Mock para DirEntry
type MockDirEntry struct {
	name string
}

func (m *MockDirEntry) Name() string {
	return m.name
}

func (m *MockDirEntry) IsDir() bool {
	return false
}

func (m *MockDirEntry) Type() fs.FileMode {
	return 0
}

func (m *MockDirEntry) Info() (fs.FileInfo, error) {
	return &MockFileInfo{name: m.name}, nil
}

// Mock para FileInfo
type MockFileInfo struct {
	name string
}

func (m *MockFileInfo) Name() string {
	return m.name
}

func (m *MockFileInfo) Size() int64 {
	return 0
}

func (m *MockFileInfo) Mode() fs.FileMode {
	return 0
}

func (m *MockFileInfo) ModTime() time.Time {
	return time.Now()
}

func (m *MockFileInfo) IsDir() bool {
	return false
}

func (m *MockFileInfo) Sys() interface{} {
	return nil
}

func TestNewAppManager(t *testing.T) {
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	// Mock file systems
	staticFS := &MockFS{
		files: map[string]string{
			"static/css/style.css": "body { color: red; }",
			"static/js/app.js":     "console.log('hello');",
		},
	}

	templateFS := &MockFS{
		files: map[string]string{
			"templates/base.html":         "<html>{{block \"content\" .}}{{end}}</html>",
			"templates/index.html":        `{{define "content"}}<h1>Hello</h1>{{end}}`,
			"templates/layouts/main.html": "<main>{{block \"content\" .}}{{end}}</main>",
		},
	}

	am := NewAppManager(logger, cfg, staticFS, templateFS)

	assert.NotNil(t, am)
	assert.NotNil(t, am.apps)
	assert.NotNil(t, am.logger)
	assert.NotNil(t, am.router)
	assert.Equal(t, cfg, am.cfg)
	assert.Equal(t, staticFS, am.staticFS)
	assert.Equal(t, templateFS, am.templateFS)
}

func TestAppManager_RegisterApp(t *testing.T) {
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	am := NewAppManager(logger, cfg, &MockFS{}, &MockFS{})

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
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	am := NewAppManager(logger, cfg, &MockFS{}, &MockFS{})

	mockApp1 := &MockApp{name: "test-app", version: "1.0.0"}
	mockApp1.On("Initialize").Return(nil)

	mockApp2 := &MockApp{name: "test-app", version: "2.0.0"}
	mockApp2.On("Initialize").Return(nil)

	err1 := am.RegisterApp(mockApp1)
	err2 := am.RegisterApp(mockApp2)

	assert.NoError(t, err1)
	assert.Error(t, err2)
	assert.Contains(t, err2.Error(), "já está registrado")
	assert.Len(t, am.apps, 1)
}

func TestAppManager_RegisterApp_InitializeError(t *testing.T) {
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	am := NewAppManager(logger, cfg, &MockFS{}, &MockFS{})

	mockApp := &MockApp{
		name:    "test-app",
		version: "1.0.0",
	}
	expectedError := errors.New("initialization failed")
	mockApp.On("Initialize").Return(expectedError)

	err := am.RegisterApp(mockApp)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Len(t, am.apps, 1) // App ainda é registrado mesmo com erro na inicialização
}

func TestAppManager_GetAllApps(t *testing.T) {
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	am := NewAppManager(logger, cfg, &MockFS{}, &MockFS{})

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

	am := NewAppManager(logger, cfg, &MockFS{}, &MockFS{})

	mode := am.GetMode()

	assert.Equal(t, cfg.GetMode(), mode)
}

func TestAppManager_GetLogger(t *testing.T) {
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	am := NewAppManager(logger, cfg, &MockFS{}, &MockFS{})

	retrievedLogger := am.GetLogger()

	assert.Equal(t, logger, retrievedLogger)
}

func TestAppManager_RegisterAllRoutes(t *testing.T) {
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	am := NewAppManager(logger, cfg, &MockFS{}, &MockFS{})

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

func TestAppManager_SetupStatic(t *testing.T) {
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	staticFS := &MockFS{
		files: map[string]string{
			"static/css/style.css":   "body { color: red; }",
			"static/js/app.js":       "console.log('hello');",
			"static/images/logo.png": "fake-png-content",
		},
	}

	am := NewAppManager(logger, cfg, staticFS, &MockFS{})

	// Test if static routes are set up by making a request
	req := httptest.NewRequest("GET", "/static/css/style.css", nil)
	w := httptest.NewRecorder()

	am.router.ServeHTTP(w, req)

	// The route should exist (even if file might not be served due to mock)
	assert.NotEqual(t, http.StatusNotFound, w.Code)
}

func TestAppManager_SetupStatic_NilFS(t *testing.T) {
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	am := NewAppManager(logger, cfg, nil, &MockFS{})

	// This should not panic and should log a warning
	assert.NotNil(t, am)
}

func TestAppManager_ConcurrentAccess(t *testing.T) {
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	am := NewAppManager(logger, cfg, &MockFS{}, &MockFS{})

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
	// Test RELEASE mode
	t.Setenv("MODE", "RELEASE")
	logger := utils.NewLogger()
	cfg := config.NewConfig()

	am := NewAppManager(logger, cfg, &MockFS{}, &MockFS{})

	assert.Equal(t, utils.RELEASE, am.GetMode())

	// Test DEBUG mode
	t.Setenv("MODE", "DEBUG")
	cfg2 := config.NewConfig()
	am2 := NewAppManager(logger, cfg2, &MockFS{}, &MockFS{})

	assert.Equal(t, utils.DEBUG, am2.GetMode())
}

// Test helper function
func TestOpenBrowser(t *testing.T) {
	// This is a simple test to ensure the function doesn't panic
	// Actual browser opening is hard to test in CI
	assert.NotPanics(t, func() {
		openBrowser("http://localhost:8006")
	})
}
