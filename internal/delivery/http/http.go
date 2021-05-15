package http

import (
	"bank-system-go/internal/config"
	"bank-system-go/pkg/logger"
	"net/http"

	ginLogger "github.com/gin-contrib/logger"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	engine *gin.Engine
	logger logger.Logger
}

func NewHttpServer(config config.Config, logger logger.Logger) *HttpServer {
	if config.Release {
		gin.SetMode(gin.ReleaseMode)
	}
	httpServer := &HttpServer{
		engine: gin.Default(),
		logger: logger,
	}
	httpServer.setRouter()

	return httpServer
}

func (server *HttpServer) setRouter() {
	server.engine.Use(gin.Recovery())
	server.engine.Use(gzip.Gzip(gzip.DefaultCompression))
	server.engine.Use(ginLogger.SetLogger(ginLogger.Config{
		Logger: &server.logger.Logger,
		UTC:    true,
	}))
	server.engine.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Page Not Found")
	})
	server.engine.NoMethod(func(c *gin.Context) {
		c.String(http.StatusMethodNotAllowed, "Method Not Allowed")
	})
}

func (server *HttpServer) Run(addr ...string) error {
	return server.engine.Run(addr...)
}
