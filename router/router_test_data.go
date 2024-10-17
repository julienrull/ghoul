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

func GetServer() *Router {
    app := New()
    app.Renderer = renderer.NewRenderer("./", ".html")
    app.Get("/simpleroute", func(ctx Ctx) error {
        ctx.Response.Header().Add("Content-Type", "text/html")
        ctx.Render("test", map[string]any{})
        return nil
    })
    app.Post("/simpleroute", func(ctx Ctx) error {
        ctx.Response.Write([]byte("simpleroutepost"))
        return nil
    })

    nestedRoute := app.Group("/nestedRoute", func(c Ctx) error {
        c.Next()
        return nil
    })
    nestedRoute.Get("/nested", func(ctx Ctx) error {
        ctx.Response.Write([]byte("nestedget"))
        return nil
    })
    nestedRoute.Post("/nested", func(ctx Ctx) error {
        ctx.Response.Write([]byte("nestedpost"))
        return nil
    })

    subnestedroute := nestedRoute.Group("/subnestedroute", func(c Ctx) error {
        c.Next()
        return nil
    })
    subnestedroute.Get("/subnested", func(ctx Ctx) error {
        ctx.Response.Write([]byte("subnestedget"))
        return nil
    })
    subnestedroute.Post("/subnested", func(ctx Ctx) error {
        ctx.Response.Write([]byte("subnestedpost"))
        return nil
    })

    simplemiddleware := app.Group("/simplemiddleware", func(ctx Ctx) error {
        ctx.Response.Write([]byte("simplemiddleware"))
        //ctx.Next()
        return nil
    })
    simplemiddleware.Get("/notexecuted", func(ctx Ctx) error {
        ctx.Response.Write([]byte("notexecutedget"))
        return nil
    })
    simplemiddleware.Post("/notexecuted", func(ctx Ctx) error {
        ctx.Response.Write([]byte("notexecutedpost"))
        return nil
    })

    firstmiddleware := app.Group("/firstmiddleware", func(ctx Ctx) error {
        fmt.Println("GROUP")
        //ctx.Next()
        return nil
    })
    //firstmiddleware.Use(func(ctx Ctx) error {
    //    fmt.Println("USE")
    //    ctx.Next()
    //    return nil
    //})

    secondmiddleware := firstmiddleware.Group("/secondmiddleware")
    secondmiddleware.Get("/notexecuted", func(ctx Ctx) error {
        ctx.Response.Write([]byte("notexecutedget"))
        return nil
    })
    secondmiddleware.Post("/notexecuted", func(ctx Ctx) error {
        ctx.Response.Write([]byte("notexecutedpost"))
        return nil
    })
    app.Get("/redirectroute", func(ctx Ctx) error {
        ctx.Redirect("/redirecttarget", http.StatusSeeOther)
        return nil
    })
    app.Get("/redirecttarget", func(ctx Ctx) error {
        ctx.Response.Write([]byte("redirecttarget"))
        return nil
    })
    app.Use(func(ctx Ctx) error {
        ctx.Next()
        return nil
    })
    nextroute := app.Group("/nextroute", func(ctx Ctx) error {
        ctx.Response.Write([]byte("subnestedpost"))
        return nil
    })
    nextroute.Use(func(ctx Ctx) error {
        ctx.Next()
        return nil
    })
    nextroute.Get("/subnextroute", func(ctx Ctx) error {
        ctx.Response.Write([]byte("subnextroute"))
        return nil
    })
    app.PostInit()
    return app
}
