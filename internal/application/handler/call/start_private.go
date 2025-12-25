package call

import (
	"encoding/json"
	"fmt"
	"gochat/api/api/call"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/infrastructure/models"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	globalUtils "gochat/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StartPrivateCall 发起私聊通话
func (h *CallHandler) StartPrivateCall(ctx *gin.Context) {
	req, err := analysis.BindParameter[call.StartPrivateCallRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.startPrivateCallLogic(ctx, &req)
	if code != 0 {
		h.response.JsonErrorFixation(ctx, code)
		return
	}
	if err != nil {
		h.response.JsonError(ctx, err, err.Error())
		return
	}

	h.response.JsonSuccess(ctx, resp)
}

func (h *CallHandler) startPrivateCallLogic(ctx *gin.Context, req *call.StartPrivateCallRequest) (resp *call.StartPrivateCallResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	receiverID := int64(req.GetReceiverId())
	callType := req.GetType()

	// 验证参数
	if receiverID <= 0 {
		return nil, code_msg.BadRequest, nil
	}
	if callType != models.CallTypeVoice && callType != models.CallTypeVideo {
		return nil, code_msg.BadRequest, nil
	}

	// 验证是否为好友关系（私聊需要是好友）
	friendDao := dao.NewFriendDao(globalUtils.DB)
	friend, err := friendDao.CheckIsFriend(userID, receiverID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if friend == nil {
		return nil, code_msg.NotFriend, nil
	}

	// 检查接收者是否在线
	wsHub := getWebSocketHub()
	if wsHub == nil {
		return nil, code_msg.ServerError, fmt.Errorf("WebSocket Hub未初始化")
	}

	// 使用事务创建通话房间
	var roomID int64
	var roomToken string

	err = h.dao.GetDB().Transaction(func(tx *gorm.DB) error {
		callDao := dao.NewCallDao(tx)

		// 生成房间令牌
		roomToken = callDao.GenerateRoomToken()

		// 创建通话房间
		now := time.Now().Unix()
		room := &models.CallRoom{
			Type:      callType,
			CallType:  models.CallScenePrivate,
			CreatorID: userID,
			GroupID:   0,
			RoomToken: roomToken,
			Status:    models.RoomStatusCalling,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := callDao.CreateRoom(room); err != nil {
			return err
		}
		roomID = room.ID

		// 创建发起者参与者记录
		creatorParticipant := &models.CallParticipant{
			RoomID:    roomID,
			UserID:    userID,
			Status:    models.ParticipantStatusJoined, // 发起者直接加入
			JoinedAt:  now,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := callDao.CreateParticipant(creatorParticipant); err != nil {
			return err
		}

		// 创建接收者参与者记录
		receiverParticipant := &models.CallParticipant{
			RoomID:    roomID,
			UserID:    receiverID,
			Status:    models.ParticipantStatusInvited,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := callDao.CreateParticipant(receiverParticipant); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 通过WebSocket发送通话邀请
	inviteMessage := map[string]interface{}{
		"type":        "call_invite",
		"room_id":     roomID,
		"room_token":  roomToken,
		"caller_id":   userID,
		"call_type":   "private",
		"media_type":  callType,
		"receiver_id": receiverID,
	}

	// 获取发送者信息（用于显示）
	userDao := dao.NewUserDao(globalUtils.DB)
	caller, err := userDao.GetUserByID(userID)
	if err == nil && caller != nil {
		inviteMessage["caller_name"] = caller.Name
		inviteMessage["caller_avatar"] = caller.Avatar
	}

	inviteBytes, _ := json.Marshal(inviteMessage)
	wsHub.SendToUser(receiverID, inviteBytes)

	// 更新接收者状态为响铃中
	callDao := dao.NewCallDao(globalUtils.DB)
	err = callDao.UpdateParticipantStatus(roomID, receiverID, models.ParticipantStatusRinging)
	if err != nil {
		// 记录错误但不影响整体流程（房间已创建，消息已发送）
		fmt.Printf("Failed to update participant status: %v\n", err)
	}

	err = callDao.UpdateRoomStatus(roomID, models.RoomStatusRinging)
	if err != nil {
		// 记录错误但不影响整体流程
		fmt.Printf("Failed to update room status: %v\n", err)
	}

	// 启动通话超时定时器
	if timeoutMgr := globalUtils.GetCallTimeoutManager(); timeoutMgr != nil {
		if manager, ok := timeoutMgr.(*websocket.CallTimeoutManager); ok {
			manager.StartTimeout(roomID)
		}
	}

	resp = &call.StartPrivateCallResponse{
		RoomId:    uint64(roomID),
		RoomToken: roomToken,
		Success:   true,
	}

	return resp, 0, nil
}
