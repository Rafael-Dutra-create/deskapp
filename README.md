# 🖥️ DeskApp

**DeskApp** é uma aplicação desktop/web modular desenvolvida em **Go (Golang)** com arquitetura em múltiplos aplicativos (apps), suporte a templates dinâmicos (multitemplate) e sistema de arquivos estáticos embutido.  
O projeto foi pensado para ser escalável, limpo e extensível, com integração facilitada entre módulos, templates e recursos estáticos.

---

## 🚀 Funcionalidades Principais

- 🔧 **Arquitetura modular** (apps independentes como `core`, `dash`, etc.)
- 🧩 **Sistema de templates dinâmicos** (com suporte a layouts base)
- 💾 **Conexão com PostgreSQL** (via pool de conexões)
- ⚙️ **Servidor HTTP com Gin** (rápido e fácil de estender)
- 📦 **Assets estáticos embutidos** (via `embed.FS`)
- 🧰 **Ferramentas internas** para criação de apps e mapeamento de tabelas
- 🧪 **Testes automatizados** para camadas de aplicação

---

## 🗂️ Estrutura do Projeto

```
deskapp/
├── Makefile                # Comandos de build, run e testes
├── go.mod / go.sum         # Dependências Go
├── src/
│   ├── main.go             # Ponto de entrada da aplicação
│   ├── app/                # Infraestrutura base para os apps
│   ├── apps/               # Módulos da aplicação (core, dash, etc.)
│   ├── internal/           # Pacotes internos (config, db, utils, etc.)
│   ├── static/             # Arquivos estáticos (CSS, JS, imagens)
│   └── templates/          # Templates HTML (com suporte a layouts)
```

---

## 🧩 Estrutura dos Apps

Cada app (ex: `core`, `dash`) segue uma estrutura semelhante:

```
src/apps/<nome_do_app>/
├── controller/             # Controladores (handlers de rota)
├── model/                  # Modelos de dados (opcional)
├── app.go                  # Registro do app e inicialização
├── routes.go               # Definição de rotas
```

---

## ⚙️ Requisitos

- Go 1.21+
- PostgreSQL 14+
- `make` (para facilitar execução dos comandos)

---

## 🧾 Configuração

A configuração de conexão com o banco de dados é definida no arquivo `src/internal/config/config.go`:

```go
postgresql://postgres:123456@localhost:5432/pydata?sslmode=disable
```

Você pode alterar o host, porta, usuário e senha conforme seu ambiente.

---

## ▶️ Como Executar

### 1. Clonar o repositório
```bash
git clone https://github.com/seuusuario/deskapp.git
cd deskapp
```

### 2. Rodar a aplicação
```bash
make run
```

O servidor será iniciado em:

```
http://localhost:8006
```

Você verá logs como:
```
✅ Sistema multitemplate configurado com sucesso!
✅ Sistema de arquivos estáticos configurado em /static
Servidor rodando em http://localhost:8006
```

### 3. Parar a aplicação
Pressione `CTRL + C`.

---

## 🧠 Estrutura de Templates

Os templates são carregados automaticamente pelo sistema **multitemplate**.  
Cada página é composta por um layout base (`base.html`) e um conteúdo específico, por exemplo:

```
templates/
├── base.html
├── index.html
├── about.html
└── dash/
    └── dash_index.html
```

Exemplo de herança de layout:

```html
{{ define "content" }}
  <h1>Sobre o DeskApp</h1>
  <p>Esta é a página About.</p>
{{ end }}
```

---

## 🧰 Ferramentas Internas

### Criar novo app

Há uma ferramenta CLI para gerar novos módulos automaticamente:

```bash
make createapp
```

Isso criará toda a estrutura básica do app (controllers, views, routes, etc.).

### Mapear tabelas do banco

```bash
make tablemap
```

---

## 🧪 Testes

Para rodar todos os testes:

```bash
make test
```

Gerar relatório de cobertura:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## 📦 Build

Para gerar o binário da aplicação:

```bash
make build
```

O binário será gerado em `./bin/deskapp`.

---


---

## 🪪 Licença

Este projeto é distribuído sob a licença **MIT**.  
Consulte o arquivo `LICENSE` para mais detalhes.

---

## 💬 Créditos

Desenvolvido por **Victor Gomes** 🧠  
💡 *"Código limpo, modular e elegante — como toda boa aplicação Go deve ser."*
