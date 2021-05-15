package http

import (
	"bank-system-go/internal/config"
	"bank-system-go/internal/controller"
	"bank-system-go/internal/model"
	"bank-system-go/pkg/logger"
	"context"
	"net/http"
	"os"

	ginLogger "github.com/gin-contrib/logger"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type HttpServer struct {
	hostname   string
	engine     *gin.Engine
	controller *controller.GatewayController
	logger     logger.Logger
}

func NewHttpServer(config config.Config, logger logger.Logger, controller *controller.GatewayController) (*HttpServer, error) {
	if config.Release {
		gin.SetMode(gin.ReleaseMode)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	httpServer := &HttpServer{
		hostname:   "gateway_" + hostname,
		engine:     gin.Default(),
		controller: controller,
		logger:     logger,
	}

	httpServer.setRouter()

	return httpServer, nil
}

func (server *HttpServer) setRouter() {
	server.engine.Use(gin.Recovery())
	server.engine.Use(gzip.Gzip(gzip.DefaultCompression))
	server.engine.Use(ginLogger.SetLogger(ginLogger.Config{
		Logger: &server.logger.Logger,
		UTC:    true,
	}))
	server.engine.Use(server.requestIDMiddleware)
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

func (server *HttpServer) requestIDMiddleware(ctx *gin.Context) {
	uuid := uuid.New().String()
	if requestID := ctx.GetHeader("X-Request-Id"); len(requestID) > 0 {
		uuid = requestID
	}
	ctx.Set("X-Request-Id", uuid)
	ctx.Header("X-Request-Id", uuid)

	ctx.Next()
}

func (server *HttpServer) Run(addr ...string) error {
	errg := errgroup.Group{}
	errg.Go(func() error {
		return server.controller.GatewayCallback(context.Background(), server.hostname)
	})
	errg.Go(func() error {
		return server.engine.Run(addr...)
	})
	return errg.Wait()
}

func (server *HttpServer) Register(c *gin.Context) {
	requestID, _ := c.Get("X-Request-Id")
	req := RegisterUserRequest{}

	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	code, user, err := server.controller.RegisterUser(c, requestID.(string), server.hostname, model.User{
		Name:     req.Name,
		Account:  req.Account,
		Password: []byte(req.Password),
	})
	if err != nil {
		c.AbortWithError(code, err)
		return
	}
	c.JSON(code, user)
}
