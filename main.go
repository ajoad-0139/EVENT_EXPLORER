package main

import (
	_ "event-explorer/routers"
	"event-explorer/secrets"
	"os"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	// must load secrets
	err := secrets.MustLoad()
	if err != nil {
		logs.Error("Failed loading secrets", err)
		os.Exit(1)
	}

	//start the server
	beego.Run()
}
