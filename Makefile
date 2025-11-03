# Makefile para DeskApp
.PHONY: run create-app help

# Variáveis
APP_NAME=deskapp
SRC_DIR=./src
SCRIPTS_DIR=./src/internal/scripts
CMD_DIR=./src/internal/cmd

# Comando padrão
.DEFAULT_GOAL := help

# Executar a aplicação
run:
	@echo "🚀 Iniciando $(APP_NAME)..."
	go run $(SRC_DIR)

# Criar um novo app
create-app:
	@echo "📱 Criando novo app..."
	go run $(CMD_DIR)/createapp

# Comando alternativo para criar app (usando scripts diretamente)
create-app-alt:
	@echo "📱 Criando novo app (alternativo)..."
	go run $(SCRIPTS_DIR)/create_app.go

# Instalar dependências
deps:
	@echo "📦 Verificando dependências..."
	go mod tidy
	go mod download

# Build da aplicação
build:
	@echo "🔨 Buildando $(APP_NAME)..."
	go build -o bin/$(APP_NAME) $(SRC_DIR)

# Limpar binários
clean:
	@echo "🧹 Limpando binários..."
	rm -rf bin/

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
	go test ./...

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
	@echo "  run           - Executar a aplicação (go run ./src)"
	@echo "  create-app    - Criar um novo app (go run ./src/internal/cmd/createapp)"
	@echo "  create-app-alt- Criar app alternativo (go run ./src/internal/scripts/create_app.go)"
	@echo "  deps          - Instalar/atualizar dependências"
	@echo "  build         - Build da aplicação"
	@echo "  clean         - Limpar binários"
	@echo "  dev           - Desenvolvimento com auto-reload (air)"
	@echo "  test          - Executar testes"
	@echo "  fmt           - Verificar formatação do código"
	@echo "  lint          - Executar linter"
	@echo "  help          - Mostrar esta ajuda"
	@echo ""
	@echo "Exemplos:"
	@echo "  make run"
	@echo "  make create-app"
	@echo "  make dev"