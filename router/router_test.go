package main 

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

    res1, _ := c.TestGetQuery("/simpleroute")
    trimedRes1 := strings.TrimSpace(res1)
    res2, _ := c.TestPostQuery("/simpleroute")

    res1exp := "<h1>simplerouteget</h1>"
    res2exp := "simpleroutepost"
    if trimedRes1 != res1exp {
        t.Errorf("expected res to be %s got %s", res1exp, trimedRes1)
    }
    if res2 != res2exp {
        t.Errorf("expected res to be %s got %s", res2exp, res2)
    }
}


func TestNestedRoute(t *testing.T) {
    s := GetServer()
    defer s.Server.Close()   
    st := httptest.NewServer(s.Handle)
    defer st.Close()   
    c := NewClient(st.URL)
    res1, _ := c.TestGetQuery("/nestedRoute/nested")
    res2, _ := c.TestPostQuery("/nestedRoute/nested")
    res3, _ := c.TestGetQuery("/nestedRoute/subnestedroute/subnested")
    res4, _ := c.TestPostQuery("/nestedRoute/subnestedroute/subnested")
    res1exp := "nestedget"
    res2exp := "nestedpost"
    res3exp := "subnestedget"
    res4exp := "subnestedpost"
    if res1 != res1exp {
        t.Errorf("expected res to be %s got %s", res1exp, res1)
    }
    if res2 != res2exp {
        t.Errorf("expected res to be %s got %s", res2exp, res2)
    }
    if res3 != res3exp {
        t.Errorf("expected res to be %s got %s", res3exp, res3)
    }
    if res4 != res4exp {
        t.Errorf("expected res to be %s got %s", res4exp, res4)
    }
}


func TestSimpleMiddlewares(t *testing.T) {
    s := GetServer()
    defer s.Server.Close()   
    st := httptest.NewServer(s.Handle)
    defer st.Close()   
    c := NewClient(st.URL)
    res1, _ := c.TestGetQuery("/simplemiddleware")
    res2, _ := c.TestPostQuery("/simplemiddleware")
    res1exp := "simplemiddleware"
    res2exp := "simplemiddleware"
    if res1 != res1exp {
        t.Errorf("expected res to be %s got %s", res1exp, res1)
    }
    if res2 != res2exp {
        t.Errorf("expected res to be %s got %s", res2exp, res2)
    }
}

func TestSubMiddlewares(t *testing.T) {
    s := GetServer()
    defer s.Server.Close()   
    st := httptest.NewServer(s.Handle)
    defer st.Close()   
    c := NewClient(st.URL)
    res1, _ := c.TestGetQuery("/firstmiddleware/secondmiddleware")
    res2, _ := c.TestPostQuery("/firstmiddleware/secondmiddleware")
    res1exp := "submiddleware"
    res2exp := "submiddleware"
    if res1 != res1exp {
        t.Errorf("expected res to be %s got %s", res1exp, res1)
    }
    if res2 != res2exp {
        t.Errorf("expected res to be %s got %s", res2exp, res2)
    }
}

func TestNext(t *testing.T) {
    s := GetServer()
    defer s.Server.Close()   
    st := httptest.NewServer(s.Handle)
    defer st.Close()   
    c := NewClient(st.URL)
    res2, _ := c.TestGetQuery("/nextroute/subnextroute")
    res2exp := "subnextroute"
    if res2 != res2exp {
        t.Errorf("expected res to be %s got %s", res2exp, res2)
    }
}

func TestRedirection(t *testing.T) {
    s := GetServer(); defer s.Server.Close()   
    st := httptest.NewServer(s.Handle)
    defer st.Close()   
    c := NewClient(st.URL)
    res1, _ := c.TestGetQuery("/redirectroute")
    res1exp := "redirecttarget"
    if res1 != res1exp {
        t.Errorf("expected res to be %s got %s", res1exp, res1)
    }
}
