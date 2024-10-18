package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/julienrull/ghoul/v1/renderer"
)

type Client struct {
    url string
}

func NewClient(url string) Client {
    return Client{url}
}

func (c Client) TestGetQuery(path string) (string, error) {
    res, _ := http.Get(c.url + path)
    defer res.Body.Close()
    out, _ := io.ReadAll(res.Body)
    return string(out), nil
}

func (c Client) TestPostQuery(path string) (string, error) {
    res, _ := http.Post(c.url + path, "", nil)
    defer res.Body.Close()
    out, _ := io.ReadAll(res.Body)
    return string(out), nil
}


var is_auth = false

func auth_guard_middleware(ctx Ctx) error {
   fmt.Println(ctx.Request.URL.RequestURI())
   if !is_auth {
       if ctx.Request.URL.RequestURI() == "/user/home" {
        ctx.Redirect("/guest/signin", http.StatusSeeOther)
       }
   }else{
    if ctx.Request.URL.RequestURI() == "/guest/signin" {
       ctx.Redirect("/user/home", http.StatusSeeOther)
    }
   }
   ctx.Next()
   return nil 
}

func GetServer() *Router {
    is_auth = false
    app := New()
    app.Renderer = renderer.NewRenderer("./", ".html")
    app.Get("/landing", func(ctx Ctx) error {
        ctx.Response.Write([]byte("landing"))
        return nil
    })
    guest := app.Group("/guest")
    guest.Get("/signin", func(ctx Ctx) error {
        ctx.Response.Write([]byte("signin"))
        return nil
    }).Post("/signin", func(ctx Ctx) error {
        is_auth = true
        ctx.Redirect("/user/home", http.StatusSeeOther)
        return nil
    }).Use(auth_guard_middleware)
    user := app.Group("/user", auth_guard_middleware)
    user.Get("/home", func(ctx Ctx) error {
        ctx.Response.Write([]byte("home"))
        return nil
    })
    app.PostInit()
    return app
}
