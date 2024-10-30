package ghoul

import (
	"net/http"
)

type ContextHandler = func(Ctx) error

type Ctx struct {
    r           *http.Request
    w           http.ResponseWriter
    handler      http.Handler
    renderer    *Renderer
}

func (c *Ctx) Request() *http.Request {
    return c.r
}
func (c *Ctx) Response() http.ResponseWriter {
    return c.w
}

func (c *Ctx) Next() error {
    if c.handler != nil {
        c.handler.ServeHTTP(c.w, c.r)
        return nil
    }
    return NewCtxErrorNext()
}

func (c *Ctx) Flush() bool {
    f, ok := c.Response().(http.Flusher)
    if ok {
        f.Flush()
    }
    return ok
}

func (c *Ctx) Write(data []byte) (int, error) {
    return c.Response().Write(data)
}

func (c *Ctx) Send(data []byte) error {
    _, err := c.Response().Write(data)
    c.Flush()    
    return err
}

func (c *Ctx) Status(code int) {
    c.Response().WriteHeader(code)
}

func (c *Ctx) Redirect(path string, status int) {
    http.Redirect(c.w, c.r, path, status)
}

func (c *Ctx) Render(tmplName string, data map[string]any, layouts ...string) {
    c.renderer.Render(c.w, tmplName, data, layouts...) 
}
