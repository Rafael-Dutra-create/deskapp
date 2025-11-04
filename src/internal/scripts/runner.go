package main

import (
	"deskapp/src/internal/utils"
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/joho/godotenv"
)

// registerScripts registra todos os scripts disponíveis
func registerScripts() {
    // Registra os scripts aqui
    Register(&CreateAppScript{})
    Register(&TableMapScript{})
}
var logger *utils.Logger

func init() {
    logger = utils.NewLogger()
    err := godotenv.Load()
    if err != nil {
        logger.Warning("error no load dotenv")
    }
}


func main() {
    // Registra todos os scripts
    registerScripts()
    
    if err := run(); err != nil {
        log.Printf("Erro: %v\n", err)
        os.Exit(1)
    }
}

func run() error {
    flag.Parse()
    args := flag.Args()
    
    if len(args) == 0 {
        return printUsage()
    }

    command := args[0]
    
    switch command {
    case "help", "-h", "--help":
        return printUsage()
    case "list", "ls":
        return printScriptsList()
    default:
        script := GetScript(command)
        if script == nil {
            return fmt.Errorf("script não encontrado: %s\nUse 'list' para ver scripts disponíveis", command)
        }
        
        // Passa os argumentos restantes para o script
        scriptArgs := []string{}
        if len(args) > 1 {
            scriptArgs = args[1:]
        }
        
        return script.Execute(scriptArgs)
    }
}

func printUsage() error {
    fmt.Printf("Uso: go run main.go <comando> [args...]\n\n")
    fmt.Printf("Comandos:\n")
    fmt.Printf("  help, -h, --help  Mostra esta ajuda\n")
    fmt.Printf("  list, ls          Lista todos os scripts disponíveis\n")
    fmt.Printf("  <script-name>     Executa um script específico\n\n")
    fmt.Printf("Exemplos:\n")
    fmt.Printf("  go run main.go list\n")
    fmt.Printf("  go run main.go create-app --name minha-app\n")
    return nil
}

func printScriptsList() error {
    scripts := ListScripts()
    if len(scripts) == 0 {
        fmt.Println("Nenhum script registrado")
        return nil
    }

    w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
    fmt.Fprintln(w, "NOME\tDESCRIÇÃO")
    
    for _, name := range scripts {
        script := GetScript(name)
        fmt.Fprintf(w, "%s\t%s\n", name, script.Description())
    }
    
    w.Flush()
    return nil
}
