package ghoul

import "net/http"


type ContextHandler = func(Ctx) error

type Ctx struct {
    Request     *http.Request
    Response    http.ResponseWriter
    Status      int
    Handle      http.Handler
    Renderer    *Renderer
}

func (c *Ctx) Next() {
    if c.Handle != nil {
        c.Handle.ServeHTTP(c.Response, c.Request)
    }
}

func (c *Ctx) Redirect(path string, status int) {
    http.Redirect(c.Response, c.Request, path, status)
}

func (c *Ctx) Render(tmplName string, data map[string]any, layouts ...string) {
    c.Renderer.Render(c.Response, tmplName, data, layouts...) 
}
