package main

import "github.com/julienrull/ghoul/v1"


func main() {
    app := ghoul.New()
    app.Get("/hello", func(c ghoul.Ctx) error {
        c.Response.Write([]byte("Hello"))
        return nil
    })
    app.ListenAndServe()
}
