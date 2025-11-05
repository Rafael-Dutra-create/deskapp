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

// (Structs ScriptBase, ColumnInfo, StructField permanecem iguais)
type TableMapScript struct {
	ScriptBase
}

func (s *TableMapScript) Name() string { return "tablemap" }
func (s *TableMapScript) Description() string {
	return "Cria uma struct de uma tabela do banco"
}
func (s *TableMapScript) Execute(args []string) error { return MapTableToStruct() }

type ColumnInfo struct {
	ColumnName string
	DataType   string
	IsNullable string
}
type StructField struct {
	GoName   string
	GoType   string
	JSONName string
}

// StructConfig foi atualizada para incluir os novos campos
type StructConfig struct {
	AppName               string // ex: dash
	ModelName             string // ex: User (Capitalizado)
	TableName             string // ex: users
	PackageName           string // ex: entities (para o modelo)
	SchemaName            string // ex: public
	EntitiesPackagePath   string // ex: deskapp/src/apps/dash/model/entities
	RepositoryPackageName string // ex: usuario (minúsculo)
	Fields                []StructField
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
	fmt.Print("📦 Nome do App: ")
	appName, _ := reader.ReadString('\n')
	appName = strings.TrimSpace(appName)

	fmt.Print("📜 Schema (public): ")
	schemaName, _ := reader.ReadString('\n')
	schemaName = strings.TrimSpace(schemaName)
	if schemaName == "" {
		schemaName = "public"
	}

	fmt.Print("🧾 Nome da Tabela: ")
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

	fmt.Printf("🔍 Encontradas %d colunas. Gerando arquivos...\n", len(columns))

	// 4. Montar configuração do Struct
	modelName := snakeToCamel(tableName)
	capitalizedModelName := titler.String(modelName)
	entitiesPackagePath := filepath.Join("deskapp/src/apps", appName, "model", "entities")

	config := StructConfig{
		AppName:     appName,
		ModelName:   capitalizedModelName, // ex: Usuario
		TableName:   tableName,            // ex: usuario
		SchemaName:  schemaName,
		PackageName: "entities", // Pacote para o arquivo de entidade
		// Garante barras no padrão Go (linux)
		EntitiesPackagePath:   strings.ReplaceAll(entitiesPackagePath, "\\", "/"),
		RepositoryPackageName: strings.ToLower(tableName), // ex: usuario
		Fields:                make([]StructField, 0),
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

	// === 5. GERAR ARQUIVO DE ENTIDADE (MODELO) ===
	modelFileName := fmt.Sprintf("%s.go", tableName)
	modelTargetPath := filepath.Join("src", "apps", appName, "model", "entities", modelFileName)

	if err := generateModelFile(modelTargetPath, config, imports); err != nil {
		return fmt.Errorf("falha ao gerar arquivo de model: %v", err)
	}
	fmt.Printf("✅ Entidade '%s' gerada em: %s\n", config.ModelName, modelTargetPath)

	// === 6. GERAR PACOTE DO REPOSITÓRIO (INTERFACE + REPOSITORY) ===
	repoPackagePath := filepath.Join("src", "apps", appName, "model", "repository", config.RepositoryPackageName)

	if err := generateRepositoryPackage(repoPackagePath, config); err != nil {
		return fmt.Errorf("falha ao gerar pacote de repositório: %v", err)
	}
	fmt.Printf("✅ Repositório '%s' gerado em: %s/\n", config.ModelName, repoPackagePath)

	return nil
}

// (inspectTable, mapPostgresTypeToGoType, snakeToCamel, titler permanecem iguais)
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
			return "json.RawMessage"
		}
		return "json.RawMessage"
	case "uuid":
		if isNullableBool {
			return "sql.NullString"
		}
		return "string"
	default:
		return "interface{}"
	}
}

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
	s = result.String()
	if s == "Id" {
		return "ID"
	}
	if strings.HasSuffix(s, "Id") {
		return strings.TrimSuffix(s, "Id") + "ID"
	}
	return s
}

var titler = cases.Title(language.Portuguese)

// --- FUNÇÕES DE GERAÇÃO DE ARQUIVO ---

// generateModelFile prepara o template e os imports para o arquivo de entidade
func generateModelFile(targetPath string, config StructConfig, imports map[string]bool) error {
	const modelTemplate = `package {{.PackageName}}

// <IMPORT_BLOCK> // Placeholder para importações dinâmicas

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
	// Adiciona imports dinâmicos
	actualImports := map[string]bool{
		"database/sql": true, // Sempre necessário para DBScanner
	}
	if imports["time"] {
		actualImports["time"] = true
	}
	if imports["encoding/json"] {
		actualImports["encoding/json"] = true
	}

	// (Adicionado do seu exemplo) Referência ao DBScanner do core
	coreEntitiesPath := strings.ReplaceAll(filepath.Join("deskapp/src/apps", "core", "model", "entities"), "\\", "/")
	actualImports[coreEntitiesPath] = true

	// Constrói o bloco de importação
	var importStr strings.Builder
	importStr.WriteString("import (\n")
	for imp := range actualImports {
		importStr.WriteString(fmt.Sprintf("\t\"%s\"\n", imp))
	}
	importStr.WriteString(")\n")

	// Substitui o placeholder
	finalTemplate := strings.Replace(modelTemplate, "// <IMPORT_BLOCK>", importStr.String(), 1)

	// Atualiza a assinatura do ScanRow para usar o DBScanner do core
	finalTemplate = strings.Replace(finalTemplate, "func (m *{{.ModelName}}) ScanRow(row DBScanner) error {", "func (m *{{.ModelName}}) ScanRow(row entities.DBScanner) error {", 1)
	finalTemplate = strings.Replace(finalTemplate, "type DBScanner interface {", "// DBScanner define a interface para Scan, implementada por *sql.Row e *sql.Rows.\ntype DBScanner interface {", 1)


	return generateFile(targetPath, finalTemplate, config)
}

// generateRepositoryPackage cria o diretório e os arquivos (interface.go, repository.go)
func generateRepositoryPackage(targetPath string, config StructConfig) error {
	// Garante que o diretório (ex: .../repository/usuario) exista
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório do pacote de repositório %s: %v", targetPath, err)
	}

	// --- 1. Gerar repository.go ---
	// Baseado em: repository.go
	const repositoryTemplate = `package {{.RepositoryPackageName}}

import (
	"database/sql"
	"deskapp/src/apps/core/model/repository"
	"{{.EntitiesPackagePath}}"
)

// {{.ModelName}}Repository é o repositório para a entidade {{.ModelName}}
type {{.ModelName}}Repository struct {
	*repository.BaseRepository[entities.{{.ModelName}}, *entities.{{.ModelName}}]
}

// New{{.ModelName}}Repository cria um novo {{.ModelName}}Repository
func New{{.ModelName}}Repository(db *sql.DB) *{{.ModelName}}Repository {
	base := repository.NewBaseRepository[entities.{{.ModelName}}](db, "{{.TableName}}", "{{.SchemaName}}")
	return &{{.ModelName}}Repository{
		BaseRepository: base,
	}
}
`
	repoFilePath := filepath.Join(targetPath, "repository.go")
	if err := generateFile(repoFilePath, repositoryTemplate, config); err != nil {
		return err
	}

	// --- 2. Gerar interface.go ---
	// Baseado em: interface.go
	const interfaceTemplate = `package {{.RepositoryPackageName}}

import (
	"context"
	"deskapp/src/apps/core/model/repository"
	"{{.EntitiesPackagePath}}"
)

type I{{.ModelName}}QueryBuilder = repository.IQueryBuilder[entities.{{.ModelName}}, *entities.{{.ModelName}}]

type I{{.ModelName}}Repository interface {
	// Where retorna o QueryBuilder específico.
	Where(ctx context.Context, queryFragment string, arg any) I{{.ModelName}}QueryBuilder
}
`
	ifaceFilePath := filepath.Join(targetPath, "interface.go")
	if err := generateFile(ifaceFilePath, interfaceTemplate, config); err != nil {
		return err
	}

	return nil
}

// generateFile é um helper refatorado para criar um arquivo a partir de um template
func generateFile(targetPath string, templateContent string, config StructConfig) error {
	tmpl, err := template.New(targetPath).Parse(templateContent)
	if err != nil {
		return fmt.Errorf("erro ao parsear template para %s: %v", targetPath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return fmt.Errorf("erro ao executar template para %s: %v", targetPath, err)
	}

	// Formata o código gerado
	formattedSource, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Printf("⚠️  Aviso: falha ao formatar o código para %s: %v\n", targetPath, err)
		formattedSource = buf.Bytes()
	}

	// Escrever o arquivo formatado
	if err := os.WriteFile(targetPath, formattedSource, 0644); err != nil {
		return fmt.Errorf("erro ao escrever arquivo %s: %v", targetPath, err)
	}

	return nil
}