package main

import (
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-msg-db-server/pkg/app"
	"github.com/thk-im/thk-im-msg-db-server/pkg/handler"
)

func main() {
	configPath := "etc/msg_db_server.yaml"
	config := conf.LoadConfig(configPath)

	appCtx := &app.Context{}
	appCtx.Init(config)
	handler.RegisterMsgDbHandlers(appCtx)

	appCtx.StartServe()
}
