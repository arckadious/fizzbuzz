// This package creates and run a rest API server, using Gin framework
package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"time"

	cst "github.com/arckadious/fizzbuzz/constant"

	"github.com/arckadious/fizzbuzz/container"
	"github.com/arckadious/fizzbuzz/model"
	"github.com/arckadious/fizzbuzz/response"
	"github.com/arckadious/fizzbuzz/util"

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

	//Middlewares for Gin
	middlewares := []gin.HandlerFunc{}

	//Set Gin mode
	if s.container.Conf.Env != "localhost" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
		middlewares = append(middlewares, gin.Logger()) //debug logger local development
		var f *os.File
		f, err := os.OpenFile("gin.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			logrus.Fatal(err)
		}
		gin.DefaultWriter = f
	}

	router := gin.New()
	router.HandleMethodNotAllowed = true

	// CustomRecovery handles unexpected 'panic' call during process, and return a custom 500 internal server error.
	middlewares = append(middlewares, gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		logrus.Error(err)
		response.New(
			http.StatusInternalServerError,
			cst.StatusError,
			[]response.ApiError{
				{
					Code:    cst.ErrorInternalServerError,
					Message: "Internal Server Error, oups !",
				},
			}, nil).WriteJSONResponse(c.Writer)
	}))

	// NoRoute handles API server response when endpoint couldn't be found
	router.NoRoute(gin.HandlerFunc(func(c *gin.Context) {
		response.New(
			http.StatusNotFound,
			cst.StatusError,
			[]response.ApiError{
				{
					Code:    cst.ErrorRouteNotFound,
					Message: "Route Endpoint Not Found.",
				},
			}, nil).WriteJSONResponse(c.Writer)
	}))

	// NoMethod handles API server response when endpoint method is not allowed
	router.NoMethod(gin.HandlerFunc(func(c *gin.Context) {
		response.New(
			http.StatusMethodNotAllowed,
			cst.StatusError,
			[]response.ApiError{
				{
					Code:    cst.ErrorMethodNotAllowed,
					Message: "Method is not allowed, boy.",
				},
			}, nil).WriteJSONResponse(c.Writer)
	}))

	router.Use(middlewares...)

	return router
}

// Logger send requests and response to database, and generate checksum if needed
func (s *Server) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Generate unique ID to make link between request and its associated response (stored in a different table)
		corID, _ := util.GenerateUID()

		body, err := util.ExtractBody(c.Request)
		if err != nil {
			logrus.Error("Logger coudn't send request data: ", err)
			return
		}

		//Create a checksum for the current request, only if it's the main endpoint (/fizzbuzz) and data is valid.
		var data model.Input
		checksum := ""
		if c.Request.RequestURI == URLPrefixVersion+FizzBaseURI && json.Unmarshal(body, &data) == nil && s.container.Validator.Struct(data) == nil {
			checksum = util.GetMD5Hash(data.String())
		}

		// Create copy to be used inside the goroutine - See Gin documentation : https://gin-gonic.com/docs/examples/goroutines-inside-a-middleware/
		cCp := c.Copy()

		go func() {
			if err := s.container.Repo.LogToDB("request", string(body), Scheme+"://"+path.Join(cCp.Request.Host, cCp.Request.RequestURI), corID, checksum, ""); err != nil {
				logrus.Error(err)
			}
		}()

		// Intercept Writer in order to get response body
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		status := blw.Status()
		respBody := blw.body.String()
		go func() {
			if err := s.container.Repo.LogToDB("response", respBody, "", corID, checksum, strconv.Itoa(status)); err != nil {
				logrus.Error(err)
			}
		}()
	}
}
