package ghoul

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type Client struct {
    url string
}

func NewClient(url string) Client {
    return Client{url}
}


func (c Client) TestQuery(method string, path string) (string, error) {
	req, err := http.NewRequest(method, c.url + path, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
    return string(resBody), nil
}


var is_auth = false
var is_admin = false

func auth_guard_middleware(ctx Ctx) error {
   if !is_auth {
       if ctx.Request().URL.RequestURI() == "/users/1/home" {
        ctx.Redirect("/guest/signin", http.StatusSeeOther)
       }
   }else{
    if ctx.Request().URL.RequestURI() == "/guest/signin" {
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
        ctx.Write([]byte("landing"))
        return nil
    }).Head("/landing", func(ctx Ctx) error {
        ctx.Write([]byte("landinghead"))
        return nil
    }).Put("/landing", func(ctx Ctx) error {
        ctx.Write([]byte("landingput"))
        return nil
    }).Delete("/landing", func(ctx Ctx) error {
        ctx.Write([]byte("landingdelete"))
        return nil
    }).Connect("/landing", func(ctx Ctx) error {
        ctx.Write([]byte("landingconnect"))
        return nil
    }).Options("/landing", func(ctx Ctx) error {
        ctx.Write([]byte("landingoptions"))
        return nil
    }).Trace("/landing", func(ctx Ctx) error {
        ctx.Write([]byte("landingtrace"))
        return nil
    }).Patch("/landing", func(ctx Ctx) error {
        ctx.Write([]byte("landingpatch"))
        return nil
    }).All("/all", func(ctx Ctx) error {
        ctx.Write([]byte("all"))
        return nil
    })

    guest := app.Group("/guest")
    guest.Get("/signin", func(ctx Ctx) error {
        ctx.Write([]byte("signin"))
        return nil
    }).Post("/signin", func(ctx Ctx) error {
        is_auth = true
        ctx.Redirect("/users/1/home", http.StatusSeeOther)
        return nil
    }).Use(auth_guard_middleware).Use("/admin", log_middleware, admin_middleware).Use([]string{"/stats", "/secret"}, log_middleware)

    app.Get("/users", func(ctx Ctx) error {
        ctx.Write([]byte("users"))
        return nil
    }, auth_guard_middleware)

    users := app.Group("/users/{userid}", auth_guard_middleware)

    users.Get("/home", func(ctx Ctx) error {
        ctx.Write([]byte("user n°" + ctx.Request().PathValue("userid")))
        return nil
    })
    posts := users.Group("/posts")
    posts.Get("/{postid}", func(ctx Ctx) error {
        //ctx.Response.Write([]byte("post n°" + ctx.Request.PathValue("postid")))
        ctx.Render("body", map[string]any{"postid": ctx.Request().PathValue("postid")}, "layouts/main")
        return nil
    }, log_middleware)
    posts.Get("/archives/{id}", func(ctx Ctx) error {
        ctx.Write([]byte("post n°" + ctx.Request().PathValue("archivepostid")))
        return nil
    }, log_middleware, log_middleware)
    app.PostInit()
    return app
}
