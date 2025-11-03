package view
import (
	"deskapp/src/internal/utils"
	"html/template"
	"net/http"
	"sync"
)

type View struct {
	mode utils.MODE
}

func NewView( mode utils.MODE)*View {
	return &View{
		mode: mode,
	}
}

type TemplateCache struct {
	templates map[string]*template.Template
	mu        sync.RWMutex
}

var cache = &TemplateCache{templates: make(map[string]*template.Template)}

// LoadTemplate carrega e cacheia um template, com suporte a layout base.
func (v *View) LoadTemplate(name string) (*template.Template, error) {
	if v.mode == utils.RELEASE {
		cache.mu.RLock()
		t, ok := cache.templates[name]
		cache.mu.RUnlock()
		if ok {
			return t, nil
		}
	}

	cache.mu.Lock()
	defer cache.mu.Unlock()

	// Recarregar caso não esteja em cache
	baseTemplate := "src/templates/base.html"
	pageTemplate := "src/templates/" + name

	// Carrega o template específico junto com o base
	t, err := template.ParseFiles(baseTemplate, pageTemplate)
	if err != nil {
		return nil, err
	}
	cache.templates[name] = t
	return t, nil
}

// Render executa o template nomeado dentro do layout base.
func (v *View) Render(w http.ResponseWriter, name string, data any) {
	t, err := v.LoadTemplate(name)
	if err != nil {
		http.Error(w, "Erro carregando templates: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, "Erro renderizando: "+err.Error(), http.StatusInternalServerError)
	}
}