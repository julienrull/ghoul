<h1>Ghoul</h1>
<img width="200px" height="auto" src="https://raw.githubusercontent.com/MariaLetta/free-gophers-pack/refs/heads/master/characters/png/1.png"/>

> [!WARNING]
> This library is still a work in progress and primarily tailored to my personal use cases.

A lightweight, no-fuss router built on top of Go's standard HTTP module.

## Routes

You can create routes and nest them within other routes. A group represents a router that serves sub-routes or sub-groups, functioning like a tree structure. In this structure, groups are the nodes, and routes are the leaves. The server acts as the root group or router.

```go
    var server = router.NewServer()

    server.Get("/helloghoul", func(ctx router.Context) (router.Context, error){
        ctx.Response.Write([]byte("Hello Ghoul !"))
        return ctx, nil
    })
    
    greatghoul := server.Group("/greatghoul")

    greatghoul.Get("/midghoul", func(ctx router.Context) (router.Context, error){
        ctx.Response.Write([]byte("Mid Ghoul"))
        return ctx, nil
    })
    greatghoul.Get("/smallghoul", func(ctx router.Context) (router.Context, error){
        ctx.Response.Write([]byte("Small Ghoul"))
        return ctx, nil
    })
    server.ListenAndServe()
```

## Middlewares

You can apply middlewares to groups of routes. The Use() method on a group attaches middleware to that group, affecting only the group and its child routes. Middlewares can be nested, with the most recently registered middleware applied first. If a middleware doesn't call the Next() method after executing its logic, the group's routes and any nested routes or middlewares will not be processed.

```go
    var server = router.NewServer()

    server.Get("/helloghoul", func(ctx router.Context) (router.Context, error){
        ctx.Response.Write([]byte("Hello Ghoul !"))
        return ctx, nil
    })
    server.Use("/helloghoul", func(ctx router.Context) (router.Context, error){
        ctx.Next()
        return ctx, nil
    })
    
    server.ListenAndServe()
```

## Redirection

You can redirect one route to another route or a nested route. For example, you can handle sign-in failures by redirecting a POST request back to the same route as a GET request.

```go
    var server = router.NewServer()
    server.Get("/signin", func(ctx router.Context) (router.Context, error){
        ctx.Response.Write([]byte("Signin"))
        return ctx, nil
    })
    server.Post("/signin", func(ctx router.Context) (router.Context, error){
        ctx.Redirect("/guest/signin", http.StatusSeeOther) 
        return ctx, nil
    })
    server.ListenAndServe()
```

## Templating

The library includes its own templating system built on top of Go templates. You can define layouts that wrap around the rendered templates.

```go
    var server = router.NewServer()
    server.Get("/signin", func(ctx router.Context) (router.Context, error){
        ctx.Render("signin", map[string]any{"data": "some data"}, "layouts/main")
        return ctx, nil
    })
    server.ListenAndServe()
```
