package main

import (
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-msgdb-server/pkg/app"
	"github.com/thk-im/thk-im-msgdb-server/pkg/handler"
)

func main() {
	configPath := "etc/msg_db_server.yaml"
	config := &conf.Config{}
	if err := conf.LoadConfig(configPath, config); err != nil {
		panic(err)
	}

	appCtx := &app.Context{}
	appCtx.Init(config)
	handler.RegisterMsgDbHandlers(appCtx)

	appCtx.StartServe()
}
