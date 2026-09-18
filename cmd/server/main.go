package main

import (
	"flag"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	port := flag.String("port", "8000", "server port")
	flag.Parse()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		time.Sleep(time.Second * 5)
		c.JSON(200, gin.H{
			"message": "pong",
			"port":    *port,
		})
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"port":    *port,
		})
	})

	r.Run(":" + *port)
}
