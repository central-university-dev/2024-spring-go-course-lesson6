package http

import "github.com/gin-gonic/gin"

func setupRouter(r *gin.Engine, _ UseCases, _ *WebSocketHandler) {
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}
