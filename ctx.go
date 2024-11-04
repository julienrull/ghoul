package ghoul

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func NewCtx(
    r           *http.Request,
    w           http.ResponseWriter,
    handler      http.Handler,
    renderer    *Renderer,
) Ctx {
    return Ctx{
        r,
        w,
        handler,
        renderer,
    }
}
func (c Ctx) Exec(handler ContextHandler){
    err := handler(c)
    HandleError(c, err)
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
    return fmt.Errorf("can't process nil next handler")
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

func (c *Ctx) Json(v any) error{
    c.Response().Header().Set("Content-type", "application/json")
    res, json_err := json.Marshal(v)
    if json_err != nil {
        return json_err
    }
    _, write_err := c.Write(res)
    if write_err != nil {
        return write_err
    } 
    return nil
}

func (c *Ctx) Status(code int) {
    c.Response().WriteHeader(code)
}

func (c *Ctx) Redirect(path string, status int) error {
    http.Redirect(c.w, c.r, path, status)
    return nil
}

func (c *Ctx) Render(tmplName string, data map[string]any, layouts ...string) error {
    c.renderer.Render(c.w, tmplName, data, layouts...) 
    return nil
}
