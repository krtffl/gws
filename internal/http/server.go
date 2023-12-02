package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/krtffl/gws/internal/cookie"
	"github.com/krtffl/gws/internal/http/middlewares"
	"github.com/krtffl/gws/internal/http/webui"
	"github.com/krtffl/gws/internal/logger"
)

type Server struct {
	ctx        context.Context
	shutdownFn context.CancelFunc
	port       uint
	handler    *webui.Handler
	cookie     *cookie.Service
}

func New(
	port uint,
	handler *webui.Handler,
	cookie *cookie.Service,
) *Server {
	ctx, shutdownFn := context.WithCancel(context.Background())
	return &Server{
		ctx:        ctx,
		shutdownFn: shutdownFn,
		port:       port,
		handler:    handler,
		cookie:     cookie,
	}
}

func (srv *Server) Run() error {
	logger.Info("[HTTP] - Starting to listen on port %d", srv.port)
	r := chi.NewMux()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/healthcheck", func(r chi.Router) {
		r.Get("/", handleHealthcheck)
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", srv.handler.Index)
		r.Get("/form", srv.handler.Form)
		r.Post("/form", srv.handler.Upload)
		r.Get("/challenge", srv.handler.Challenge)
		r.Post("/challenge", srv.handler.Solve)
	})

	r.Route("/secure", func(r chi.Router) {
		r.Use(middlewares.Solved(srv.cookie))
		r.Get("/memories", srv.handler.Memories)
	})

	fs := http.FileServer(http.Dir("public"))
	r.Handle("/public/*", http.StripPrefix("/public/", fs))

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", srv.port),
		Handler: r,
	}

	go func() {
		<-srv.ctx.Done()
		if err := httpServer.Shutdown(srv.ctx); err != nil {
			logger.Error("[HTTP Server] - Failed to shutdown on port %d. %v", srv.port, err)
		}
		logger.Info("[HTTP Server] - Server on port %d has shutdown", srv.port)
	}()

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("[HTTP Server] - %v", err)
	}

	return nil
}

func (srv *Server) Shutdown() {
	logger.Info("[HTTP Server] - Server on port %d shutting down", srv.port)
	srv.shutdownFn()
	logger.Info("[HTTP Server] - Server on port %d shutted down", srv.port)
}
