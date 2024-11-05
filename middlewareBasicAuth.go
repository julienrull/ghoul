package ghoul

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func NewBasicAuth(config ...BasicAuthConfig) MiddlewareHandler {
    current_config := DefaultBasicAuthConfig
    if len(config) > 0 {
        current_config = config[0]
    }
    return func(c Ctx) error {
        if current_config.Next != nil {
            if current_config.Next(c) {
                return c.Next()
            }
        }
        auth_header := strings.Split(c.Request().Header.Get("Authorization"), " ")
        if len(auth_header) > 1 && auth_header[0] == "Basic" {
            raw, base64_err := base64.StdEncoding.DecodeString(auth_header[1])
            if base64_err != nil {
                return fmt.Errorf(base64_err.Error())
            }
            creds := strings.Split(string(raw), ":")
            if len(creds) > 1 {
                if current_config.Authorizer != nil {
                    if current_config.Authorizer(creds[0], creds[1]) {
                        return c.Next() 
                    }
                }else {
                    value, ok := current_config.Users[creds[0]]
                    if ok {
                        if value == creds[1] {
                            c.Request().WithContext(context.WithValue(c.Request().Context(), current_config.ContextUsername,  creds[0]))
                            c.Request().WithContext(context.WithValue(c.Request().Context(), current_config.ContextPassword,  creds[1]))
                            return c.Next() 
                        }
                    }
                }
            }
        }
        c.Response().Header().Set("WWW-Authenticate", "Basic realm=\""+ current_config.Realm +"\", charset=\"UTF-8\"")
        c.Status(http.StatusUnauthorized)
        if current_config.Unauthorized != nil {
            return current_config.Unauthorized(c)
        }else {
            return c.Send([]byte(http.StatusText(http.StatusUnauthorized)))
        }
    }
}
