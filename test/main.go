package main

import (
	"flag"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	host := flag.String("host", ":8080", "The host of the application.")
	r := gin.Default()
	room := newRoom()
	go room.run()
	r.LoadHTMLGlob("templates/*.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", gin.H{
			"Host": "localhost" + *host,
		})
	})
	r.GET("/room", room.ServeHTTP)
	r.Run(*host) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
