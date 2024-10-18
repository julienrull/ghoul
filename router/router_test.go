package main 

import (
	"net/http/httptest"
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
    res2, _ := c.TestGetQuery("/user/home")

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

    res1, _ := c.TestGetQuery("/user/home")
    res2, _ := c.TestPostQuery("/guest/signin")
    res3, _ := c.TestGetQuery("/guest/signin")

    res1exp := "signin"
    res2exp := "home"
    res3exp := "home"

    if res1 != res1exp {
        t.Errorf("expected res to be %s got %s", res1exp, res1)
    }
    if res2 != res2exp {
        t.Errorf("expected res to be %s got %s", res2exp, res2)
    }
    if res3 != res3exp {
        t.Errorf("expected res to be %s got %s", res3exp, res3)
    }
}
