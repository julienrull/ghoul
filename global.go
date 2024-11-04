package ghoul

import (
	"context"
	"log"
	"net"
	"net/http"
)

var defaultConfiguration = Config{
    Addr:                           ":http",
    Handler:                        http.DefaultServeMux,
    DisableGeneralOptionsHandler:   false,
    TLSConfig:                      nil,
    ReadTimeout:                    0,
    ReadHeaderTimeout:              0,
    WriteTimeout:                   0,
    IdleTimeout:                    0,
    MaxHeaderBytes:                 http.DefaultMaxHeaderBytes,
    TLSNextProto:                   nil,
    ConnState:                      nil,
    ErrorLog:                       log.Default(),
    BaseContext:                    func(l net.Listener) context.Context {return context.Background()} ,
    ConnContext:                    nil,
}
