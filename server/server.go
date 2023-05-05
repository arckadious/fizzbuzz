// This package creates and run a rest API server, using Gin framework
package server

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/arckadious/fizzbuzz/container"

	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"
)

const (
	// Timeout delay and graceful shutdown deadline
	Timeout     = time.Second * 180
	IdleTimeout = 60

	FizzBaseURI       = "/fizzbuzz"
	FizzStatisticsURI = "/statistics"

	URLPrefixVersion = "/v1"
	Scheme           = "http"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Server class
type Server struct {
	container *container.Container
}

// New constructor Server
func New(cntnr *container.Container) *Server {
	return &Server{
		container: cntnr,
	}
}

// Run launch http server
func (s *Server) Run() {

	host := "0.0.0.0:"

	//Windows OS : avoid firewall asking to allow network connections every execution time"
	// if s.container.Conf.Env == "localhost" {
	// 	host = "localhost:"
	// }

	server := &http.Server{
		Addr: host + strconv.Itoa(s.container.Conf.Port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: Timeout,
		ReadTimeout:  Timeout,
		IdleTimeout:  IdleTimeout,
		Handler:      s.handler(),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatal(err)
		}
	}()

	// Shutdown Database connections gracefully
	defer s.container.Db.Shutdown()

	// Process signals channel
	sigChannel := make(chan os.Signal, 1)

	// Graceful shutdown via SIGINT
	signal.Notify(sigChannel, os.Interrupt)

	logrus.Info("Server running...")
	<-sigChannel // Block until SIGINT received

	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	server.Shutdown(ctx)

	logrus.Warn("Server shutdown")
}

// handler configures gin and set all REST API endpoints
func (s *Server) handler() *gin.Engine {

	router := s.configureGin()

	//api doc (excluded from loggerMiddleware)
	router.Static("/swagger", "./swaggerui")

	loggerRouter := router.Group("")
	loggerRouter.Use(s.Logger())

	//ping
	loggerRouter.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "Ping OK !")
	})

	//api subrouter v1
	v1 := loggerRouter.Group(URLPrefixVersion)

	//fizzbuzz routes
	v1.POST(FizzBaseURI, func(c *gin.Context) {
		s.container.FizzAction.HandleFizz(c.Writer, c.Request)
	})
	v1.GET(FizzStatisticsURI, func(c *gin.Context) {
		s.container.FizzAction.HandleStatistics(c.Writer, c.Request)
	})

	return router

}

// configureGin set gin configuration.
func (s *Server) configureGin() *gin.Engine {

	router := gin.New()
	if s.container.Conf.Env == "localhost" {
		router.Use(gin.Logger())
	}

	router.HandleMethodNotAllowed = true

	router.Use(gin.CustomRecovery(s.recoveryHandler))
	router.NoRoute(gin.HandlerFunc(s.notFoundHandler))
	router.NoMethod(gin.HandlerFunc(s.methodNotAllowedHandler))

	return router
}
