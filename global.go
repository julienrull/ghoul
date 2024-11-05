package ghoul

import (
	"context"
	"fmt"
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

var DefaultBasicAuthConfig  = BasicAuthConfig {
    Next:           nil,
    Users:           map[string]string{},
    Realm:           "Restricted",
    Authorizer:      nil,
    Unauthorized:    nil,
    ContextUsername: "username",
    ContextPassword: "password",
}

var DefaultKeyAuthConfig  = KeyAuthConfig {
    Next:           nil,
    SuccessHandler: func(c Ctx) error{
	    return c.Next()
	},
    ErrorHandler:   func(c Ctx, err error) error {
		if err != nil {
            c.Status(http.StatusUnauthorized)
			return c.Send([]byte("Invalid or expired API Key"))
		}
		return err
		//if err == ErrMissingOrMalformedAPIKey {
		//	return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
		//}
		//return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired API Key")
    },
  	KeyLookup:      "header:Authorization",
    CustomKeyLookup: nil,
	AuthScheme:     "Bearer",  
    Validator:       func(c Ctx, key string) (bool, error) {
	    return false, fmt.Errorf("Invalid Key.")
    },
    ContextKey:     "token",
}

