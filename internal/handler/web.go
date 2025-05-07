package handler

import (
	"exaroton-wa-bot/internal/config"
	"exaroton-wa-bot/internal/dto"
	"exaroton-wa-bot/internal/middleware"
	"exaroton-wa-bot/internal/service"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Web struct {
	Router  *echo.Echo
	cfg     *config.Cfg
	session dto.WebSession

	middleware *middleware.Middleware
	svc        *service.Service
}

func NewWeb(cfg *config.Cfg, svc *service.Service) *Web {
	router := config.NewEcho()
	router.HTTPErrorHandler = errorHandler()

	// mount pages
	pagesDir, err := filepath.Abs(cfg.MustString(config.KeyPagesDir))
	if err != nil {
		panic(fmt.Sprintf("failed to get pages dir (key: %s)", config.KeyPagesDir))
	}

	hotReload := false
	if cfg.Args.Env == config.EnvDevelopment {
		hotReload = true
	}

	router.Renderer = NewRenderer(pagesDir, hotReload)

	webSession := dto.NewWebSession()

	router.Validator = config.NewValidator()
	middleware := middleware.NewMiddleware(cfg, svc.AuthService, webSession)

	web := &Web{
		Router:     router,
		cfg:        cfg,
		middleware: middleware,
		svc:        svc,
		session:    webSession,
	}

	web.LoadRoutes()

	return web
}

func (h *Web) RunHTTP(port int) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: h.Router,
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
