package ghoul

import (
	"io"
	"net/http"
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
var is_admin = false

func auth_guard_middleware(ctx Ctx) error {
   if !is_auth {
       if ctx.Request.URL.RequestURI() == "/users/1/home" {
        ctx.Redirect("/guest/signin", http.StatusSeeOther)
       }
   }else{
    if ctx.Request.URL.RequestURI() == "/guest/signin" {
       ctx.Redirect("/users/1/home", http.StatusSeeOther)
    }
   }
   ctx.Next()
   return nil 
}


func admin_middleware(ctx Ctx) error {
   ctx.Next()
   return nil 
}

func log_middleware(ctx Ctx) error {
   ctx.Next()
   return nil 
}

func GetServer() *Router {
    is_auth = false
    is_admin = false
    app := New()
    app.Renderer = NewRenderer("./views/", ".html")
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
        ctx.Redirect("/users/1/home", http.StatusSeeOther)
        return nil
    }).Use(auth_guard_middleware)

    app.Get("/users", func(ctx Ctx) error {
        ctx.Response.Write([]byte("users"))
        return nil
    }, auth_guard_middleware)

    users := app.Group("/users/{userid}", auth_guard_middleware)

    users.Get("/home", func(ctx Ctx) error {
        ctx.Response.Write([]byte("user n°" + ctx.Request.PathValue("userid")))
        return nil
    })
    posts := users.Group("/posts")
    posts.Get("/{postid}", func(ctx Ctx) error {
        //ctx.Response.Write([]byte("post n°" + ctx.Request.PathValue("postid")))
        ctx.Render("body", map[string]any{"postid": ctx.Request.PathValue("postid")}, "layouts/main")
        return nil
    }, log_middleware)
    posts.Get("/archives/{id}", func(ctx Ctx) error {
        ctx.Response.Write([]byte("post n°" + ctx.Request.PathValue("archivepostid")))
        return nil
    }, log_middleware, log_middleware)
    app.PostInit()
    return app
}
