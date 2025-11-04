package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"deskapp/src/internal/config"
	"deskapp/src/internal/database"
	"deskapp/src/internal/utils"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type TableMapScript struct {
	ScriptBase
}

func (s *TableMapScript) Name() string {
	return "tablemap"
}

func (s *TableMapScript) Description() string {
	return "Cria uma struct de uma tabela do banco"
}

func (s *TableMapScript) Execute(args []string) error {
	return MapTableToStruct()
}

// ColumnInfo armazena os metadados de uma coluna do DB
type ColumnInfo struct {
	ColumnName string
	DataType   string
	IsNullable string
}

// StructField armazena as propriedades do campo Go gerado
type StructField struct {
	GoName   string // Nome em CamelCase (ex: CreatedAt)
	GoType   string // Tipo Go (ex: sql.NullTime)
	JSONName string // Nome em snake_case (ex: created_at)
}

// StructConfig é passado para o template para gerar o arquivo
type StructConfig struct {
	AppName             string // ex: dash
	ModelName           string // ex: User
	TableName           string // ex: users
	PackageName         string // ex: model
	SchemaName          string // ex: public
	EntitiesPackagePath string
	Fields              []StructField // Lista de campos
}

// MapTableToStruct é a função principal que executa o script
func MapTableToStruct() error {
	cfg := config.NewConfig()
	db, err := database.InitDB(cfg.DBDSN)
	logger := utils.NewLogger()
	if err != nil {
		logger.Errorf("Database URL: %s", cfg.DBDSN)
		return fmt.Errorf("falha ao abrir conexão com DB: %v", err)
	}
	reader := bufio.NewReader(os.Stdin)

	// 1. Coletar informações do usuário
	fmt.Print("📦 Nome do App (ex: dash): ")
	appName, _ := reader.ReadString('\n')
	appName = strings.TrimSpace(appName)

	fmt.Print("📜 Schema (ex: public): ")
	schemaName, _ := reader.ReadString('\n')
	schemaName = strings.TrimSpace(schemaName)

	fmt.Print("🧾 Nome da Tabela (ex: users): ")
	tableName, _ := reader.ReadString('\n')
	tableName = strings.TrimSpace(tableName)

	defer db.Close()

	// 3. Inspecionar a Tabela
	columns, err := inspectTable(db, schemaName, tableName)
	if err != nil {
		return fmt.Errorf("falha ao inspecionar tabela: %v", err)
	}

	if len(columns) == 0 {
		return fmt.Errorf("tabela '%s.%s' não encontrada ou está vazia", schemaName, tableName)
	}

	fmt.Printf("🔍 Encontradas %d colunas. Gerando struct...\n", len(columns))

	// 4. Montar configuração do Struct
	config := StructConfig{
		AppName:     appName,
		ModelName:   snakeToCamel(tableName), // ex: users -> User
		TableName:   tableName,
		PackageName: "entities",
		Fields:      make([]StructField, 0),
	}

	imports := make(map[string]bool)

	for _, col := range columns {
		goType := mapPostgresTypeToGoType(col.DataType, col.IsNullable)
		goName := snakeToCamel(col.ColumnName)

		// Adicionar imports necessários
		if strings.HasPrefix(goType, "sql.") {
			imports["database/sql"] = true
		}
		if goType == "time.Time" || goType == "sql.NullTime" {
			imports["time"] = true
		}

		config.Fields = append(config.Fields, StructField{
			GoName:   goName,
			GoType:   goType,
			JSONName: col.ColumnName,
		})
	}

	// 5. Gerar o arquivo a partir do template
	modelFileName := fmt.Sprintf("%s.go", tableName)
	targetPath := filepath.Join("src", "apps", appName, "model", "entities", modelFileName)

	if err := generateModelFile(targetPath, config, imports); err != nil {
		return fmt.Errorf("falha ao gerar arquivo de model: %v", err)
	}

	capitalizedModelName := titler.String(config.ModelName)

	config.ModelName = capitalizedModelName
	// Define o caminho de importação das entidades
	entitiesPackagePath := filepath.Join("deskapp/src/apps", appName, "model", "entities")
	// Garante barras no padrão Go (linux) e não Windows
	config.EntitiesPackagePath = strings.ReplaceAll(entitiesPackagePath, "\\", "/")

	fmt.Printf("✅ Struct '%s' gerado com sucesso em: %s\n", config.ModelName, targetPath)

	repoFileName := fmt.Sprintf("%s_repository.go", tableName)
	repoTargetPath := filepath.Join("src", "apps", appName, "model", "repository", repoFileName)

	if err := generateRepositoryFile(repoTargetPath, config); err != nil {
		return fmt.Errorf("falha ao gerar arquivo de repository: %v", err)
	}
	return nil
}

// inspectTable busca os metadados das colunas
func inspectTable(db *sql.DB, schemaName, tableName string) ([]ColumnInfo, error) {
	query := `
	SELECT column_name, data_type, is_nullable
	FROM information_schema.columns
	WHERE table_schema = $1
	  AND table_name = $2
	ORDER BY ordinal_position;
	`

	rows, err := db.Query(query, schemaName, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var col ColumnInfo
		if err := rows.Scan(&col.ColumnName, &col.DataType, &col.IsNullable); err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}
	return columns, nil
}

// mapPostgresTypeToGoType converte tipos PG para Go (incluindo nulos)
func mapPostgresTypeToGoType(pgType string, isNullable string) string {
	isNullableBool := strings.ToUpper(isNullable) == "YES"

	switch strings.ToLower(pgType) {
	case "character varying", "varchar", "text", "character", "char", "bpchar":
		if isNullableBool {
			return "sql.NullString"
		}
		return "string"

	case "integer", "int", "int4":
		if isNullableBool {
			return "sql.NullInt32"
		}
		return "int"

	case "smallint", "int2":
		if isNullableBool {
			return "sql.NullInt16"
		}
		return "int16"

	case "bigint", "int8":
		if isNullableBool {
			return "sql.NullInt64"
		}
		return "int64"

	case "boolean", "bool":
		if isNullableBool {
			return "sql.NullBool"
		}
		return "bool"

	case "numeric", "decimal", "real", "float4", "double precision", "float8":
		if isNullableBool {
			return "sql.NullFloat64"
		}
		return "float64"

	case "timestamp", "timestamp without time zone", "timestamp with time zone", "date", "time":
		if isNullableBool {
			return "sql.NullTime"
		}
		return "time.Time"

	case "json", "jsonb":
		if isNullableBool {
			// json.RawMessage pode ser nulo por padrão
			return "json.RawMessage"
		}
		return "json.RawMessage" // Ou []byte

	case "uuid":
		if isNullableBool {
			// Pode requerer uma lib (google/uuid) ou tratar como sql.NullString
			return "sql.NullString"
		}
		return "string"

	default:
		return "interface{}"
	}
}

// snakeToCamel converte snake_case para CamelCase
func snakeToCamel(s string) string {
	var result strings.Builder
	capitalizeNext := true

	for _, r := range s {
		if r == '_' {
			capitalizeNext = true
		} else {
			if capitalizeNext {
				result.WriteRune(unicode.ToUpper(r))
				capitalizeNext = false
			} else {
				result.WriteRune(r)
			}
		}
	}
	// Tratar IDs
	if result.String() == "Id" {
		return "ID"
	}
	return result.String()
}

// Titler para capitalizar o nome do modelo
var titler = cases.Title(language.Portuguese)

// generateModelFile cria o arquivo .go final
func generateModelFile(targetPath string, config StructConfig, imports map[string]bool) error {

	// ⬇️ === TEMPLATE ATUALIZADO ===
	// Adicionamos o método Columns()
	const modelTemplate = `package {{.PackageName}}

// <IMPORT_BLOCK> // Placeholder para importações dinâmicas

// DBScanner define a interface para Scan, implementada por *sql.Row e *sql.Rows.
type DBScanner interface {
	Scan(dest ...any) error
}

// {{.ModelName}} representa a tabela {{.TableName}} do banco de dados
type {{.ModelName}} struct {
{{- range .Fields}}
	{{.GoName}} {{.GoType}} ` + "`" + `json:"{{.JSONName}}"` + "`" + `
{{- end}}
}

// Columns retorna a lista de colunas na ordem exata do ScanRow.
func (m *{{.ModelName}}) Columns() []string {
	return []string{
{{- range .Fields}}
		"{{.JSONName}}",
{{- end}}
	}
}

// ScanRow implementa a lógica de scan para um DBScanner (*sql.Row ou *sql.Rows).
func (m *{{.ModelName}}) ScanRow(row DBScanner) error {
	return row.Scan(
{{- range .Fields}}
		&m.{{.GoName}},
{{- end}}
	)
}
`
	// --- O resto da sua função permanece exatamente igual ---

	// As importações reais necessárias...
	actualImports := map[string]bool{
		"database/sql": true,
	}

	for _, field := range config.Fields {
		if strings.Contains(field.GoType, " time.") || strings.Contains(field.GoType, "sql.NullTime") {
			actualImports["time"] = true
		}
		if strings.Contains(field.GoType, " sql.") {
			actualImports["database/sql"] = true
		}
		if strings.Contains(field.GoType, " json.") {
			actualImports["encoding/json"] = true
		}
	}

	// Criar o conteúdo das importações dinamicamente
	var importStr strings.Builder
	importStr.WriteString("import (\n")
	for imp := range actualImports {
		importStr.WriteString(fmt.Sprintf("\t\"%s\"\n", imp))
	}
	importStr.WriteString(")\n")

	// Substituir o placeholder <IMPORT_BLOCK> pelo bloco de importações
	finalTemplate := strings.Replace(modelTemplate, "// <IMPORT_BLOCK>", importStr.String(), 1)

	// Capitalizar ModelName
	config.ModelName = titler.String(config.ModelName)

	tmpl, err := template.New("model").Parse(finalTemplate)
	if err != nil {
		return fmt.Errorf("erro ao parsear template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return fmt.Errorf("erro ao executar template: %v", err)
	}

	// Formata o código gerado
	formattedSource, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Printf("⚠️  Aviso: falha ao formatar o código gerado: %v\n", err)
		formattedSource = buf.Bytes()
	}

	// Garantir que o diretório existe
	dir := filepath.Dir(targetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório %s: %v", dir, err)
	}

	// Escrever o arquivo formatado
	if err := os.WriteFile(targetPath, formattedSource, 0644); err != nil {
		return fmt.Errorf("erro ao escrever arquivo: %v", err)
	}

	return nil
}

// generateRepositoryFile cria o arquivo .go do repositório
// <<< INÍCIO DA FUNÇÃO ATUALIZADA >>>

// generateRepositoryFile cria o arquivo .go do repositório especializado
func generateRepositoryFile(targetPath string, config StructConfig) error {
	// NOTA: O import do BaseRepository parece fixo com base nos seus arquivos.
	// Ajuste "deskapp/src/apps/core/model/repository" se este caminho for dinâmico.
	const repositoryTemplate = `package repository

import (
	"database/sql"
	"deskapp/src/apps/core/model/repository"
)

// {{.ModelName}}Repository é o repositório para a entidade {{.ModelName}}
type {{.ModelName}}Repository struct {
	*repository.BaseRepository
}

// New{{.ModelName}}Repository cria um novo {{.ModelName}}Repository
func New{{.ModelName}}Repository(db *sql.DB) *{{.ModelName}}Repository {
	base := repository.NewBaseRepository(db, "{{.TableName}}", "{{.SchemaName}}")
	return &{{.ModelName}}Repository{
		BaseRepository: base,
	}
}
`
	// O ModelName já deve vir capitalizado da função MapTableToStruct

	tmpl, err := template.New("repository").Parse(repositoryTemplate)
	if err != nil {
		return fmt.Errorf("erro ao parsear template do repositório: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return fmt.Errorf("erro ao executar template do repositório: %v", err)
	}

	// Formata o código gerado
	formattedSource, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Printf("⚠️  Aviso: falha ao formatar o código do repositório gerado: %v\n", err)
		formattedSource = buf.Bytes()
	}

	// Garantir que o diretório existe
	dir := filepath.Dir(targetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório %s: %v", dir, err)
	}

	// Escrever o arquivo formatado
	if err := os.WriteFile(targetPath, formattedSource, 0644); err != nil {
		return fmt.Errorf("erro ao escrever arquivo do repositório: %v", err)
	}

	return nil
}