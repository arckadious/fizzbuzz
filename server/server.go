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
		Handler:      s.Handler(),
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

// Handler configures all endpoints and middlewares with Gin
func (s *Server) Handler() *gin.Engine {

	router := gin.New()
	router.HandleMethodNotAllowed = true

	if s.container.Conf.Env == "localhost" {
		router.Use(gin.Logger()) //debug logger local development
	}

	router.Use(gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		logrus.Error(err)
		s.panicRecoveryHandler(c.Writer, c.Request)
	}))
	router.NoRoute(gin.HandlerFunc(func(c *gin.Context) {
		s.notFoundHandler(c.Writer, c.Request)
	}))
	router.NoMethod(gin.HandlerFunc(func(c *gin.Context) {
		s.methodNotAllowedHandler(c.Writer, c.Request)
	}))

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

// notFoundHandler handles API server response when endpoint couldn't be found
func (s *Server) notFoundHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	res := response.New(http.StatusNotFound, response.StatusError, make([]response.ApiError, 0), nil)
	res.SetErrorMessages(res.GetErrorMessageSlice(response.ErrorRouteNotFound, "", "Route Not Found."))
	res.WriteJSONResponse(w)

	return
}

// panicRecoveryHandler : in case of unexpected 'panic' call during process, return a custom 500 internal server error instead of nothing.
func (s *Server) panicRecoveryHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	gin.Recovery()

	res := response.New(http.StatusInternalServerError, response.StatusError, make([]response.ApiError, 0), nil)
	res.SetErrorMessages(res.GetErrorMessageSlice(response.ErrorInternalServerError, "", "Internal Server Error"))
	res.WriteJSONResponse(w)

	return
}

// methodNotAllowedHandler handles API server response when endpoint method is not allowed
func (s *Server) methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	res := response.New(http.StatusMethodNotAllowed, response.StatusError, make([]response.ApiError, 0), nil)
	res.SetErrorMessages(res.GetErrorMessageSlice(response.ErrorMethodNotAllowed, "", "Method is not allowed."))
	res.WriteJSONResponse(w)

	return
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
