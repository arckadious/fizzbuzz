package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	cst "github.com/arckadious/fizzbuzz/constant"
	"github.com/arckadious/fizzbuzz/container"
	"github.com/arckadious/fizzbuzz/response"

	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"
)

const (
	// Timeout delay and graceful shutdown deadline
	Timeout     = time.Second * 180
	IdleTimeout = 60
	pingURL     = "/ping"
)

// Server httpServer struct
type Server struct {
	container *container.Container
}

// New create server
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

func (s *Server) Handler() *gin.Engine {

	router := gin.New()
	router.HandleMethodNotAllowed = true
	router.Use(s.Logger())
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

	//ping
	router.GET(pingURL, func(c *gin.Context) {
		c.String(http.StatusOK, "Ping OK !")
	})

	//documentation de l'api
	router.Static("/swagger", "./swaggerui")

	//api subrouter v1
	v1 := router.Group(cst.URLPrefixVersion)

	//fizzbuzz routes
	v1.POST(cst.FizzBaseURI, func(c *gin.Context) {
		s.container.FizzAction.HandleFizz(c.Writer, c.Request)
	})
	v1.GET(cst.FizzStatisticsURI, func(c *gin.Context) {
		s.container.FizzAction.HandleStatistics(c.Writer, c.Request)
	})

	return router

}

// notFoundHandler
func (s *Server) notFoundHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	res := response.NewApiResponse(http.StatusNotFound, response.StatusError, make([]response.ApiError, 0), nil)
	res.SetErrorMessages(res.GetErrorMessageSlice(response.ErrorRouteNotFound, "", "Route Not Found."))
	res.WriteJsonResponse(w)

	return
}

// panicRecoveryHandler (in case of 'panic' call during process, return a custom 500 internal server error instead of nothing.)
func (s *Server) panicRecoveryHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	gin.Recovery()

	res := response.NewApiResponse(http.StatusInternalServerError, response.StatusError, make([]response.ApiError, 0), nil)
	res.SetErrorMessages(res.GetErrorMessageSlice(response.ErrorInternalServerError, "", "Internal Server Error"))
	res.WriteJsonResponse(w)

	return
}

// methodNotAllowedHandler
func (s *Server) methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	res := response.NewApiResponse(http.StatusMethodNotAllowed, response.StatusError, make([]response.ApiError, 0), nil)
	res.SetErrorMessages(res.GetErrorMessageSlice(response.ErrorMethodNotAllowed, "", "Method is not allowed."))
	res.WriteJsonResponse(w)

	return
}

// audit requests and response to DB
func (s *Server) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {

		// before request
		//TO DO
		//go s.container.RepoFizz.AuditRequest()
		logrus.Info("TEST")
		logrus.Info("TEST2")
		// c.Next()
		logrus.Info("TEST3")
		logrus.Info("TEST4")
		// after request
		//TO DO
		//go s.container.RepoFizz.AuditResponse()

	}
}
