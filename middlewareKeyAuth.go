package ghoul

import (
	"context"
	"fmt"
	"strings"
)

// TODO: Create Errors


func NewKeyAuth(config ...KeyAuthConfig) MiddlewareHandler {
    current_config := DefaultKeyAuthConfig
    if len(config) > 0 {
        current_config = config[0]
    }
    return func(c Ctx) error {
        key := ""
        if current_config.Next != nil {
            if current_config.Next(c) {
                return c.Next()
            }
        }
        if current_config.CustomKeyLookup != nil {
            val, lookup_error := current_config.CustomKeyLookup(c)
            if lookup_error != nil {
                return lookup_error
            }
            key = val
        }else {
            keyLookup := strings.Split(current_config.KeyLookup, ":")
            scheme := ""
            if len(keyLookup) != 2 {
                return current_config.ErrorHandler(c, fmt.Errorf("Invalid KeyLookup format."))
            }
            switch keyLookup[0] {
            case "header": 
                scheme = c.Request().Header.Get(keyLookup[1])
            case "cookie": 
                cookie, _ := c.Request().Cookie(keyLookup[1])
                scheme = cookie.Value
            case "query":
                scheme = c.Request().URL.Query().Get(keyLookup[1])
            case "form": 
                scheme = c.Request().FormValue(keyLookup[1])
            default:
                return current_config.ErrorHandler(c, fmt.Errorf("Malformed Key."))
            }
            schemeKey := strings.Split(scheme, " ")
            if len(schemeKey) == 2 {
                if schemeKey[0] != current_config.AuthScheme {
                    return current_config.ErrorHandler(c, fmt.Errorf("Malformed Key."))
                }
                key = schemeKey[1]
            }else if len(schemeKey) == 1 {
                if current_config.AuthScheme != "" {
                    return current_config.ErrorHandler(c, fmt.Errorf("Malformed Key."))
                }
                key = schemeKey[0]
            }else {
                return fmt.Errorf("Malformed Key.")
            }
        }
        valid, validator_err := current_config.Validator(c, key)
        if !valid || validator_err != nil {
            return current_config.ErrorHandler(c, validator_err)
        }
        c.Request().WithContext(context.WithValue(c.Request().Context(), current_config.ContextKey, key))
        return current_config.SuccessHandler(c)
    }
}
