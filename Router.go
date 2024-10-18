package ghoul

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type Middleware = func(http.Handler) http.Handler
type MiddlewareHandler = ContextHandler
type ContextHandler = func(Ctx) error

type Ctx struct {
    Request     *http.Request
    Response    http.ResponseWriter
    Status      int
    Handle      http.Handler
    Renderer    *Renderer
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
    Root        *Router 
    isRoot      bool
    BaseUrl     string
    Handle      http.Handler
    Mux         *http.ServeMux
    Renderer    *Renderer
}

func New() *Router {
    mux := http.NewServeMux()
    router := &Router{
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
        isRoot: true,
    }
    router.Root = router
    return router
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

func (r *Router) Group(path string, middlewares ...ContextHandler) *Router{
    return r.Add("GROUP", path, nil, middlewares...)
}

func (r *Router) Get(path string, handler ContextHandler, middlewares ...ContextHandler) *Router {
    return r.Add("GET", path, handler, middlewares...)
}

func (r *Router) Post(path string, handler ContextHandler, middlewares ...ContextHandler) *Router {
    return r.Add("POST", path, handler, middlewares...)
}

func (r *Router) Use(args ...any) *Router {
	var prefix string
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
        r.Add("USE", prefix, nil, handlers...)
	}
	return r
}


func (r *Router) Add(method string, path string, handler ContextHandler, middlewares ...ContextHandler) *Router {
    var router *Router = r
    var handle http.Handler = r.Handle
    var stack Middleware = nil
    path = r.BaseUrl + path 
    var main_handler ContextHandler = nil
    var midd_handlers []ContextHandler = []ContextHandler{}

    if handler != nil {
        main_handler = handler
    }
    if middlewares != nil {
        midd_handlers = middlewares
    }


    if main_handler != nil && len(midd_handlers) > 0 {
        stack = r.CreateStack(midd_handlers...) 
        path = method + " " + path
        handle = stack(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
            main_handler(Ctx{
                Response: w,
                Request:  req,
                Renderer: r.Renderer,
                Handle: nil,
            })
        }))
    } else if main_handler != nil {
        //stack = r.CreateStack(handlers...) 
        path = method + " " + path
        handle = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
            main_handler(Ctx{
                Response: w,
                Request:  req,
                Renderer: r.Renderer,
                Handle: nil,
            })
        })
    }else{
        if method == "GROUP" && len(midd_handlers) > 0 {
            stack = r.CreateStack(midd_handlers...) 
            mux := http.NewServeMux()
            router = &Router{
                Handle: mux,
                Mux: mux,
                BaseUrl: path,
                Server: r.Server,
                isRoot: false,
                Root: r.Root,
            }
            path += "/"
            handle = stack(router.Handle)
        }else if method == "GROUP"{
            mux := http.NewServeMux()
            router = &Router{
                Handle: mux,
                Mux: mux,
                BaseUrl: path,
                Server: r.Server,
                isRoot: false,
                Root: r.Root,
            }
            path += "/"
            handle = router.Handle
        }else { // method == "USE"
            stack = r.CreateStack(midd_handlers...) 
            root := r.Root.Handle
            r.Root.Handle = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                fmt.Println(req.URL.RequestURI(), path)
               if strings.HasPrefix(req.URL.RequestURI(), path) {
                    next := stack(root)
                    next.ServeHTTP(w, req)
                    return
               }
               root.ServeHTTP(w, req)
               return
            })
            return router
        }
    }    
    r.Mux.Handle(path, handle)
    return router
}


func (r *Router) PostInit() {
    if r.isRoot {
        r.Server.Handler = r.Handle
        return
    } 
    panic("Can't listen and serve from sub router")
}

func (r *Router) ListenAndServe() {
    r.PostInit()
    r.Server.ListenAndServe()
}
