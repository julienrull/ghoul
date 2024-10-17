package main

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/julienrull/ghoul/v1/renderer"
)

type Middleware = func(http.Handler) http.Handler
type MiddlewareHandler = ContextHandler
type ContextHandler = func(Ctx) error

type Ctx struct {
    Request     *http.Request
    Response    http.ResponseWriter
    Status      int
    Handle      http.Handler
    Renderer    *renderer.Renderer
}

func (c *Ctx) Next() {
    if c.Handle != nil {
        c.Handle.ServeHTTP(c.Response, c.Request)
    }
}

func (c *Ctx) Redirect(path string, status int) {
    http.Redirect(c.Response, c.Request, path, status)
}

func (c *Ctx) Render(tmplName string, data map[string]any, layouts ...string) {
    c.Renderer.Render(c.Response, tmplName, data, layouts...) 
}



type Router struct {
    Server      *http.Server
    BaseUrl     string
    Handle      http.Handler
    Mux         *http.ServeMux
    Childs      []*Router     
    Renderer    *renderer.Renderer
}

func New() *Router {
    mux := http.NewServeMux()
    return &Router{
        Handle: mux,
        Mux: mux,
        BaseUrl: "",
        Server: &http.Server{
            Addr:           "localhost:3000",
            ReadTimeout:    10 * time.Second,
            WriteTimeout:   10 * time.Second,
            MaxHeaderBytes: 1 << 20,
            Handler: nil,
        },
    }
}

func (r *Router) Group(path string, middlewares ...ContextHandler) *Router{
    var handle ContextHandler = nil
    var handles []ContextHandler = nil
    if len(middlewares) > 1 {
       handle =  middlewares[0]
       handles =  middlewares[1:]
    } else if len(middlewares) == 1 {
       handle =  middlewares[0]
    }    
    return r.Add("GROUP", path, handle, handles...)
}

func (r *Router) Get(path string, handler ContextHandler, middlewares ...ContextHandler) *Router {
    return r.Add("GET", path, handler, middlewares...)
}

func (r *Router) Post(path string, handler ContextHandler, middlewares ...ContextHandler) *Router {
    return r.Add("POST", path, handler, middlewares...)
}

func (r *Router) CreateStack(middlewares ...ContextHandler) Middleware {
    return func(next http.Handler) http.Handler  {
        for i := len(middlewares)-1; i >= 0; i-- {
           h := middlewares[i]
           prenext := next
           next = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
               h(Ctx{
                   Response: w,
                   Request:  req,
                   Renderer: r.Renderer,
                   Handle: prenext,
               }) 
           })
        }
        return next
    }
}

func (r *Router) register(handlers []ContextHandler) http.Handler {
    var handler http.Handler = nil
    if handlers != nil {
        if len(handlers) > 1 {
            stack := r.CreateStack(handlers[1:]...)
            handler = stack(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                handlers[0](Ctx{
                    Response: w,
                    Request:  req,
                    Renderer: r.Renderer,
                    Handle: nil,
                }) 
            }))
        }else if len(handlers) == 1{
            handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                handlers[0](Ctx{
                    Response: w,
                    Request:  req,
                    Renderer: r.Renderer,
                    Handle: nil,
                }) 
            })
        }
    }
    if handler == nil {
        panic("no handlers provided")
    }
    return handler
}


func (r *Router) Add(method string, path string, handler ContextHandler, middlewares ...ContextHandler) *Router {
    handlers := append([]ContextHandler{handler}, middlewares...)
    handle := r.register(handlers)
    if method == "USE"{
        path = r.BaseUrl + path + "/"
    } else if method == "GROUP" {
        path = r.BaseUrl + path 
        mux := http.NewServeMux()
        gr := &Router{
            Handle: mux,
            Mux: mux,
            BaseUrl: path,
        }
        gr.Mux.Handle(path + "/", handle)
        r.Mux.Handle(path + "/", gr.Handle)
        return gr
    }else{
        path = method + " " + r.BaseUrl + path
    }
    r.Mux.Handle(path, handle)
    return r
}

func (r *Router) Use(args ...any) *Router {
	var prefix string = "/"
	//var subRouter *Router
	var prefixes []string
	var handlers []ContextHandler

	for i := 0; i < len(args); i++ {
		switch arg := args[i].(type) {
		case string:
			prefix = arg
		//case *Router:
		//	subRouter = arg
		case []string:
			prefixes = arg
		case ContextHandler:
			handlers = append(handlers, arg)
		default:
			panic(fmt.Sprintf("use: invalid handler %v\n", reflect.TypeOf(arg)))
		}
	}
	if len(prefixes) == 0 {
		prefixes = append(prefixes, prefix)
	}
	for _, prefix := range prefixes {
        r.Add("USE", prefix, handlers[0], handlers[1:]...)
	}
	return r
}

func (r *Router) PostInit() {
    if r.Server != nil {
        r.Server.Handler = r.Mux
        return
    } 
    panic("Can't listen and serve from sub router")
}

func (r *Router) ListenAndServe() {
    r.PostInit()
    r.Server.ListenAndServe()
}

func main() {
    app := New()
    app.Get("/hello", func(ctx Ctx) error {
        ctx.Response.Write([]byte("Hello Ghoul"))
        return nil
    }, func(c Ctx) error {
        c.Next() 
        return nil
    })
    good := app.Group("/good", func(ctx Ctx) error {
        ctx.Next()
        return nil
    })
    good.Get("/so6", func(ctx Ctx) error {
        ctx.Response.Write([]byte("Hello Good So6"))
        return nil
    })
    app.ListenAndServe()
}
