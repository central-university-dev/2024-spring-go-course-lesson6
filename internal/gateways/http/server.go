package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Server struct {
	host   string
	port   uint16
	router *gin.Engine
}

// TODO реализовать web-сервис согласно заданию
func NewServer(options ...func(*Server)) *Server {
	r := gin.Default()
	setupRouter(r)

	s := &Server{router: r, host: "localhost", port: 8080}
	for _, o := range options {
		o(s)
	}

	return s
}

func WithHost(host string) func(*Server) {
	return func(s *Server) {
		s.host = host
	}
}

func WithPort(port uint16) func(*Server) {
	return func(s *Server) {
		s.port = port
	}
}

func (s *Server) Run() error {
	return s.router.Run(s.host + ":" + strconv.FormatUint(uint64(s.port), 10))
}
