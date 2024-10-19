package ghoul

import (
	"net/http/httptest"
	"strings"
	"testing"
)


func TestSimpleRoute(t *testing.T) {
    s := GetServer()
    st := httptest.NewServer(s.Handle)
    defer st.Close()   
    c := NewClient(st.URL)

    res1, _ := c.TestGetQuery("/landing")
    res2, _ := c.TestPostQuery("/landing")

    res1exp := "landing"
    res2exp := "Method Not Allowed\n"

    if res1 != res1exp {
        t.Errorf("expected res to be %s got %s", res1exp, res1)
    }
    if res2 != res2exp {
        t.Errorf("expected res to be %s got %s", res2exp, res2)
    }
}


func TestNestedRoute(t *testing.T) {
    s := GetServer()
    st := httptest.NewServer(s.Handle)
    defer st.Close()   
    c := NewClient(st.URL)

    res1, _ := c.TestGetQuery("/guest/signin")
    res2, _ := c.TestGetQuery("/users/1/home")

    res1exp := "signin"
    res2exp := "signin"

    if res1 != res1exp {
        t.Errorf("expected res to be %s got %s", res1exp, res1)
    }
    if res2 != res2exp {
        t.Errorf("expected res to be %s got %s", res2exp, res2)
    }
}


func TestMiddleware(t *testing.T) {
    s := GetServer()
    st := httptest.NewServer(s.Handle)
    defer st.Close()   
    c := NewClient(st.URL)

    res1, _ := c.TestGetQuery("/users/1/home")
    res2, _ := c.TestPostQuery("/guest/signin")
    res3, _ := c.TestGetQuery("/guest/signin")
    res4, _ := c.TestGetQuery("/users/1/posts/1")

    res1exp := "signin"
    res2exp := "user n°1"
    res3exp := "user n°1"
    res4exp := `<!doctype html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <title>test</title>
    </head>
    <body>
        <h1>Main layout</h1> 
        <main >
            <h1>Post n°1</h1>
        </main>
    </body>
</html>`

    if res1 != res1exp {
        t.Errorf("expected res to be %s got %s", res1exp, res1)
    }
    if res2 != res2exp {
        t.Errorf("expected res to be %s got %s", res2exp, res2)
    }
    if res3 != res3exp {
        t.Errorf("expected res to be %s got %s", res3exp, res3)
    }

    res4expTrimed := strings.Replace(res4exp, " ", "", -1)
    res4expTrimed = strings.Replace(res4expTrimed, "\n", "", -1)
    res4Trimed := strings.Replace(res4, " ", "", -1)
    res4Trimed = strings.Replace(res4Trimed, "\n", "", -1)
    if res4Trimed != res4expTrimed {
        t.Errorf("expected res to be %s got %s", res4expTrimed, res4Trimed)
    }
}
