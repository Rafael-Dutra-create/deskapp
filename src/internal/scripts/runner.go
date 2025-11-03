package scripts

import (
    "flag"
    "fmt"
)

// Run executa scripts baseado em argumentos de linha de comando
func Run() error {
    flag.Parse()
    args := flag.Args()
    
    if len(args) == 0 {
        return fmt.Errorf("uso: go run scripts/runner.go <comando>\ncomandos disponíveis: create-app")
    }
    
    switch args[0] {
    case "create-app":
        return CreateApp()
    default:
        return fmt.Errorf("comando não reconhecido: %s", args[0])
    }
}