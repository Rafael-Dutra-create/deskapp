package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// CreateAppScript implementação do script create-app
type CreateAppScript struct {
    ScriptBase
}

func (s *CreateAppScript) Name() string {
    return "create-app"
}

func (s *CreateAppScript) Description() string {
    return "Cria uma nova aplicação com estrutura básica"
}

func (s *CreateAppScript) Execute(args []string) error {
	return CreateApp()
}


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
		filepath.Join(baseAppPath, "model", "entities"),
		filepath.Join(baseAppPath, "model", "repository"),
		filepath.Join(baseAppPath, "model", "action"),
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
		appTemplate, // ATUALIZADO
		config,
	); err != nil {
		return err
	}

	// Criar routes.go
	if err := createFileFromTemplate(
		filepath.Join(baseAppPath, "routes.go"),
		routesTemplate, // ATUALIZADO
		config,
	); err != nil {
		return err
	}

	// Criar controller principal
	if err := createFileFromTemplate(
		filepath.Join(baseAppPath, "controller", fmt.Sprintf("%s_controller.go", config.LowerName)),
		controllerTemplate, // ATUALIZADO
		config,
	); err != nil {
		return err
	}

	// Criar template base
	// ATUALIZADO: Nome do arquivo mudou de index.html para {{.LowerName}}_index.html
	if err := createFileFromTemplate(
		filepath.Join(basePath, "templates", config.LowerName, fmt.Sprintf("%s_index.html", config.LowerName)),
		templateIndex, // Sem mudanças no conteúdo, mas o nome do arquivo sim
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

	var contentToWrite []byte
	var err error

	// Se for HTML, não use o motor de template, pois ele remove os blocos {{define}}
	// Apenas faça substituições simples.
	if strings.HasSuffix(filePath, ".html") {
		tempContent := strings.ReplaceAll(templateContent, "{{.Name}}", config.Name)
		tempContent = strings.ReplaceAll(tempContent, "{{.LowerName}}", config.LowerName)
		tempContent = strings.ReplaceAll(tempContent, "{{.UpperName}}", config.UpperName)
		tempContent = strings.ReplaceAll(tempContent, "{{.Version}}", config.Version)
		contentToWrite = []byte(tempContent)

	} else {
		// Para arquivos .go, .css, .js, use o motor de template padrão
		tmpl, terr := template.New(filepath.Base(filePath)).Parse(templateContent)
		if terr != nil {
			return fmt.Errorf("erro ao criar template para %s: %v", filePath, terr)
		}

		var buf bytes.Buffer
		if err = tmpl.Execute(&buf, config); err != nil {
			return fmt.Errorf("erro ao executar template em %s: %v", filePath, err)
		}
		contentToWrite = buf.Bytes()
	}

	// Escrever o conteúdo no arquivo
	if err = os.WriteFile(filePath, contentToWrite, 0644); err != nil {
		return fmt.Errorf("erro ao criar arquivo %s: %v", filePath, err)
	}

	fmt.Printf("📄 Criado arquivo: %s\n", filePath)
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
	// ATUALIZADO: Troca 'mode' por 'cfg'
	newRegistration := fmt.Sprintf(`	app.RegisterApp(%s.New%sApp(logger, cfg))`, config.LowerName, config.UpperName)

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

	// Se não encontrou RegisterApp, procurar por SetupStatic (ou outro ponto de referência)
	if !registrationAdded {
		for i, line := range lines {
			// Tente encontrar um ponto de referência, ex: router.Run
			// Ajuste "router.Run" se seu ponto de injeção for outro
			if strings.Contains(line, "router.Run") || strings.Contains(line, "app.SetupStatic()") {
				// Inserir antes desta linha
				newLines := make([]string, 0)
				newLines = append(newLines, lines[:i]...)
				newLines = append(newLines, newRegistration)
				newLines = append(newLines, lines[i:]...)
				lines = newLines
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

// =================================================================
// TEMPLATES ATUALIZADOS
// =================================================================

// ATUALIZADO:
// - Usa *gin.Context em vez de http.ResponseWriter/Request
// - Usa ctx.HTML() e ctx.JSON()
// - Chama o template "{{.LowerName}}_index"
// - Remove verificação de método (o Gin cuida disso no roteamento)
const controllerTemplate = `package controller

import (
	"deskapp/src/app"
	"deskapp/src/apps/core/controller"
	"deskapp/src/apps/core/view"
	"net/http"
	
	"github.com/gin-gonic/gin"
)

type {{.UpperName}}Controller struct {
	*controller.BaseController
}

func New{{.UpperName}}Controller(app app.AppInterface, view *view.View) *{{.UpperName}}Controller {
	base := controller.NewBaseController(app, "{{.LowerName}} controller")
	return &{{.UpperName}}Controller{
		BaseController: base,
	}
}

func (c *{{.UpperName}}Controller) Index(ctx *gin.Context) {
	data := map[string]interface{}{
		"Title":       "{{.Name}}",
		"Page":        "{{.LowerName}}",
		"ActiveMenu":  "{{.LowerName}}",
		"Message":     "Bem-vindo ao app {{.Name}}",
		"Name":        "{{.Name}}",
        "LowerName":   "{{.LowerName}}",
	}
	// ATENÇÃO: O nome do template agora é "{{.LowerName}}_index"
	// e usamos ctx.HTML, não c.Render
	ctx.HTML(http.StatusOK, "{{.LowerName}}_index", data)
}
`

// ATUALIZADO:
// - Usa *gin.Engine em vez de *http.ServeMux
// - Usa router.GET, router.POST, etc.
const routesTemplate = `package {{.LowerName}}

import (
	"deskapp/src/apps/{{.LowerName}}/controller"
	
	"github.com/gin-gonic/gin"
)

func (a *{{.UpperName}}App) RegisterRoutes(router *gin.Engine) {
    controllers := a.GetControllers()

	dashGroup := router.Group("/{{.LowerName}}")
    
    for _, controllerInterface := range controllers {
        switch ctrl := controllerInterface.(type) {
        case *controller.{{.UpperName}}Controller:
            // Rotas básicas
            dashGroup.GET("/", ctrl.Index)
        }
    }
}
`

// ATUALIZADO:
// - Recebe *config.Config em vez de utils.MODE
// - Passa cfg para NewBaseApp
const appTemplate = `package {{.LowerName}}

import (
	"deskapp/src/app"
	"deskapp/src/apps/{{.LowerName}}/controller"
	"deskapp/src/internal/config" // Import adicionado
	"deskapp/src/internal/utils"
)

type {{.UpperName}}App struct {
    *app.BaseApp
}

func New{{.UpperName}}App(logger *utils.Logger, cfg *config.Config) *{{.UpperName}}App {
    baseApp := app.NewBaseApp("{{.LowerName}}", "{{.Version}}", logger, cfg)
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

// SEM MUDANÇAS NO CONTEÚDO
// (O nome do arquivo gerado mudou para {{.LowerName}}_index.html)
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
            <h2>✅ App Criado com Sucesso! (Modelo Gin)</h2>
            <p>Seu app <strong>{{.Name}}</strong> foi criado e está pronto para desenvolvimento.</p>
            
            <div class="core-features">
                <h3>🚀 Integrado com a Estrutura Core (Gin)</h3>
                <p>Este app está integrado com a estrutura Core do projeto DeskApp.</p>
                <ul>
                    <li><strong>Roteamento Gin</strong> - Rotas definidas com <code>router.GET</code>, <code>router.POST</code>, etc.</li>
                    <li><strong>Handlers Gin</strong> - Controllers usam <code>*gin.Context</code></li>
                    <li><strong>Renderização Gin</strong> - Respostas com <code>ctx.HTML()</code> e <code>ctx.JSON()</code></li>
                    <li><strong>Configuração Centralizada</strong> - App recebe <code>*config.Config</code></li>
                </ul>
            </div>
            
            <div class="next-steps">
                <h3>📝 Próximos Passos para Desenvolvimento:</h3>
                <ul>
                    <li><strong>Implementar Lógica:</strong> Edite <code>app.go</code> para adicionar funcionalidades específicas</li>
                    <li><strong>Configurar Rotas:</strong> Adicione novas rotas em <code>routes.go</code></li>
                    <li><strong>Criar Templates:</strong> Desenvolve templates em <code>templates/{{.LowerName}}/</code></li>
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
                        <span class="stat-number">1+</span>
                        <span class="stat-label">Rotas Configuradas</span>
                    </div>
                </div>
            </div>

            <div class="code-example">
                <h3>💻 Exemplo de Uso (Gin):</h3>
                <pre><code>// No controller {{.LowerName}}_controller.go
func (c *{{.UpperName}}Controller) CustomAction(ctx *gin.Context) {
    data := map[string]interface{}{
        "Title": "Página Customizada",
        "Data":  "Seus dados aqui",
    }
    ctx.HTML(http.StatusOK, "{{.LowerName}}_custom", data)
}</code></pre>
            </div>
        </div>
    </main>
</div>
{{end}}`

// SEM MUDANÇAS
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

.next-steps, .core-features, .app-stats, .quick-actions, .code-example {
    margin-top: 25px;
    padding-top: 20px;
    border-top: 1px solid #eee;
}

.next-steps h3, .core-features h3, .app-stats h3, .quick-actions h3, .code-example h3 {
    color: #333;
    margin-bottom: 15px;
}

.next-steps ul, .core-features ul {
    list-style-type: none;
    padding: 0;
}

.next-steps li, .core-features li {
    padding: 8px 0;
    border-bottom: 1px solid #f5f5f5;
}

.next-steps li:last-child, .core-features li:last-child {
    border-bottom: none;
}

.next-steps code {
    background: #f4f4f4;
    padding: 2px 6px;
    border-radius: 3px;
    font-family: 'Courier New', monospace;
    color: #d63384;
}

.stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    gap: 15px;
    text-align: center;
}

.stat-item {
    background: #f9f9f9;
    padding: 15px;
    border-radius: 5px;
}

.stat-item .stat-number {
    display: block;
    font-size: 2em;
    font-weight: bold;
    color: #007bff;
}

.stat-item .stat-label {
    font-size: 0.9em;
    color: #555;
}

.actions-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 15px;
}

.action-btn {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 12px 15px;
    background: #e9f5ff;
    border: 1px solid #b6dfff;
    border-radius: 5px;
    text-decoration: none;
    color: #0056b3;
    font-weight: 500;
    transition: background-color 0.2s, box-shadow 0.2s;
}

.action-btn:hover {
    background-color: #dcf0ff;
    box-shadow: 0 2px 4px rgba(0,0,0,0.05);
}

.action-icon {
    font-size: 1.2em;
}

.code-example pre {
    background: #2d2d2d;
    color: #f1f1f1;
    padding: 15px;
    border-radius: 5px;
    overflow-x: auto;
}

.code-example code {
    font-family: 'Courier New', monospace;
}
`

// SEM MUDANÇAS
const jsTemplate = `// JavaScript para o app {{.Name}}
console.log('App {{.Name}} (Gin) carregado!');

document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM carregado para o app {{.Name}}');
    
    // Adicione a lógica JavaScript do seu app aqui
    
    // Exemplo: interação básica
    const welcomeCard = document.querySelector('.welcome-card');
    if (welcomeCard) {
        welcomeCard.addEventListener('mouseenter', function() {
            this.style.boxShadow = '0 8px 12px rgba(0, 0, 0, 0.15)';
        });
		welcomeCard.addEventListener('mouseleave', function() {
            this.style.boxShadow = '0 4px 6px rgba(0, 0, 0, 0.1)';
        });
    }
});
`
