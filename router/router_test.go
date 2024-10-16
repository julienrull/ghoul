package router 

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

func GetServer() *Server {
    s := NewServer()
    s.Renderer.Folder = "./"
    s.Get("/simpleroute", func(ctx Context) (Context, error) {
        ctx.Response.Header().Add("Content-Type", "text/html")
        ctx.Render("test", map[string]any{})
        return ctx, nil
    })
    s.Post("/simpleroute", func(ctx Context) (Context, error) {
        ctx.Response.Write([]byte("simpleroutepost"))
        return ctx, nil
    })

    nestedRoute := s.Group("/nestedRoute")
    nestedRoute.Get("/nested", func(ctx Context) (Context, error) {
        ctx.Response.Write([]byte("nestedget"))
        return ctx, nil
    })
    nestedRoute.Post("/nested", func(ctx Context) (Context, error) {
        ctx.Response.Write([]byte("nestedpost"))
        return ctx, nil
    })

    subnestedroute := nestedRoute.Group("/subnestedroute")
    subnestedroute.Get("/subnested", func(ctx Context) (Context, error) {
        ctx.Response.Write([]byte("subnestedget"))
        return ctx, nil
    })
    subnestedroute.Post("/subnested", func(ctx Context) (Context, error) {
        ctx.Response.Write([]byte("subnestedpost"))
        return ctx, nil
    })

    simplemiddleware := s.Group("/simplemiddleware")
    simplemiddleware.Get("/notexecuted", func(ctx Context) (Context, error) {
        ctx.Response.Write([]byte("notexecutedget"))
        return ctx, nil
    })
    simplemiddleware.Post("/notexecuted", func(ctx Context) (Context, error) {
        ctx.Response.Write([]byte("notexecutedpost"))
        return ctx, nil
    })
    simplemiddleware.Use(func(ctx Context) (Context, error) {
        ctx.Response.Write([]byte("simplemiddleware"))
        return ctx, nil
    })

    firstmiddleware := s.Group("/firstmiddleware")
    secondmiddleware := firstmiddleware.Group("/secondsubmiddleware")
    firstmiddleware.Use(func(ctx Context) (Context, error) {
        ctx.Next()
        return ctx, nil
    })
    secondmiddleware.Get("/notexecuted", func(ctx Context) (Context, error) {
        ctx.Response.Write([]byte("notexecutedget"))
        return ctx, nil
    })
    secondmiddleware.Post("/notexecuted", func(ctx Context) (Context, error) {
        ctx.Response.Write([]byte("notexecutedpost"))
        return ctx, nil
    })
    secondmiddleware.Use(func(ctx Context) (Context, error) {
        ctx.Response.Write([]byte("submiddleware"))
        return ctx, nil
    })
    s.Get("/redirectroute", func(ctx Context) (Context, error) {
        ctx.Redirect("/redirecttarget", http.StatusSeeOther)
        return ctx, nil
    })
    s.Get("/redirecttarget", func(ctx Context) (Context, error) {
        ctx.Response.Write([]byte("redirecttarget"))
        return ctx, nil
    })
    s.Use(func(ctx Context) (Context, error) {
        ctx.Next()
        return ctx, nil
    })
    nextroute := s.Group("/nextroute")
    nextroute.Use(func(ctx Context) (Context, error) {
        ctx.Next()
        return ctx, nil
    })
    nextroute.Get("/subnextroute", func(ctx Context) (Context, error) {
        ctx.Response.Write([]byte("subnextroute"))
        return ctx, nil
    })
    s.Run()
    return s
}

func TestSimpleRoute(t *testing.T) {
    s := GetServer()
    st := httptest.NewServer(s.Handle)
    defer st.Close()   
    c := NewClient(st.URL)

    res1, _ := c.TestGetQuery("/simpleroute")
    trimedRes1 := strings.TrimSpace(res1)
    res2, _ := c.TestPostQuery("/simpleroute")

    res1exp := "<h1>simplerouteget</h1>"
    res2exp := "simpleroutepost"
    if trimedRes1 != res1exp {
        t.Errorf("expected res to be %s got %s", res1exp, trimedRes1)
    }
    if res2 != res2exp {
        t.Errorf("expected res to be %s got %s", res2exp, res2)
    }
}


func TestNestedRoute(t *testing.T) {
    s := GetServer()
    defer s.Server.Close()   
    st := httptest.NewServer(s.Handle)
    defer st.Close()   
    c := NewClient(st.URL)
    res1, _ := c.TestGetQuery("/nestedRoute/nested")
    res2, _ := c.TestPostQuery("/nestedRoute/nested")
    res3, _ := c.TestGetQuery("/nestedRoute/subnestedroute/subnested")
    res4, _ := c.TestPostQuery("/nestedRoute/subnestedroute/subnested")
    res1exp := "nestedget"
    res2exp := "nestedpost"
    res3exp := "subnestedget"
    res4exp := "subnestedpost"
    if res1 != res1exp {
        t.Errorf("expected res to be %s got %s", res1exp, res1)
    }
    if res2 != res2exp {
        t.Errorf("expected res to be %s got %s", res2exp, res2)
    }
    if res3 != res3exp {
        t.Errorf("expected res to be %s got %s", res3exp, res3)
    }
    if res4 != res4exp {
        t.Errorf("expected res to be %s got %s", res4exp, res4)
    }
}


func TestSimpleMiddlewares(t *testing.T) {
    s := GetServer()
    defer s.Server.Close()   
    st := httptest.NewServer(s.Handle)
    defer st.Close()   
    c := NewClient(st.URL)
    res1, _ := c.TestGetQuery("/simplemiddleware")
    res2, _ := c.TestPostQuery("/simplemiddleware")
    res1exp := "simplemiddleware"
    res2exp := "simplemiddleware"
    if res1 != res1exp {
        t.Errorf("expected res to be %s got %s", res1exp, res1)
    }
    if res2 != res2exp {
        t.Errorf("expected res to be %s got %s", res2exp, res2)
    }
}

func TestSubMiddlewares(t *testing.T) {
    s := GetServer()
    defer s.Server.Close()   
    st := httptest.NewServer(s.Handle)
    defer st.Close()   
    c := NewClient(st.URL)
    res1, _ := c.TestGetQuery("/firstmiddleware/secondmiddleware")
    res2, _ := c.TestPostQuery("/firstmiddleware/secondmiddleware")
    res1exp := "submiddleware"
    res2exp := "submiddleware"
    if res1 != res1exp {
        t.Errorf("expected res to be %s got %s", res1exp, res1)
    }
    if res2 != res2exp {
        t.Errorf("expected res to be %s got %s", res2exp, res2)
    }
}

func TestNext(t *testing.T) {
    s := GetServer()
    defer s.Server.Close()   
    st := httptest.NewServer(s.Handle)
    defer st.Close()   
    c := NewClient(st.URL)
    res2, _ := c.TestGetQuery("/nextroute/subnextroute")
    res2exp := "subnextroute"
    if res2 != res2exp {
        t.Errorf("expected res to be %s got %s", res2exp, res2)
    }
}

func TestRedirection(t *testing.T) {
    s := GetServer(); defer s.Server.Close()   
    st := httptest.NewServer(s.Handle)
    defer st.Close()   
    c := NewClient(st.URL)
    res1, _ := c.TestGetQuery("/redirectroute")
    res1exp := "redirecttarget"
    if res1 != res1exp {
        t.Errorf("expected res to be %s got %s", res1exp, res1)
    }
}
