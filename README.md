<h1>Ghoul</h1>
<img width="200px" height="auto" src="https://cdn.pixabay.com/photo/2016/03/31/20/56/evil-1296097_1280.png"/>

> [!WARNING]
> This library is unfinished and for my personal use cases first.

Simple router library on top of GO standard http module.

## Routes

You can create routes and nested routes.
A group represent a router serving sub routes or sub groups.
It works like a tree struct, groups are nodes and routes leafs.
The Server is the root group/router.

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

You can wrap groups in middlewares.
Group Use() method apply the middleware on itself.
Only it and its children will be concerned by the middlewares.
Middlewares can be nested, the last registered is the top one.
Group and their sub routes or sub middlewares will not be served if Next() method isn't called after middleware logic.

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

You can redirect a route to another routes or nested route.
Here is a redirection you can perform when handle signin failures.

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

Ghoul came with its own templating system on top of GO template.
You can define layouts template that will wrap the rendered layout.

```go
    var server = router.NewServer()
    server.Get("/signin", func(ctx router.Context) (router.Context, error){
        ctx.Render("signin", map[string]any{"data": "some data"}, "layouts/main")
        return ctx, nil
    })
    server.ListenAndServe()
```
