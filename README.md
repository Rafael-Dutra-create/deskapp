# 🖥️ DeskApp

**DeskApp** é uma aplicação desktop/web modular desenvolvida em **Go (Golang)** com arquitetura em múltiplos aplicativos (apps), suporte a templates dinâmicos (multitemplate) e sistema de arquivos estáticos embutido.  
O projeto foi pensado para ser escalável, limpo e extensível, com integração facilitada entre módulos, templates e recursos estáticos.

---

## 🚀 Funcionalidades Principais

- 🔧 **Arquitetura modular** (apps independentes como `core`, `dash`, etc.)
- 🧩 **Sistema de templates dinâmicos** (com suporte a layouts base e components reutilizáveis)
- 🧮 **Funções personalizadas em templates** (via `template.FuncMap`)
- 💾 **Conexão com PostgreSQL**
- ⚙️ **Servidor HTTP com Gin**
- 📦 **Assets embutidos** (`embed.FS`)
- 🧰 **Ferramentas internas CLI**
- 🧪 **Testes automatizados**

---

## 🗂️ Estrutura do Projeto

```
deskapp/
├── Makefile
├── go.mod / go.sum
├── src/
│   ├── main.go
│   ├── app/
│   ├── apps/
│   ├── internal/
│   ├── static/
│   └── templates/
│       ├── base.html
│       ├── components/
│       │   ├── chart.html
│       │   └── modulo.html
│       └── pages/
│           ├── index.html
│           └── about.html
```

---

## 🧩 Estrutura dos Apps

Cada app (ex: `core`, `dash`) segue a estrutura:

```
src/apps/<nome_do_app>/
├── controller/
├── model/
├── app.go
└── routes.go
```

---

## 🧠 Estrutura de Templates

O sistema **multitemplate** permite combinar **layouts**, **páginas** e **components reutilizáveis**.

### Exemplo de Página (`about.html`)

```html
{{ define "content" }}
  <h1>Sobre o DeskApp</h1>
  <p>Esta é a página About.</p>
  {{ template "chart" . }}
  {{ template "modulo" . }}
{{ end }}
```

---

## 🧱 Exemplos de Components

### 📊 **chart.html**

```html
{{ define "chart" }}
<div style="background-color: aliceblue;">
    <canvas id="myChart"></canvas>
</div>
{{ end }}
```

Uso:
```html
{{ template "chart" . }}
```

---

### 🧩 **modulo.html**

```html
{{ define "modulo" }}
<div>
    <p>{{.Modulo.Segmento}} - {{default .Modulo.Area "VAZIO"}} / {{.Modulo.Modulo}}</p>
</div>
{{ end }}
```

Uso:
```html
{{ template "modulo" . }}
```

---

## 🧮 Adicionando Funções ao Template

O DeskApp permite registrar **funções personalizadas** para uso direto nos templates HTML.  
Essas funções são mapeadas em um `template.FuncMap`, definido no pacote `functemplates`.

### 📁 Arquivo: `src/internal/functemplates/register.go`

```go
package functemplates

import "html/template"

var funcMap template.FuncMap

func register(name string, fn any) {
	funcMap[name] = fn
}

func init() {
	funcMap = make(template.FuncMap)
	register("default", defaultFunc)
}

func GetFuncMap() template.FuncMap {
	return funcMap
}
```

### 🧩 Exemplo de função (`defaultFunc`)

```go
func defaultFunc(value interface{}, fallback string) string {
	if value == nil || value == "" {
		return fallback
	}
	return fmt.Sprintf("%v", value)
}
```

Essa função é usada no template `modulo.html`:

```html
{{ default .Modulo.Area "VAZIO" }}
```

---

### ➕ Como adicionar novas funções

1. **Defina a função** no mesmo pacote (`functemplates`):  
   ```go
   func upperCase(s string) string {
       return strings.ToUpper(s)
   }
   ```

2. **Registre a função** dentro do `init()`:
   ```go
   func init() {
       funcMap = make(template.FuncMap)
       register("default", defaultFunc)
       register("upper", upperCase)
   }
   ```

3. **Use no template**:
   ```html
   <p>{{ upper "deskapp" }}</p>
   <!-- saída: DESKAPP -->
   ```

Dessa forma, qualquer função registrada pode ser usada diretamente nos templates, tornando-os muito mais expressivos e reutilizáveis.

---

## 🧰 Ferramentas Internas

- **Criar app:** `make app`  
- **Mapear tabelas:** `make tablemap`  
- **Gerar DTO:** `make dto`

---

## 🧪 Testes

```bash
make test
```

---

## 📦 Build

```bash
make build
```

Binário gerado em `./bin/deskapp`.

---

## 🪪 Licença

Licença **MIT** — veja `LICENSE`.

---

## 💬 Créditos

Desenvolvido por **Victor Gomes** 🧠  
💡 *"Código limpo, modular e elegante — como toda boa aplicação Go deve ser."*
