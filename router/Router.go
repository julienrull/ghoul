package router

import (
	"net/http"
	"time"
	"github.com/julienrull/ghoul/renderer"
)

type Middleware func(http.Handler) http.Handler

type MiddlewareHandler ContextHandler

//func CreateStack(xs ...Middleware) Middleware {
//    return func(next http.Handler) http.Handler  {
//        for i := len(xs)-1; i >= 0; i-- {
//            x := xs[i]
//            next  = x(next)
//        }
//        return next
//    }
//}

type Context struct {
    Request     *http.Request
    Response    http.ResponseWriter
    Status      int
    Handle      http.Handler
    Renderer    *renderer.Renderer
}

func (c *Context) Next() {
    c.Handle.ServeHTTP(c.Response, c.Request)
}

func (c *Context) Redirect(path string, status int) {
    http.Redirect(c.Response, c.Request, path, status)
}

func (c *Context) Render(tmplName string, data map[string]any, layouts ...string) {
    c.Renderer.Render(c.Response, tmplName, data, layouts...) 
}

type ContextHandler func(Context) (Context, error)

type ServerConfig struct {
    Host                string
    Port                string
    RequestTimeout      time.Duration 
    ResponseTimeout     time.Duration
    RequestHeaderSize   int
}

var serverConfig = ServerConfig {
    Host:               "localhost", 
    Port:               "3000", 
    RequestTimeout:     10 * time.Second,
    ResponseTimeout:    10 * time.Second,
    RequestHeaderSize:  1 << 20,
}

type Server struct {
    Renderer *renderer.Renderer
    Server  http.Server     
    Mux     *http.ServeMux     
    Handle  http.Handler     
    Childs  []*Router     
} 

func NewServer() *Server {
    mux := http.NewServeMux()
    renderer := renderer.NewRenderer("./views/", ".html")
    var handle http.Handler = mux
    server := &Server{
        Renderer: renderer,
        Server: http.Server{
            Addr:           serverConfig.Host + ":" + serverConfig.Port,
            ReadTimeout:    serverConfig.RequestTimeout,
            WriteTimeout:   serverConfig.ResponseTimeout,
            MaxHeaderBytes: serverConfig.RequestHeaderSize,
            Handler: handle,
        },
        Mux: mux,
        Handle: handle,
        Childs: make([]*Router, 0, 16),
    }
    return server
}

func (s *Server) Register(routers []*Router) []*Router {
    if len(routers) > 0{
        for _, router := range routers {
            subRouters := s.Register(router.Childs) 
            if len(subRouters) > 0 {
                for _, subRouter := range subRouters {
                    router.Mux.Handle(router.BaseUrl + "/", subRouter.Handle)
                }
            }
        }
    }
    return routers
}

func (s *Server) Run(){
    subRouters := s.Register(s.Childs)
    if len(subRouters) > 0 {
        for _, subRouter := range subRouters {
            s.Mux.Handle(subRouter.BaseUrl + "/", subRouter.Handle)
        }
    }
}

func (s *Server) ListenAndServe(){
    s.Run()
    s.Server.ListenAndServe()
}

func (s *Server) Route(method string, path string, handler ContextHandler) {
    handle := s.Handle
    s.Mux.HandleFunc(method + " " + path, func(w http.ResponseWriter, req *http.Request) {
        handler(Context{
            Response: w,
            Request:  req,
            Handle: handle,
            Renderer: s.Renderer,
        }) 
    })
}

func (s *Server) Get(path string, handler ContextHandler) {
    s.Route("GET", path, handler)
}

func (s *Server) Post(path string, handler ContextHandler) {
    s.Route("POST", path, handler)
}

func (s *Server) Use(middleware MiddlewareHandler) {
    handle := s.Handle
    s.Handle = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        middleware(Context{
            Response: w,
            Request:  r,
            Handle: handle,
            Renderer: s.Renderer,
        }) 
    })
}

func (s *Server) Group(path string) *Router{
   mux := http.NewServeMux()
   router := &Router{
        BaseUrl: path,
        Handle: mux,
        Mux: mux,
        Renderer: s.Renderer,
        Childs: make([]*Router, 0, 16),
   }
   s.Childs = append(s.Childs, router)
   return router
}

type Router struct {
    BaseUrl     string
    Handle      http.Handler
    Mux         *http.ServeMux
    Childs      []*Router     
    Renderer    *renderer.Renderer
}

func (r *Router) Group(path string) *Router{
    mux := http.NewServeMux()
    router := &Router{
         BaseUrl: r.BaseUrl + path,
         Handle: mux,
         Mux: mux,
         Renderer: r.Renderer,
         Childs: make([]*Router, 0, 16),
    }
    r.Childs = append(r.Childs, router)
    return router
}

func (r *Router) Route(method string, path string, handler ContextHandler) {
    handle := r.Handle
    r.Mux.HandleFunc(method + " " + r.BaseUrl + path, func(w http.ResponseWriter, req *http.Request) {
        handler(Context{
            Response: w,
            Request:  req,
            Handle: handle,
            Renderer: r.Renderer,
        }) 
    })
}

func (r *Router) Get(path string, handler ContextHandler) {
    r.Route("GET", path, handler)
}

func (r *Router) Post(path string, handler ContextHandler) {
    r.Route("POST", path, handler)
}

func (r *Router) Use(middleware MiddlewareHandler) {
    handle := r.Handle
    r.Handle = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        middleware(Context{
            Response: w,
            Request:  req,
            Handle: handle,
            Renderer: r.Renderer,
        }) 
    })
}
