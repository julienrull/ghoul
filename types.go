package ghoul

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type Config struct {
    Addr                            string
    ReadTimeout                     time.Duration
    WriteTimeout                    time.Duration
    IdleTimeout                     time.Duration
    ReadHeaderTimeout               time.Duration
    MaxHeaderBytes                  int
    Handler                         http.Handler
    DisableGeneralOptionsHandler    bool
    TLSConfig                       *tls.Config
    TLSNextProto                    map[string]func(*http.Server, *tls.Conn, http.Handler)
    ConnState                       func(net.Conn, http.ConnState)
    ErrorLog                        *log.Logger
    BaseContext                     func(net.Listener) context.Context
    ConnContext                     func(ctx context.Context, c net.Conn) context.Context
}
type Router struct {
    Server      *http.Server
    signalOut   chan os.Signal
    Root        *Router 
    isRoot      bool
    BaseUrl     string
    Handle      http.Handler
    Mux         *http.ServeMux
    Renderer    *Renderer
}
type APIError struct {
    StatusCode  int `json:"statusCode"`
    Msg         any `json:"msg"`
}
type Renderer struct {
    Folder string
    Ext    string 
}
type Ctx struct {
    r           *http.Request
    w           http.ResponseWriter
    handler      http.Handler
    renderer    *Renderer
}
type BasicAuthConfig struct {
    OnSuccess   MiddlewareHandler
    OnFailure   MiddlewareHandler
}

type ContextHandler = func(Ctx) error
type ErrorHandler = func (Ctx, error)
type Middleware = func(http.Handler) http.Handler
type MiddlewareHandler = func(Ctx) error

