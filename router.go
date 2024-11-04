package ghoul

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

func (r *Router) Group(path string, middlewares ...ContextHandler) *Router{
    return r.Add("GROUP", path, nil, middlewares...)
}

func (r *Router) Get(path string, handler ContextHandler, middlewares ...ContextHandler) *Router {
    return r.Add(http.MethodGet, path, handler, middlewares...)
}

func (r *Router) Head(path string, handler ContextHandler, middlewares ...ContextHandler) *Router {
    return r.Add(http.MethodHead, path, handler, middlewares...)
}

func (r *Router) Post(path string, handler ContextHandler, middlewares ...ContextHandler) *Router {
    return r.Add(http.MethodPost, path, handler, middlewares...)
}

func (r *Router) Put(path string, handler ContextHandler, middlewares ...ContextHandler) *Router {
    return r.Add(http.MethodPut, path, handler, middlewares...)
}

func (r *Router) Delete(path string, handler ContextHandler, middlewares ...ContextHandler) *Router {
    return r.Add(http.MethodDelete, path, handler, middlewares...)
}

func (r *Router) Connect(path string, handler ContextHandler, middlewares ...ContextHandler) *Router {
    return r.Add(http.MethodConnect, path, handler, middlewares...)
}

func (r *Router) Options(path string, handler ContextHandler, middlewares ...ContextHandler) *Router {
    return r.Add(http.MethodOptions, path, handler, middlewares...)
}

func (r *Router) Trace(path string, handler ContextHandler, middlewares ...ContextHandler) *Router {
    return r.Add(http.MethodTrace, path, handler, middlewares...)
}

func (r *Router) Patch(path string, handler ContextHandler, middlewares ...ContextHandler) *Router {
    return r.Add(http.MethodPatch, path, handler, middlewares...)
}

func (r *Router) All(path string, handler ContextHandler, middlewares ...ContextHandler) *Router {
    return r.Add("", path, handler, middlewares...)
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

func (r *Router) CreateStack(middlewares ...ContextHandler) Middleware {
    return func(next http.Handler) http.Handler  {
        for i := len(middlewares)-1; i >= 0; i-- {
           h := middlewares[i]
           prenext := next
           next = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                NewCtx(
                    req,
                    w,
                    prenext,
                    r.Renderer,
                ).Exec(h)
           })
        }
        return next
    }
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
        if method != "" {
            path = method + " " + path
        }
        handle = stack(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
            NewCtx(
                req,
                w,
                nil,
                r.Renderer,
            ).Exec(main_handler)
        }))
    } else if main_handler != nil {
        //stack = r.CreateStack(handlers...) 
        path = method + " " + path
        handle = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
            NewCtx(
                req,
                w,
                nil,
                r.Renderer,
            ).Exec(main_handler)
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
                Renderer: r.Renderer,
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
                Renderer: r.Renderer,
            }
            path += "/"
            handle = router.Handle
        }else { // method == "USE"
            stack = r.CreateStack(midd_handlers...) 
            root := r.Root.Handle
            r.Root.Handle = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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

func (r *Router) ListenAndServeTLS(certFile, keyFile string) {
    r.PostInit()
	go func() {
		if err := r.Server.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
    r.Exit()
}

func (r *Router) ListenAndServe() {
    r.PostInit()
	go func() {
		if err := r.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
    r.Exit()
}

func (r *Router) Exit() {
    fmt.Println(`
 ▗▄▄▖▗▖ ▗▖ ▗▄▖ ▗▖ ▗▖▗▖   
▐▌   ▐▌ ▐▌▐▌ ▐▌▐▌ ▐▌▐▌   
▐▌▝▜▌▐▛▀▜▌▐▌ ▐▌▐▌ ▐▌▐▌   
▝▚▄▞▘▐▌ ▐▌▝▚▄▞▘▝▚▄▞▘▐▙▄▄▖`)
    fmt.Printf("\nSERVE ON : http://%s\n", r.Server.Addr)

	<-r.signalOut
	fmt.Print("Server Stopped\n")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := r.Server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v\n", err)
	}
    fmt.Print("Server Exited Properly\n")
}
