package http

import (
	"bank-system-go/internal/config"
	"bank-system-go/pkg/logger"
	"bank-system-go/pkg/mq"
	"net/http"

	ginLogger "github.com/gin-contrib/logger"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	engine *gin.Engine
	mq     mq.MQ
	logger logger.Logger
}

func NewHttpServer(config config.Config, logger logger.Logger, mq mq.MQ) *HttpServer {
	if config.Release {
		gin.SetMode(gin.ReleaseMode)
	}
	httpServer := &HttpServer{
		engine: gin.Default(),
		mq:     mq,
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

	{
		apiV1 := server.engine.Group("/api/v1")
		apiV1.POST("/register", server.Register)
	}
}

func (server *HttpServer) Run(addr ...string) error {
	return server.engine.Run(addr...)
}

func (server *HttpServer) Register(c *gin.Context) {
	c.String(http.StatusOK, "Hi")
}
