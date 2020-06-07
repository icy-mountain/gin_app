package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/matryer/goblueprints/chapter1/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)

var egn = gin.Default()

// loginHandler handles the third-party login process.
func loginHandler(c *gin.Context) {
	action := c.Param("action")
	provider := c.Param("provider")
	switch action {
	case "login":

		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(c.Writer, fmt.Sprintf("Error when trying to get provider %s: %s", provider, err), http.StatusBadRequest)
			return
		}

		loginURL, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			http.Error(c.Writer, fmt.Sprintf("Error when trying to GetBeginAuthURL for %s: %s", provider, err), http.StatusInternalServerError)
			return
		}

		tracer := trace.New(os.Stdout)
		tracer.Trace("in login section!")
		c.Redirect(307, loginURL)
	case "callback":
		tracer := trace.New(os.Stdout)
		tracer.Trace("callback section!")
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(c.Writer, fmt.Sprintf("Error when trying to get provider %s: %s", provider, err), http.StatusBadRequest)
			return
		}

		// get the credentials
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(c.Request.URL.RawQuery))
		if err != nil {
			http.Error(c.Writer, fmt.Sprintf("Error when trying to complete auth for %s: %s", provider, err), http.StatusInternalServerError)
			return
		}

		// get the user
		user, err := provider.GetUser(creds)
		if err != nil {
			http.Error(c.Writer, fmt.Sprintf("Error when trying to get user from %s: %s", provider, err), http.StatusInternalServerError)
			return
		}

		// save some data
		authCookieValue := objx.New(map[string]interface{}{
			"name": user.Name(),
		}).MustBase64()
		http.SetCookie(c.Writer, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/"})
		// func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
		// c.SetCookie("auth", authCookieValue, 24*60*60, "/", "", false, false)
		_, err = c.Request.Cookie("auth")
		if err == nil {
			tracer := trace.New(os.Stdout)
			tracer.Trace("cookie OK!")
		} else {
			tracer := trace.New(os.Stdout)
			tracer.Trace("cookie NG!")
		}
		// c.Writer.Header().Set("Location", "/chat")
		// c.Writer.WriteHeader(http.StatusTemporaryRedirect)
		c.Redirect(307, "http://localhost:8080/chat")
	default:
		c.Writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(c.Writer, "Auth action %s not supported", action)
	}
}
