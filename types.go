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

// STRUCTS
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

// CONFIGS
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
type BasicAuthConfig struct {
    Next            func(Ctx) bool
    Users           map[string]string
    Realm           string 
    Authorizer      func(string, string) bool
    Unauthorized    MiddlewareHandler
    ContextUsername interface{}
    ContextPassword interface{}
}
type KeyAuthConfig struct {
    Next            func(Ctx) bool
    SuccessHandler  MiddlewareHandler
    ErrorHandler    ErrorHandler
    KeyLookup       string
    CustomKeyLookup func(Ctx) (string, error)
    AuthScheme      string
    Validator       func(Ctx, string) (bool, error)
    ContextKey      interface{}
}

// FUNCTIONS TYPES
type ContextHandler = func(Ctx) error
type ErrorHandler = func (Ctx, error) error
type Middleware = func(http.Handler) http.Handler
type MiddlewareHandler = func(Ctx) error

