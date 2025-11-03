package controller

// IController define a interface que todos os controllers devem implementar
type IController interface {
    // GetName retorna o nome do controller para logging e identificação
    GetName() string
    
}