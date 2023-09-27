package ripple

import (
	. "github.com/bmbstack/ripple/helper"
	"github.com/flosch/pongo2/v4"
	"github.com/labstack/echo/v4"
	"io"
	"path/filepath"
	"sync"
)

//
type PongorOption struct {
	// Directory to load templates. Default is "templates"
	Directory string
	// Reload to reload templates every time.
	Reload bool
}

type Renderer struct {
	PongorOption
	templates map[string]*pongo2.Template
	lock      sync.RWMutex
}

func NewRenderer(config *BaseConfig, options ...PongorOption) *Renderer {
	var opt PongorOption
	if IsNotEmpty(options) {
		opt = options[0]
	}
	if len(opt.Directory) == 0 {
		opt.Directory = config.Templates
	}
	r := &Renderer{
		PongorOption: opt,
		templates:    make(map[string]*pongo2.Template),
	}
	return r
}

func getContext(templateData interface{}) pongo2.Context {
	if templateData == nil {
		return nil
	}
	contextData, isMap := templateData.(map[string]interface{})
	if isMap {
		return contextData
	}
	return nil
}

func (r *Renderer) buildTemplatesCache(name string) (t *pongo2.Template, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	t, err = pongo2.FromFile(filepath.Join(r.Directory, name))
	if err != nil {
		return
	}
	r.templates[name] = t
	return
}

func (r *Renderer) getTemplate(name string) (t *pongo2.Template, err error) {
	if r.Reload {
		return pongo2.FromFile(filepath.Join(r.Directory, name))
	}
	r.lock.RLock()
	var ok bool
	if t, ok = r.templates[name]; !ok {
		r.lock.RUnlock()
		t, err = r.buildTemplatesCache(name)
	} else {
		r.lock.RUnlock()
	}
	return
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	template, err := r.getTemplate(name)
	if err != nil {
		return err
	}
	err = template.ExecuteWriter(getContext(data), w)
	return err
}
