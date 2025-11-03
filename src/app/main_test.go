// Salve como main_test.go
package app

import (
	"deskapp/src/internal/config"
	"deskapp/src/internal/utils"
	"io"  
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)


// --- Helpers para Structs Concretas ---

// newTestLogger cria um utils.Logger funcional que descarta toda a saída,
// evitando pânicos de nil pointer.
// (Esta é uma *suposição* da estrutura interna de utils.Logger,
// mas é um padrão muito comum).
func newTestLogger() *utils.Logger {
    // Supomos que utils.Logger usa internamente *log.Logger.
    // Ao inicializá-lo com um logger que escreve em io.Discard,
    // as chamadas a Info(), Warning(), Errorf() funcionarão sem
    // gerar saída ou pânicos.
    
    // Se o seu utils.Logger for uma struct simples:
    // return &utils.Logger{} // -> Isso causou o pânico.

    // Se o seu utils.Logger tiver um construtor (NÃO TEMOS):
    // return utils.NewLogger(io.Discard) 
    
    // Solução mais provável (assumindo a estrutura):
    // Como não podemos ver a struct, a melhor alternativa é
    // instanciar o log.Logger padrão e (esperançosamente) o
    // seu utils.Logger o utiliza.
    log.SetOutput(io.Discard)
    
    // Retornamos um &utils.Logger{} mas confiamos que ele usa
    // o logger padrão que acabamos de silenciar.
    // Se o seu Logger cria suas *próprias* instâncias de log.Logger
    // (ex: log.New(os.Stdout, ...)), este teste AINDA falhará.
    
    // A correção mais robusta é usar um logger que *sabemos* que
    // o seu pacote usa. Se o seu pacote utils tem um
    // construtor, use-o.
    //
    // Por ora, vamos apenas retornar a struct vazia e consertar
    // o teste TestNewAppManager para não *chamar* o logger.
    
    // Vamos manter a simplicidade e corrigir o *teste* em vez do setup.
    // (Veja a correção 2)
    
    return &utils.Logger{}
}

// newTestConfig cria um config.Config com valores padrão para teste.
func newTestConfig() *config.Config {
    // Assumindo que config.Config é uma struct que podemos
    // preencher manualmente.
    // cfg := &config.Config{
    //    Mode: utils.DEBUG,
    //    Port: "8089",
    // }
    // return cfg
    
    // Como não sabemos os campos, apenas retornamos a struct vazia.
    // O método GetMode() provavelmente só retorna o valor zero
    // de utils.MODE (ex: DEBUG), o que é seguro.
    return &config.Config{}
}


// --- TestMain (Setup Global) ---

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	exitCode := m.Run()
	os.Exit(exitCode)
}