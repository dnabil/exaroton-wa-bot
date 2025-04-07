package handler

import (
	"exaroton-wa-bot/internal/config"
	"exaroton-wa-bot/internal/service"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type Web struct {
	Router *echo.Echo

	svc *service.Service
	cfg *config.Cfg
}

func NewWeb(cfg *config.Cfg, svc *service.Service) *Web {
	rest := &Web{
		svc:    svc,
		Router: config.NewEcho(),
		cfg:    cfg,
	}

	return rest
}

func (w *Web) Run() error {
	pagesDir, err := filepath.Abs(w.cfg.MustString(config.KeyPagesDir))
	if err != nil {
		panic(fmt.Sprintf("failed to get pages dir (key: %s)", config.KeyPagesDir))
	}

	hotReload := false
	if w.cfg.Args.Env == config.EnvDevelopment {
		hotReload = true
	}

	w.Router.Renderer = NewRenderer(pagesDir, hotReload)
	w.Router.HTTPErrorHandler = webErrorHandler()
	w.routes()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", w.cfg.String(config.KeyPort)),
		Handler: w.Router,
	}

	return srv.ListenAndServe()
}

// ===============================================================================
// HTML Renderer

type Renderer struct {
	template  *template.Template
	location  string
	hotReload bool
}

func (t *Renderer) LoadTemplates() {
	t.template = template.New("")

	tmp, err := t.template.ParseGlob(filepath.Join(t.location, "*.tmpl"))
	if err != nil {
		panic(err)
	}

	if tmp != nil {
		t.template = tmp
	}
}

func (t *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if t.hotReload {
		t.LoadTemplates()
	}

	return t.template.ExecuteTemplate(w, name, data)
}

func NewRenderer(location string, hotReload bool) echo.Renderer {
	renderer := &Renderer{
		location:  location,
		hotReload: hotReload,
	}

	renderer.LoadTemplates()

	return renderer
}

// ==============================================================================
