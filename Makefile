# Makefile para DeskApp
.PHONY: run create-app help

# Variáveis
APP_NAME=deskapp
SRC_DIR=./src
SCRIPTS_DIR=./src/internal/scripts
DATABASE_URL=postgresql://postgres:123456@localhost:5432/pydata?sslmode=disable
CMD_DIR=./src/internal/cmd

# Comando padrão
.DEFAULT_GOAL := help

# Executar a aplicação
run:
	@echo "🚀 Iniciando $(APP_NAME)..."
	go run $(SRC_DIR)


# Instalar dependências
deps:
	@echo "📦 Verificando dependências..."
	go mod tidy
	go mod download

# Build da aplicação
build:
	@echo "🔨 Buildando $(APP_NAME)..."
	go build -o bin/$(APP_NAME) $(SRC_DIR)

# NOVO - Build para Windows (amd64)
build-windows:
	@echo "🔨 Buildando $(APP_NAME) para Windows (amd64)..."
	GOOS=windows GOARCH=amd64 go build -o bin/$(APP_NAME).exe $(SRC_DIR)

# NOVO - Build para Linux (amd64)
build-linux:
	@echo "🔨 Buildando $(APP_NAME) para Linux (amd64)..."
	GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME) $(SRC_DIR)

# NOVO - Build para macOS (amd64)
build-macos:
	@echo "🔨 Buildando $(APP_NAME) para macOS (amd64)..."
	GOOS=darwin GOARCH=amd64 go build -o bin/$(APP_NAME) $(SRC_DIR)

# Criar um novo app (executando o script alternativo, se existir)
app:
	@echo "📱 Criando novo app (alternativo)..."
	go run $(SCRIPTS_DIR) create-app
    
# Mapear tabela (executando o script)
tablemap:
	@echo "🗺️ Mapeando tabela para struct..."
	go run $(SCRIPTS_DIR) tablemap

dto:
	@echo "🗺️ Mapeando tabela para struct..."
	go run $(SCRIPTS_DIR) create-dto

migrate-up:
	go run $(SCRIPTS_DIR) migrate up

migrate-down:
	go run $(SCRIPTS_DIR) migrate down

migrate-status:
	go run $(SCRIPTS_DIR) migrate status
	


# Limpar binários
clean:
	@echo "🧹 Limpando binários..."
	rm -rf bin/
	rm *.out
	rm coverage.html
	go clean --cache

coverprofile:
	go test $$(go list -f '{{if .TestGoFiles}}{{.ImportPath}}{{end}}' ./...) --cover -coverprofile cover.out
	go tool cover -html=cover.out -o=coverage.html


# Desenvolvimento com auto-reload (se tiver air instalado)
dev:
	@if command -v air > /dev/null; then \
		echo "🔥 Iniciando desenvolvimento com auto-reload..."; \
		air; \
	else \
		echo "❌ Air não instalado. Instale com: go install github.com/cosmtrek/air@latest"; \
		echo "💡 Ou executando: make run"; \
	fi

# Testes
test:
	@echo "🧪 Executando testes..."
	go test $$(go list -f '{{if .TestGoFiles}}{{.ImportPath}}{{end}}' ./...)
# Verificar formatação
fmt:
	@echo "🎨 Verificando formatação..."
	go fmt ./...

# Lint
lint:
	@echo "🔍 Executando lint..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "❌ golangci-lint não instalado. Instale com: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Help
help:
	@echo "DeskApp - Makefile Commands"
	@echo ""
	@echo "Comandos disponíveis:"
	@echo "  run             "
	@echo "  createapp       "
	@echo "  tablemap        "
	@echo "  migrate-up      "
	@echo "  migrate-down    "
	@echo "  migrate-status  "
	@echo "  deps          "
	@echo "  build         "
	@echo "  clean         "
	@echo "  dev           "
	@echo "  test          "
	@echo "  fmt           "
	@echo "  lint          "
	@echo "  help          "
	@echo ""
	@echo "Exemplos:"
	@echo "  make run"
	@echo "  make create-app"
	@echo "  make dev"