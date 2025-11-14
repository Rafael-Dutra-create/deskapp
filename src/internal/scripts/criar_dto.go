package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)


// CreateDTOScript é o seu script para criar DTOs
type CreateDTOScript struct {
	ScriptBase
}

func (s *CreateDTOScript) Name() string {
	return "create-dto"
}

func (s *CreateDTOScript) Description() string {
	return "Cria uma nova DTO com tags json e form"
}

// prompt é uma função utilitária para ler a entrada do console
func prompt(reader *bufio.Reader, text string, required bool) (string, error) {
	for {
		fmt.Print(text)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		input = strings.TrimSpace(input)

		if input == "" && required {
			fmt.Println("Este campo é obrigatório. Tente novamente.")
			continue
		}
		return input, nil
	}
}

// Regex para conversão para snake_case
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// toSnakeCase converte uma string camelCase para snake_case
func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// toGoName converte um nome de campo (ex: 'firstName') para um nome de struct Go (ex: 'FirstName')
func toGoName(str string) string {
	if str == "" {
		return ""
	}
	// Converte a primeira letra para maiúscula
	return strings.ToUpper(string(str[0])) + str[1:]
}

// Execute é onde a mágica acontece
func (s *CreateDTOScript) Execute(args []string) error {
	reader := bufio.NewReader(os.Stdin)

	// 1. Perguntar o nome do App (para o pacote e diretório)
	appName, err := prompt(reader, "Nome do App (ex: users, products): ", true)
	if err != nil {
		return fmt.Errorf("falha ao ler nome do app: %w", err)
	}
	packageName := strings.ToLower(appName)

	// 2. Perguntar o nome do DTO
	dtoName, err := prompt(reader, "Nome do DTO (ex: CreateUser): ", true)
	if err != nil {
		return fmt.Errorf("falha ao ler nome do DTO: %w", err)
	}
	// Precisamos dos nomes Go e snake_case
	goDtoName := toGoName(dtoName)
	snakeDtoName := toSnakeCase(dtoName)

	// 3. Perguntar os campos em loop
	type Field struct {
		Name     string
		Type     string
		Required bool
	}
	var fields []Field

	fmt.Println("\nDigite os campos (ex: 'firstName' e 'string').")
	fmt.Println("Pressione [Enter] no 'Nome do Campo' para finalizar.")

	for {
		fieldName, err := prompt(reader, "  Nome do Campo (ex: firstName): ", false)
		if err != nil {
			return fmt.Errorf("falha ao ler nome do campo: %w", err)
		}

		// Se o nome do campo for vazio, o usuário terminou
		if fieldName == "" {
			break
		}

		fieldType, err := prompt(reader, "    Tipo do Campo (ex: string, int, *bool): ", true)
		if err != nil {
			return fmt.Errorf("falha ao ler tipo do campo: %w", err)
		}

		isRequired, err := prompt(reader, "    Obrigatório (1,0): ", true)
		if err != nil {
			return fmt.Errorf("falha ao ler isRequired: %w", err)
		}

		fields = append(fields, Field{Name: fieldName, Type: fieldType, Required: isRequired == "1"})
	}

	if len(fields) == 0 {
		fmt.Println("Nenhum campo adicionado. Saindo.")
		return nil
	}

	// 4. Gerar APENAS as linhas de string dos novos campos
	var newFieldsBuilder strings.Builder
	for _, field := range fields {
		goName := toGoName(field.Name)
		snakeName := toSnakeCase(field.Name)

		// Formata a tag
		tag := fmt.Sprintf("`json:\"%s\" form:\"%s\"", snakeName, snakeName)
		if field.Required {
			tag += ` binding:"required"`
		}
		tag += "`"
		// Escreve a linha do campo
		newFieldsBuilder.WriteString(fmt.Sprintf("\t%s %s %s\n", goName, field.Type, tag))
	}
	newFieldsString := newFieldsBuilder.String()

	// 5. Escrever o arquivo (Lógica de Criar OU Atualizar)
	dirPath := fmt.Sprintf("src/apps/%s/model/dtos", packageName)
	fileName := fmt.Sprintf("%s.go", snakeDtoName) // Nome do arquivo é snake_case
	filePath := fmt.Sprintf("%s/%s", dirPath, fileName)

	// Criar o diretório se não existir (necessário em ambos os casos)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return fmt.Errorf("falha ao criar diretório '%s': %w", dirPath, err)
	}

	// Tentar ler o arquivo
	fileContent, err := os.ReadFile(filePath)

	if os.IsNotExist(err) {
		// --- ARQUIVO NÃO EXISTE: Criar novo ---
		fmt.Println("Arquivo não encontrado. Criando novo...")
		var sb strings.Builder

		// Define o pacote (usando o nome do app)
		sb.WriteString("package dto\n")
		// Define a struct (usando o nome Go + "DTO")
		sb.WriteString(fmt.Sprintf("type %sDTO struct {\n", goDtoName))
		sb.WriteString(newFieldsString) // Adiciona os campos gerados
		sb.WriteString("}\n")

		// Escrever o novo arquivo
		err = os.WriteFile(filePath, []byte(sb.String()), 0644)
		if err != nil {
			return fmt.Errorf("falha ao escrever novo arquivo '%s': %w", filePath, err)
		}
		fmt.Printf("\n✅ Novo DTO criado com sucesso em: %s\n", filePath)

	} else if err == nil {
		// --- ARQUIVO EXISTE: Atualizar ---
		fmt.Println("Arquivo encontrado. Adicionando campos...")
		contentStr := string(fileContent)

		// +++ LÓGICA DE VERIFICAÇÃO DE DUPLICIDADE +++
		var duplicateFields []string
		for _, field := range fields {
			goName := toGoName(field.Name)
			// A verificação mais segura é procurar pelo GoName (ex: "FirstName")
			// seguido por um espaço, e precedido por um tab.
			// Ex: "\tFirstName "
			searchString := "\t" + goName + " "

			if strings.Contains(contentStr, searchString) {
				duplicateFields = append(duplicateFields, field.Name)
			}
		}

		if len(duplicateFields) > 0 {
			// Encontrou duplicados. Retorna um erro.
			return fmt.Errorf("falha ao atualizar DTO: os seguintes campos já existem: %s",
				strings.Join(duplicateFields, ", "))
		}


		// Encontrar a última '}' no arquivo.
		// Esta é uma abordagem simples que assume que a '}' da struct é o último caractere '}'
		lastBraceIndex := strings.LastIndex(contentStr, "}")
		if lastBraceIndex == -1 {
			return fmt.Errorf("formato de arquivo inválido. Não foi possível encontrar a '}' de fechamento da struct em '%s'", filePath)
		}

		// Inserir os novos campos ANTES da última '}'
		var sb strings.Builder
		sb.WriteString(contentStr[:lastBraceIndex]) // Tudo antes da '}'
		sb.WriteString(newFieldsString)             // Os novos campos
		sb.WriteString(contentStr[lastBraceIndex:]) // A '}' e o resto do arquivo (se houver)

		// Escrever o arquivo modificado
		err = os.WriteFile(filePath, []byte(sb.String()), 0644)
		if err != nil {
			return fmt.Errorf("falha ao atualizar arquivo '%s': %w", filePath, err)
		}
		fmt.Printf("\n✅ DTO atualizado com sucesso em: %s\n", filePath)

	} else {
		// Outro erro ao ler o arquivo
		return fmt.Errorf("falha ao verificar arquivo '%s': %w", filePath, err)
	}

	return nil
}
