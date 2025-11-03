package controller

import "net/http"

// IController define a interface que todos os controllers devem implementar
type IController interface {
    // GetName retorna o nome do controller para logging e identificação
    GetName() string
    
    // Render é o método padrão para renderizar templates
    Render(w http.ResponseWriter, template string, data map[string]interface{})
    
    // JSON envia resposta JSON (método útil comum)
    JSON(w http.ResponseWriter, statusCode int, data interface{})
    
    // Error envia resposta de erro padrão
    Error(w http.ResponseWriter, statusCode int, message string)
}