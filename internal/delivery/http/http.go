package http

import (
	"bank-system-go/internal/config"
	"bank-system-go/internal/controller"
	"bank-system-go/internal/model"
	"bank-system-go/pkg/logger"
	"context"
	"net/http"
	"strings"

	ginLogger "github.com/gin-contrib/logger"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

const (
	_requestIDHeaderName      = "X-Request-Id"
	_authenticationHeaderName = "Authentication"
)

const (
	_userContextKey = "user_from_context"
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
	httpServer := &HttpServer{
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
	server.engine.Use(server.RequestIDMiddleware)
	server.engine.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Page Not Found")
	})
	server.engine.NoMethod(func(c *gin.Context) {
		c.String(http.StatusMethodNotAllowed, "Method Not Allowed")
	})

	{
		apiV1 := server.engine.Group("/api/v1")
		apiV1.POST("/register", server.Register)
		apiV1.POST("/login", server.Login)
		apiV1.Use(server.AuthMiddleware)
		apiV1.POST("/wallet", server.CreateWallet)
		apiV1.GET("/wallets", server.ListWallet)
		apiV1.POST("/wallet/balance", server.UpdateWalletBalance)
	}
}

func (server *HttpServer) RequestIDMiddleware(ctx *gin.Context) {
	uuid := uuid.New().String()
	if requestID := ctx.GetHeader(_requestIDHeaderName); len(requestID) > 0 {
		uuid = requestID
	}
	ctx.Set(_requestIDHeaderName, uuid)
	ctx.Header(_requestIDHeaderName, uuid)

	ctx.Next()
}

func (server *HttpServer) Run(addr ...string) error {
	errg := errgroup.Group{}
	errg.Go(func() error {
		return server.controller.GatewayCallback(context.Background())
	})
	errg.Go(func() error {
		return server.engine.Run(addr...)
	})
	return errg.Wait()
}

func (server *HttpServer) RequestID(c *gin.Context) string {
	v, ok := c.Get(_requestIDHeaderName)
	if !ok {
		return uuid.New().String()
	}
	requestID, ok := v.(string)
	if !ok {
		return uuid.New().String()
	}
	return requestID
}

func (server *HttpServer) Register(c *gin.Context) {
	req := model.RegisterUserRequest{}

	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	code, resp, err := server.controller.RegisterUser(c, server.RequestID(c), req)
	if err != nil {
		c.AbortWithError(code, err)
		return
	}
	c.JSON(code, resp)
}

func (server *HttpServer) Login(c *gin.Context) {
	req := model.UserLoginRequest{}

	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	req.IP = c.ClientIP()
	code, resp, err := server.controller.Login(c, server.RequestID(c), req)
	if err != nil {
		c.AbortWithError(code, err)
		return
	}
	c.JSON(code, resp)
}

func (server *HttpServer) GetUser(c *gin.Context) model.User {
	v, ok := c.Get(_userContextKey)
	if !ok {
		return model.User{}
	}
	user, ok := v.(model.User)
	if !ok {
		return model.User{}
	}
	return user
}

func (server *HttpServer) AuthMiddleware(c *gin.Context) {
	req := model.VerifyUserRequest{
		Token: strings.TrimPrefix(c.GetHeader(_authenticationHeaderName), "Bearer "),
	}

	code, resp, err := server.controller.VerifyUser(c, server.RequestID(c), req)
	if err != nil {
		c.AbortWithError(code, err)
		return
	}

	c.Set(_userContextKey, resp.User)
	c.Next()
}

func (server *HttpServer) CreateWallet(c *gin.Context) {
	req := model.CreateWalletRequest{}

	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	code, resp, err := server.controller.CreateWallet(c, server.RequestID(c), server.GetUser(c), req)
	if err != nil {
		c.AbortWithError(code, err)
		return
	}
	c.JSON(code, resp)
}

func (server *HttpServer) ListWallet(c *gin.Context) {
	req := model.ListWalletRequest{}

	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	code, resp, err := server.controller.ListWallet(c, server.RequestID(c), server.GetUser(c), req)
	if err != nil {
		c.AbortWithError(code, err)
		return
	}
	c.JSON(code, resp)
}

func (server *HttpServer) UpdateWalletBalance(c *gin.Context) {
	req := model.UpdateWalletBalanceRequest{}

	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	code, resp, err := server.controller.UpdateWalletBalance(c, server.RequestID(c), server.GetUser(c), req)
	if err != nil {
		c.AbortWithError(code, err)
		return
	}
	c.JSON(code, resp)
}
