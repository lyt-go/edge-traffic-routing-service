// Package app 负责依赖装配。
package app

import (
	"net/http"

	"reverseproxy/internal/config"
	"reverseproxy/internal/handler"
	"reverseproxy/internal/service"
	"reverseproxy/internal/store"
	"reverseproxy/pkg/logger"
)

type App struct {
	server *handler.Server
}

func New(cfg *config.Config, log *logger.Logger) (*App, error) {
	st := store.NewMemoryStore()
	svc := service.New(st, log, cfg)
	server := handler.NewServer(svc, log, cfg)
	log.Infof("应用装配完成，配置：%s", cfg.String())
	return &App{server: server}, nil
}

func (a *App) Routes() http.Handler { return a.server.Routes() }
