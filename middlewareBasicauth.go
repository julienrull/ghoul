package ghoul

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func NewBasicAuth(config BasicAuthConfig) MiddlewareHandler {
    return func(c Ctx) error {
        auth_header := strings.Split(c.Request().Header.Get("Authorization"), " ")
        if len(auth_header) > 1 && auth_header[0] == "Basic" {
            raw, base64_err := base64.StdEncoding.DecodeString(auth_header[1])
            if base64_err != nil {
                return fmt.Errorf(base64_err.Error())
            }
            creds := strings.Split(string(raw), ":")
            fmt.Println(creds)
            return c.Next() 
        }
        c.Response().Header().Set("WWW-Authenticate", "Basic realm=\"Accès au site\", charset=\"UTF-8\"")
        c.Status(http.StatusUnauthorized)
        return c.Send([]byte(http.StatusText(http.StatusUnauthorized)))
    }
}
