package app

import (
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-base-server/server"
	"github.com/thk-im/thk-im-msgapi-server/pkg/model"
	"github.com/thk-im/thk-im-msgdb-server/pkg/loader"
)

type Context struct {
	*server.Context
}

func (c *Context) UserMessageModel() model.UserMessageModel {
	return c.Context.ModelMap["user_message"].(model.UserMessageModel)
}

func (c *Context) Init(config *conf.Config) {
	c.Context = &server.Context{}
	c.Context.Init(config)
	c.Context.ModelMap = loader.LoadModels(c.Config().Models, c.Database(), c.Logger(), c.SnowflakeNode())

}
