package main

// Registry registro global de scripts
var Registry = make(map[string]IScript)

// Register registra um script no registry
func Register(script IScript) {
    Registry[script.Name()] = script
}

// GetScript retorna um script pelo nome
func GetScript(name string) IScript {
    return Registry[name]
}

// ListScripts retorna lista de scripts disponíveis
func ListScripts() []string {
    var names []string
    for name := range Registry {
        names = append(names, name)
    }
    return names
}
