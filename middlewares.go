package ghoul

import "net/http"

type Middleware = func(http.Handler) http.Handler
type MiddlewareHandler = ContextHandler


