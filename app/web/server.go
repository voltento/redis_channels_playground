package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"redis_channels_playground/app/config"
	"redis_channels_playground/app/service"
)

type Server struct {
	cfg     config.Config
	router  *gin.Engine
	service *service.Service
	hc      *service.HealthChecker
}

func (s *Server) Router() *gin.Engine {
	return s.router
}

func NewServer(cfg config.Config, search *service.Service, hc *service.HealthChecker) *Server {
	router := gin.Default()

	s := &Server{cfg: cfg, router: router, service: search, hc: hc}
	return s
}

func (s *Server) bind() {

	s.router.GET("/api/v1/hw",
		func(c *gin.Context) {
			s.service.PostMessage()
			c.JSON(http.StatusOK, map[string]interface{}{"message": "hello world"})
		},
	)

	s.router.GET("/health", func(c *gin.Context) {
		if err := s.hc.Healthy(); err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, "healthy")
		}
	})
}

func (s *Server) Run() error {
	s.bind()
	return s.router.Run(s.cfg.ServiceHost)
}
