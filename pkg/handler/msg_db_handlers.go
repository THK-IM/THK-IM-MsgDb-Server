package handler

import (
	"encoding/json"
	"github.com/thk-im/thk-im-base-server/event"
	"github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"github.com/thk-im/thk-im-msgapi-server/pkg/model"
	"github.com/thk-im/thk-im-msgdb-server/pkg/app"
	"github.com/thk-im/thk-im-msgdb-server/pkg/errorx"
)

func RegisterMsgDbHandlers(appCtx *app.Context) {
	appCtx.MsgSaverSubscriber().Sub(func(m map[string]interface{}) error {
		return onMqSaveMsgEventReceived(m, appCtx)
	})
}

func onMqSaveMsgEventReceived(m map[string]interface{}, appCtx *app.Context) error {
	msgJsonStr, okMsg := m[event.SaveMsgEventKey].(string)
	receiversStr, okReceiver := m[event.SaveMsgUsersKey].(string)
	appCtx.Logger().Info(msgJsonStr, okMsg, receiversStr, okReceiver)
	if okMsg && okReceiver {
		message := &dto.Message{}
		err := json.Unmarshal([]byte(msgJsonStr), message)
		if err != nil {
			appCtx.Logger().Error(err)
			return errorx.ErrMessageFormat
		}
		receivers := make([]int64, 0)
		err = json.Unmarshal([]byte(receiversStr), &receivers)
		for _, r := range receivers {
			status := 0
			if r == message.FUid {
				status = model.MsgStatusAcked | model.MsgStatusRead
			}
			userMessage := &model.UserMessage{
				MsgId:      message.MsgId,
				ClientId:   message.CId,
				UserId:     r,
				SessionId:  message.SId,
				FromUserId: message.FUid,
				AtUsers:    message.AtUsers,
				MsgType:    message.Type,
				MsgContent: message.Body,
				ExtData:    message.ExtData,
				ReplyMsgId: message.RMsgId,
				Status:     status,
				CreateTime: message.CTime,
				UpdateTime: message.CTime,
				Deleted:    0,
			}
			err = appCtx.UserMessageModel().InsertUserMessage(userMessage)
			if err != nil {
				appCtx.Logger().Error(err)
				return errorx.ErrMessageFormat
			}
			// 处理原始消息
			if userMessage.ReplyMsgId != nil {
				if message.Type == model.MsgTypeRead && message.RMsgId != nil && r == message.FUid { // 发送已读的人将自己的消息标记为已读
					err = appCtx.UserMessageModel().UpdateUserMessage(r, message.SId, []int64{*message.RMsgId}, model.MsgStatusRead, nil)
					if err != nil {
						return err
					}
				}
				if message.Type == model.MsgTypeRevoke && message.RMsgId != nil {
					err = appCtx.UserMessageModel().DeleteMessages(r, message.SId, []int64{*message.RMsgId}, nil, nil)
					if err != nil {
						return err
					}
				}
				if message.Type == model.MsgTypeReedit && message.RMsgId != nil {
					err = appCtx.UserMessageModel().UpdateUserMessage(r, message.SId, []int64{*message.RMsgId}, model.MsgStatusReedit, &message.Body)
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	} else {
		appCtx.Logger().Error("okReceiver, okMsg:", okReceiver, okMsg)
		return errorx.ErrMessageFormat
	}
}
