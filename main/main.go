package main

import (
	"github.com/julienrull/ghoul"
)

func main() {
    app := ghoul.New(ghoul.Config{
        Addr: ":3000",
    })
    app.Get("/", func(c ghoul.Ctx) error {
        return c.Json("HELLO")
    })
    app.Use(ghoul.NewBasicAuth(ghoul.BasicAuthConfig{}))
    app.ListenAndServe()
}
