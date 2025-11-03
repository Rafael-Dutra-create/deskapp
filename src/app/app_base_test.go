// Salve como app_base_test.go
package app

import (
	"deskapp/src/internal/config"
	"deskapp/src/internal/utils"
	"testing"
)

func TestNewBaseApp(t *testing.T) {
	// Setup
	mockCfg := &config.Config{} // Use o MockConfig se o Config real for complexo
	
	// Assumindo que utils.Logger pode ser nil ou um logger real
	testLogger := &utils.Logger{} // Ajuste conforme necessário

	b := NewBaseApp("TestApp", "1.0.0", testLogger, mockCfg)

	if b == nil {
		t.Fatal("NewBaseApp retornou nil")
	}
	if b.Name != "TestApp" {
		t.Errorf("Esperado Name='TestApp', obtido='%s'", b.Name)
	}
	if b.Version != "1.0.0" {
		t.Errorf("Esperado Version='1.0.0', obtido='%s'", b.Version)
	}
	if b.Logger != testLogger {
		t.Error("Logger não foi atribuído corretamente")
	}
	if b.Config != mockCfg {
		t.Error("Config não foi atribuído corretamente")
	}
	// Não podemos testar o b.View sem o pacote view, mas confiamos que NewView foi chamado
	if b.View == nil {
		t.Error("View é nil, NewView falhou ou não foi chamado")
	}
}

func TestBaseAppGetters(t *testing.T) {
	// Setup
	mockCfg := &config.Config{} 
	testLogger := &utils.Logger{}

	b := NewBaseApp("GetterApp", "v2", testLogger, mockCfg)

	if name := b.GetName(); name != "GetterApp" {
		t.Errorf("GetName(): esperado 'GetterApp', obtido '%s'", name)
	}
	if version := b.GetVersion(); version != "v2" {
		t.Errorf("GetVersion(): esperado 'v2', obtido '%s'", version)
	}
	if logger := b.GetLogger(); logger != testLogger {
		t.Error("GetLogger(): não retornou o logger correto")
	}
	if cfg := b.GetConfig(); cfg != mockCfg {
		t.Error("GetConfig(): não retornou o config correto")
	}
}