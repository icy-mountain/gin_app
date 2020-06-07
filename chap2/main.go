package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/matryer/goblueprints/chapter1/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

func authCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := trace.New(os.Stdout)
		tracer.Trace("in auth check!")
		_, err := c.Request.Cookie("auth")
		if err == http.ErrNoCookie {
			tracer.Trace("redirect!")
			c.Redirect(307, "http://localhost:8080/login")
			// return
			c.AbortWithStatus(307)
		}
		if err != nil {
			// some other error
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			c.AbortWithStatus(500)
		}
		c.Next()
	}
}

func main() {
	port := flag.String("port", ":8080", "The host of the application.")
	flag.Parse()
	egn := gin.Default()
	room := newRoom()
	room.tracer = trace.New(os.Stdout)
	go room.run()

	gomniauth.SetSecurityKey("98dfbg7iu2nb4uywevihjw4tuiyub34noilk")
	gomniauth.WithProviders(
		google.New("160129032797-4cpkfrf3dd9lvfmr2ibpg473ampn4f71.apps.googleusercontent.com", "UTUNVwyoJPt8MTWrdQMmaw3T", "http://localhost:8080/auth/callback/google"),
	)
	egn.LoadHTMLGlob("templates/*.html")
	egn.GET("/chat", authCheck(), func(c *gin.Context) {
		tracer := trace.New(os.Stdout)
		tracer.Trace("in chat!")
		data := map[string]interface{}{
			"Host": c.Request.Host,
		}
		if authCookie, err := c.Request.Cookie("auth"); err == nil {
			data["UserData"] = objx.MustFromBase64(authCookie.Value)
		}
		c.HTML(http.StatusOK, "chat.html", data)
	})
	egn.GET("/room", authCheck(), room.ServeHTTP)
	egn.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	egn.GET("/auth/:action/:provider", loginHandler)
	egn.Run(*port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
