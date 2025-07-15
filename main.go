package main

import (
	"butterfly.orx.me/core"
	"butterfly.orx.me/core/app"
	"github.com/orvice/simpleproxy/internal/conf"
	"github.com/orvice/simpleproxy/internal/handler"
)

func main() {
	app := core.New(&app.Config{
		Config:  conf.Conf,
		Service: "simpleproxy",
		Router:  handler.Router,
	})
	app.Run()
}
