package loader

import (
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-base-server/snowflake"
	"github.com/thk-im/thk-im-msgapi-server/pkg/model"
	"gorm.io/gorm"
)

func LoadModels(modeConfigs []conf.Model, database *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node) map[string]interface{} {
	modelMap := make(map[string]interface{})
	for _, ms := range modeConfigs {
		var m interface{}
		if ms.Name == "user_message" {
			m = model.NewUserMessageModel(database, logger, snowflakeNode, ms.Shards)
		}
		modelMap[ms.Name] = m
	}
	return modelMap
}
