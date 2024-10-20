<h1>Ghoul</h1>
<img width="200px" height="auto" src="https://raw.githubusercontent.com/MariaLetta/free-gophers-pack/refs/heads/master/characters/png/1.png"/>

> [!WARNING]
> This library is still a work in progress and primarily tailored to my personal use cases.

A lightweight, no-fuss router built on top of Go's standard HTTP module.

## Installation

``` sh
go get -u github.com/julienrull/ghoul
```

## Routes

You can define routes and nest them within other routes to create a hierarchical structure. A group represents a router that organizes related sub-routes or sub-groups like a tree structure where groups act as internal nodes, while individual routes serve as the leaves. The root of this tree-like structure is the app route, which functions as the top-level group.

```go
    is_auth = false
    app := ghoul.New()
    app.Get("/landing", func(ctx ghoul.Ctx) error {
        ctx.Response.Write([]byte("landing"))
        return nil
    })
    guest := app.Group("/guest")
    guest.Get("/signin", func(ctx ghoul.Ctx) error {
        ctx.Response.Write([]byte("signin"))
        return nil
    }).Post("/signin", func(ctx ghoul.Ctx) error {
        is_auth = true
        ctx.Redirect("/user/home", http.StatusSeeOther)
        return nil
    })
    user := app.Group("/user")
    user.Get("/home", func(ctx ghoul.Ctx) error {
        ctx.Response.Write([]byte("home"))
        return nil
    })
    app.ListenAndServe()
```

## Middleware

Middleware allows you to execute logic between incoming requests. Multiple middlewares can be layered, with the most recently registered middleware being executed first. If a middleware does not invoke the Next() method after completing its logic, no further routes or middleware will be processed.

### Use()

The Use() method applies middleware globally, but their execution is conditional on the request URI prefix matching the path specified in Use(). The path is automatically rewritten to include the group's base URI as a prefix.

### Group()

The Group() method allows you to define middleware that applies only to the routes within that group and its nested sub-groups and sub-routes.

### Get(), Post(), etc.

Standard HTTP methods (such as Get(), Post(), etc.) can also take middlewares, but these will apply specifically to the individual route where they are defined.

```go
    func auth_guard_middleware(ctx ghoul.Ctx) error {
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
    app.Get("/landing", func(ctx ghoul.Ctx) error {
        ctx.Response.Write([]byte("landing"))
        return nil
    }, auth_guard_middleware, ...)
    guest := app.Group("/guest")
    guest.Get("/signin", func(ctx ghoul.Ctx) error {
        ctx.Response.Write([]byte("signin"))
        return nil
    }).Use(auth_guard_middleware, ...)
    user := app.Group("/user", auth_guard_middleware, ...)
    user.Get("/home", func(ctx ghoul.Ctx) error {
        ctx.Response.Write([]byte("home"))
        return nil
    })
```
## Templating

The library includes its own templating system built on top of Go templates. You can define layouts that wrap around the rendered templates.

```go
    app.Get("/signin", func(ctx ghoul.Ctx) error{
        ctx.Render("signin", map[string]any{"data": "some data"}, "layouts/main")
        return nil
    })
```
