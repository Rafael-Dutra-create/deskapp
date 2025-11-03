package scripts

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type AppConfig struct {
	Name      string
	LowerName string
	UpperName string
	Version   string
}

// CreateApp cria um novo app com estrutura completa
func CreateApp() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("📱 Digite o nome do novo app: ")
	appName, _ := reader.ReadString('\n')
	appName = strings.TrimSpace(appName)

	if appName == "" {
		return fmt.Errorf("nome do app não pode estar vazio")
	}

	config := AppConfig{
		Name:      appName,
		LowerName: strings.ToLower(appName),
		UpperName: cases.Title(language.Portuguese).String(appName),
		Version:   "1.0.0",
	}

	fmt.Printf("🎯 Criando app: %s...\n", config.Name)

	// Criar estrutura de pastas
	if err := createAppStructure(config); err != nil {
		return err
	}

	// Criar arquivos
	if err := createAppFiles(config); err != nil {
		return err
	}

	// Registrar no main.go
	if err := registerAppInMain(config); err != nil {
		return err
	}

	fmt.Printf("✅ App '%s' criado com sucesso!\n", config.Name)
	fmt.Printf("📁 Estrutura criada em: src/apps/%s/\n", config.LowerName)
	fmt.Println("🔧 Lembre-se de implementar a lógica específica do app")

	return nil
}

func createAppStructure(config AppConfig) error {
	basePath := filepath.Join("src")
	baseAppPath := filepath.Join(basePath, "apps", config.LowerName)
	dirs := []string{
		filepath.Join(baseAppPath, "controller"),
		filepath.Join(baseAppPath, "model"),
		filepath.Join(basePath, "templates", config.LowerName),
		filepath.Join(basePath, "static", config.LowerName, "css"),
		filepath.Join(basePath, "static", config.LowerName, "js"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("erro ao criar pasta %s: %v", dir, err)
		}
		fmt.Printf("📂 Criada pasta: %s\n", dir)
	}

	return nil
}

func createAppFiles(config AppConfig) error {

	basePath := filepath.Join("src")
	baseAppPath := filepath.Join(basePath, "apps", config.LowerName)

	// Criar app.go
	if err := createFileFromTemplate(
		filepath.Join(baseAppPath, "app.go"),
		appTemplate,
		config,
	); err != nil {
		return err
	}

	// Criar routes.go
	if err := createFileFromTemplate(
		filepath.Join(baseAppPath, "routes.go"),
		routesTemplate,
		config,
	); err != nil {
		return err
	}

	// Criar controller principal
	if err := createFileFromTemplate(
		filepath.Join(baseAppPath, "controller", fmt.Sprintf("%s_controller.go", config.LowerName)),
		controllerTemplate,
		config,
	); err != nil {
		return err
	}

	// Criar model principal
	if err := createFileFromTemplate(
		filepath.Join(baseAppPath, "model", fmt.Sprintf("%s.go", config.LowerName)),
		modelTemplate,
		config,
	); err != nil {
		return err
	}

	// Criar template base
	if err := createFileFromTemplate(
		filepath.Join(basePath, "templates", config.LowerName, "index.html"),
		templateIndex,
		config,
	); err != nil {
		return err
	}

	// Criar CSS básico
	if err := createFileFromTemplate(
		filepath.Join(basePath, "static", config.LowerName, "css", "style.css"),
		cssTemplate,
		config,
	); err != nil {
		return err
	}

	// Criar JS básico
	if err := createFileFromTemplate(
		filepath.Join(basePath, "static", config.LowerName, "js", "app.js"),
		jsTemplate,
		config,
	); err != nil {
		return err
	}

	return nil
}

func createFileFromTemplate(filePath, templateContent string, config AppConfig) error {
    // Garantir que o diretório existe
    dir := filepath.Dir(filePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("erro ao criar diretório %s: %v", dir, err)
    }

    // Debug específico para arquivos HTML
    if strings.HasSuffix(filePath, ".html") {
        fmt.Printf("🔍 Criando template HTML: %s\n", filePath)
        fmt.Printf("📏 Tamanho do template: %d bytes\n", len(templateContent))
    }

    tmpl, err := template.New(filepath.Base(filePath)).Parse(templateContent)
    if err != nil {
        return fmt.Errorf("erro ao criar template para %s: %v", filePath, err)
    }
    
    file, err := os.Create(filePath)
    if err != nil {
        return fmt.Errorf("erro ao criar arquivo %s: %v", filePath, err)
    }
    defer file.Close()
    
    if err := tmpl.Execute(file, config); err != nil {
        return fmt.Errorf("erro ao executar template em %s: %v", filePath, err)
    }
    
    // Verificar se o arquivo foi escrito corretamente
    if strings.HasSuffix(filePath, ".html") {
        info, _ := os.Stat(filePath)
        if info != nil {
            fmt.Printf("✅ Arquivo HTML criado: %s (%d bytes)\n", filePath, info.Size())
        }
    } else {
        fmt.Printf("📄 Criado arquivo: %s\n", filePath)
    }
    
    return nil
}

func registerAppInMain(config AppConfig) error {
	mainPath := "src/main.go"

	content, err := os.ReadFile(mainPath)
	if err != nil {
		return fmt.Errorf("erro ao ler main.go: %v", err)
	}

	mainContent := string(content)

	// Verificar se o app já está registrado
	if strings.Contains(mainContent, fmt.Sprintf(`apps/%s`, config.LowerName)) {
		fmt.Printf("⚠️  App '%s' já está registrado no main.go\n", config.Name)
		return nil
	}

	// 1. ADICIONAR IMPORT - Estratégia mais robusta
	newImport := fmt.Sprintf(`	"deskapp/src/apps/%s"`, config.LowerName)

	// Encontrar a seção de imports
	lines := strings.Split(mainContent, "\n")
	inImportSection := false
	importAdded := false

	for i, line := range lines {
		// Marcar início da seção de imports
		if strings.Contains(line, `import (`) {
			inImportSection = true
			continue
		}

		// Se estamos na seção de imports
		if inImportSection {
			// Se encontrou o fechamento da seção de imports, adicionar antes
			if strings.Contains(line, `)`) {
				// Inserir o novo import antes do fechamento
				lines[i] = newImport + "\n" + line
				importAdded = true
				break
			}

			// Se encontrou outro import de apps, adicionar após ele
			if strings.Contains(line, `"deskapp/src/apps/`) {
				// Verificar se é o último import de apps
				if i+1 < len(lines) && !strings.Contains(lines[i+1], `"deskapp/src/apps/`) {
					lines[i] = line + "\n" + newImport
					importAdded = true
					break
				}
			}
		}
	}

	// Se não conseguiu adicionar na seção de imports, tentar estratégia alternativa
	if !importAdded {
		fmt.Printf("⚠️  Não foi possível adicionar o import automaticamente. Adicione manualmente:\n%s\n", newImport)
	} else {
		mainContent = strings.Join(lines, "\n")
		fmt.Printf("✅ Import adicionado: %s\n", newImport)
	}

	// 2. ADICIONAR REGISTRO DO APP
	newRegistration := fmt.Sprintf(`	app.RegisterApp(%s.New%sApp(logger, mode))`, config.LowerName, config.UpperName)

	// Procurar pela última ocorrência de RegisterApp
	registrationAdded := false
	lines = strings.Split(mainContent, "\n")

	for i := len(lines) - 1; i >= 0; i-- {
		if strings.Contains(lines[i], "app.RegisterApp(") {
			// Inserir após esta linha
			if i+1 < len(lines) {
				newLines := make([]string, 0)
				newLines = append(newLines, lines[:i+1]...)
				newLines = append(newLines, newRegistration)
				newLines = append(newLines, lines[i+1:]...)
				lines = newLines
			} else {
				lines = append(lines, newRegistration)
			}
			registrationAdded = true
			break
		}
	}

	// Se não encontrou RegisterApp, procurar por SetupStatic
	if !registrationAdded {
		for i, line := range lines {
			if strings.Contains(line, "app.SetupStatic()") {
				// Inserir após SetupStatic
				if i+1 < len(lines) {
					newLines := make([]string, 0)
					newLines = append(newLines, lines[:i+1]...)
					newLines = append(newLines, newRegistration)
					newLines = append(newLines, lines[i+1:]...)
					lines = newLines
				} else {
					lines = append(lines, newRegistration)
				}
				registrationAdded = true
				break
			}
		}
	}

	mainContent = strings.Join(lines, "\n")

	if registrationAdded {
		fmt.Printf("✅ Registro do app adicionado: %s\n", newRegistration)
	} else {
		fmt.Printf("⚠️  Não foi possível adicionar o registro automaticamente. Adicione manualmente:\n%s\n", newRegistration)
	}

	// 3. Escrever o arquivo atualizado
	if err := os.WriteFile(mainPath, []byte(mainContent), 0644); err != nil {
		return fmt.Errorf("erro ao escrever main.go: %v", err)
	}

	fmt.Printf("📝 App '%s' processado no main.go\n", config.Name)
	return nil
}

// Templates
const controllerTemplate = `package controller

import (
	"deskapp/src/app"
	"deskapp/src/apps/core/controller"
	"deskapp/src/apps/core/view"
	"net/http"
)

type {{.UpperName}}Controller struct {
	*controller.BaseController
}

func New{{.UpperName}}Controller(app app.AppInterface, view *view.View) *{{.UpperName}}Controller {
	base := controller.NewBaseController(app, view, "{{.LowerName}} controller")
	return &{{.UpperName}}Controller{
		BaseController: base,
	}
}

func (c *{{.UpperName}}Controller) Index(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":       "{{.Name}}",
		"Page":        "{{.LowerName}}",
		"ActiveMenu":  "{{.LowerName}}",
		"Message":     "Bem-vindo ao app {{.Name}}",
	}
	c.Render(w, "{{.LowerName}}/index.html", data)
}

func (c *{{.UpperName}}Controller) Get{{.UpperName}}(w http.ResponseWriter, r *http.Request) {
	c.JSON(w, http.StatusOK, map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"id":   "1",
				"name": "Exemplo {{.Name}} 1",
			},
			{
				"id":   "2", 
				"name": "Exemplo {{.Name}} 2",
			},
		},
		"total": 2,
		"app":   "{{.LowerName}}",
	})
}

func (c *{{.UpperName}}Controller) Create{{.UpperName}}(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		c.JSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Método não permitido",
		})
		return
	}

	c.JSON(w, http.StatusCreated, map[string]interface{}{
		"message": "{{.Name}} criado com sucesso",
		"app":     "{{.LowerName}}",
		"id":      "12345",
	})
}

func (c *{{.UpperName}}Controller) Update{{.UpperName}}(w http.ResponseWriter, r *http.Request) {
	c.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "{{.Name}} atualizado com sucesso",
		"app":     "{{.LowerName}}",
	})
}

func (c *{{.UpperName}}Controller) Delete{{.UpperName}}(w http.ResponseWriter, r *http.Request) {
	c.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "{{.Name}} deletado com sucesso", 
		"app":     "{{.LowerName}}",
	})
}

// APIHandler - Exemplo de endpoint de API
func (c *{{.UpperName}}Controller) APIHandler(w http.ResponseWriter, r *http.Request) {
	c.JSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"app":     "{{.LowerName}}",
		"version": "1.0.0",
		"data":    "Dados da API do app {{.Name}}",
	})
}
`

const routesTemplate = `package {{.LowerName}}

import (
	"deskapp/src/apps/{{.LowerName}}/controller"
	"net/http"
)

func (a *{{.UpperName}}App) RegisterRoutes(mux *http.ServeMux) {
    controllers := a.GetControllers()
    
    for _, controllerInterface := range controllers {
        switch ctrl := controllerInterface.(type) {
        case *controller.{{.UpperName}}Controller:
            // Rotas básicas
            mux.HandleFunc("/{{.LowerName}}/", ctrl.Index)
            mux.HandleFunc("/{{.LowerName}}/index/", ctrl.Index)
            
            // Rotas CRUD
            mux.HandleFunc("/{{.LowerName}}/get/", ctrl.Get{{.UpperName}})
            mux.HandleFunc("/{{.LowerName}}/create/", ctrl.Create{{.UpperName}})
            mux.HandleFunc("/{{.LowerName}}/update/", ctrl.Update{{.UpperName}})
            mux.HandleFunc("/{{.LowerName}}/delete/", ctrl.Delete{{.UpperName}})
            
            // Rotas de API
            mux.HandleFunc("/api/{{.LowerName}}/", ctrl.APIHandler)
            mux.HandleFunc("/api/{{.LowerName}}/status/", ctrl.APIHandler)
        }
    }
}
`

// Os templates de app, model, etc. permanecem iguais...
const appTemplate = `package {{.LowerName}}

import (
	"deskapp/src/app"
	"deskapp/src/apps/{{.LowerName}}/controller"
	"deskapp/src/internal/utils"
)

type {{.UpperName}}App struct {
    *app.BaseApp
}

func New{{.UpperName}}App(logger *utils.Logger, mode utils.MODE) *{{.UpperName}}App {
    baseApp := app.NewBaseApp("{{.LowerName}}", "{{.Version}}", logger, mode)
    return &{{.UpperName}}App{
        BaseApp: baseApp,
    }
}

func (a *{{.UpperName}}App) Initialize() error {
    a.LogInfo("Inicializando app {{.Name}}")
    return nil
}

func (a *{{.UpperName}}App) GetControllers() []interface{} {
    return []interface{}{
       controller.New{{.UpperName}}Controller(a, a.View),
    }
}
`

const modelTemplate = `package model

import (
	"fmt"
	"time"
)

type {{.UpperName}} struct {
    ID        string    ` + "`" + `json:"id"` + "`" + `
    Name      string    ` + "`" + `json:"name"` + "`" + `
    CreatedAt time.Time ` + "`" + `json:"created_at"` + "`" + `
    UpdatedAt time.Time ` + "`" + `json:"updated_at"` + "`" + `
}

// New{{.UpperName}} cria uma nova instância
func New{{.UpperName}}(name string) *{{.UpperName}} {
    now := time.Now()
    return &{{.UpperName}}{
        Name:      name,
        CreatedAt: now,
        UpdatedAt: now,
    }
}

// Métodos do modelo podem ser adicionados aqui
func (m *{{.UpperName}}) Validate() error {
    if m.Name == "" {
        return fmt.Errorf("nome não pode estar vazio")
    }
    return nil
}
`

const templateIndex = `{{define "title"}}{{.Name}} - DeskApp{{end}}

{{define "page_css"}}
<link rel="stylesheet" href="/static/{{.LowerName}}/css/style.css">
{{end}}

{{define "content"}}
<div class="{{.LowerName}}-container">
    <header class="app-header">
        <h1>🎯 App {{.Name}}</h1>
        <p>Bem-vindo ao app {{.Name}}</p>
    </header>
    
    <main class="app-content">
        <div class="welcome-card">
            <h2>✅ App Criado com Sucesso!</h2>
            <p>Seu app <strong>{{.Name}}</strong> foi criado e está pronto para desenvolvimento.</p>
            
            <div class="core-features">
                <h3>🚀 Integrado com a Estrutura Core</h3>
                <p>Este app está integrado com a estrutura Core do projeto DeskApp.</p>
                <ul>
                    <li><strong>BaseController</strong> - Herda da estrutura base de controllers</li>
                    <li><strong>Sistema de Views</strong> - Templates unificados e organizados</li>
                    <li><strong>Logging Consistente</strong> - Sistema de logging padronizado</li>
                    <li><strong>Gerenciamento de Rotas</strong> - Sistema automático de rotas</li>
                </ul>
            </div>
            
            <div class="next-steps">
                <h3>📝 Próximos Passos para Desenvolvimento:</h3>
                <ul>
                    <li><strong>Implementar Lógica:</strong> Edite <code>app.go</code> para adicionar funcionalidades específicas</li>
                    <li><strong>Configurar Rotas:</strong> Adicione novas rotas em <code>routes.go</code></li>
                    <li><strong>Criar Templates:</strong> Desenvolva templates em <code>templates/{{.LowerName}}/</code></li>
                    <li><strong>Desenvolver Controllers:</strong> Crie controllers específicos em <code>controller/</code></li>
                    <li><strong>Definir Modelos:</strong> Implemente modelos de dados em <code>model/</code></li>
                    <li><strong>Estilizar:</strong> Personalize o CSS em <code>static/{{.LowerName}}/css/</code></li>
                </ul>
            </div>
            
            <div class="app-stats">
                <h3>📊 Estrutura Criada:</h3>
                <div class="stats-grid">
                    <div class="stat-item">
                        <span class="stat-number">6</span>
                        <span class="stat-label">Arquivos Gerados</span>
                    </div>
                    <div class="stat-item">
                        <span class="stat-number">5</span>
                        <span class="stat-label">Pastas Criadas</span>
                    </div>
                    <div class="stat-item">
                        <span class="stat-number">8+</span>
                        <span class="stat-label">Rotas Configuradas</span>
                    </div>
                </div>
            </div>

            <div class="quick-actions">
                <h3>🔗 Ações Rápidas:</h3>
                <div class="actions-grid">
                    <a href="/{{.LowerName}}/" class="action-btn">
                        <span class="action-icon">👀</span>
                        <span>Ver App em Ação</span>
                    </a>
                    <a href="/{{.LowerName}}/get/" class="action-btn">
                        <span class="action-icon">📋</span>
                        <span>Ver Lista</span>
                    </a>
                    <a href="/{{.LowerName}}/create/" class="action-btn">
                        <span class="action-icon">➕</span>
                        <span>Criar Item</span>
                    </a>
                    <a href="/api/{{.LowerName}}/" class="action-btn">
                        <span class="action-icon">🔌</span>
                        <span>Testar API</span>
                    </a>
                </div>
            </div>

            <div class="code-example">
                <h3>💻 Exemplo de Uso:</h3>
                <pre><code>// No controller {{.LowerName}}_controller.go
func (c *{{.UpperName}}Controller) CustomAction(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{
        "Title": "Página Customizada",
        "Data":  "Seus dados aqui",
    }
    c.Render(w, "{{.LowerName}}/custom.html", data)
}</code></pre>
            </div>
        </div>
    </main>
</div>
{{end}}`

const cssTemplate = `/* Estilos para o app {{.Name}} */
.{{.LowerName}}-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.app-header {
    text-align: center;
    margin-bottom: 40px;
    padding-bottom: 20px;
    border-bottom: 2px solid #e0e0e0;
}

.app-header h1 {
    color: #333;
    margin-bottom: 10px;
}

.app-header p {
    color: #666;
    font-size: 18px;
}

.welcome-card {
    background: white;
    border-radius: 10px;
    padding: 30px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    border-left: 4px solid #007bff;
}

.welcome-card h2 {
    color: #007bff;
    margin-bottom: 15px;
}

.next-steps {
    margin-top: 25px;
    padding-top: 20px;
    border-top: 1px solid #eee;
}

.next-steps h3 {
    color: #333;
    margin-bottom: 15px;
}

.next-steps ul {
    list-style-type: none;
    padding: 0;
}

.next-steps li {
    padding: 8px 0;
    border-bottom: 1px solid #f5f5f5;
}

.next-steps li:last-child {
    border-bottom: none;
}

.next-steps code {
    background: #f4f4f4;
    padding: 2px 6px;
    border-radius: 3px;
    font-family: 'Courier New', monospace;
    color: #d63384;
}
`

const jsTemplate = `// JavaScript para o app {{.Name}}
console.log('App {{.Name}} carregado!');

document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM carregado para o app {{.Name}}');
    
    // Adicione a lógica JavaScript do seu app aqui
    
    // Exemplo: interação básica
    const welcomeCard = document.querySelector('.welcome-card');
    if (welcomeCard) {
        welcomeCard.addEventListener('click', function() {
            this.style.transform = 'scale(0.98)';
            setTimeout(() => {
                this.style.transform = 'scale(1)';
            }, 150);
        });
    }
});
`
