package ghoul

import (
	"bytes"
	"strings"
	"testing"
)


func TestRender(t *testing.T) {
    r := NewRenderer("./views/", ".html")  
    buf := bytes.NewBufferString("")
    r.Render(buf, "render", map[string]any{"bodyData": "Hello Test"}, "layouts/rendermain")
exp := strings.ReplaceAll(`
<!doctype html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <title>test</title>
    </head>
    <body>
        <h1>Main layout</h1> 
        <main >
            <h1>body</h1>
            Hello Test
        </main>
    </body>
</html>`, " ", "")
    exp = strings.ReplaceAll(exp, "\n", "")
    exp = strings.ReplaceAll(exp, "\t", "")

    res := strings.ReplaceAll(buf.String(), " ", "")
    res = strings.ReplaceAll(res, "\n", "")
    res = strings.ReplaceAll(res, "\t", "")
    if res != exp {
        t.Errorf("expected res to be %s got %s", exp, res)
    }
}
