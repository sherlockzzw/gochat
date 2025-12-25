package call

import (
	"encoding/json"
	"fmt"
	"gochat/api/api/call"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	globalUtils "gochat/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StartGroupCall 发起群聊通话
func (h *CallHandler) StartGroupCall(ctx *gin.Context) {
	req, err := analysis.BindParameter[call.StartGroupCallRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.startGroupCallLogic(ctx, &req)
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

func (h *CallHandler) startGroupCallLogic(ctx *gin.Context, req *call.StartGroupCallRequest) (resp *call.StartGroupCallResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	groupID := int64(req.GetGroupId())
	callType := req.GetType()
	memberIDs := req.GetMemberIds()

	// 验证参数
	if groupID <= 0 {
		return nil, code_msg.BadRequest, nil
	}
	if callType != models.CallTypeVoice && callType != models.CallTypeVideo {
		return nil, code_msg.BadRequest, nil
	}

	// 验证群组是否存在
	groupDao := dao.NewGroupDao(globalUtils.DB)
	groupInfo, err := groupDao.GetGroupByID(groupID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if groupInfo == nil {
		return nil, code_msg.NotFound, nil
	}

	// 验证用户是否在群组中
	isMember, err := groupDao.IsMemberInGroup(groupID, userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if !isMember {
		return nil, code_msg.NotInGroup, nil
	}

	// 获取要邀请的成员列表
	var targetMemberIDs []int64
	if len(memberIDs) > 0 {
		// 选择了特定成员
		targetMemberIDs = make([]int64, 0, len(memberIDs))
		for _, id := range memberIDs {
			if int64(id) != userID { // 排除自己
				targetMemberIDs = append(targetMemberIDs, int64(id))
			}
		}
		// 验证所有选择的成员都在群组中
		for _, memberID := range targetMemberIDs {
			isMember, err := groupDao.IsMemberInGroup(groupID, memberID)
			if err != nil {
				return nil, code_msg.ServerError, err
			}
			if !isMember {
				return nil, code_msg.UserNotInGroup, nil
			}
		}
	} else {
		// 未选择成员，通知全部成员（排除自己）
		allMembers, err := groupDao.GetGroupMembers(groupID)
		if err != nil {
			return nil, code_msg.ServerError, err
		}
		targetMemberIDs = make([]int64, 0, len(allMembers))
		for _, member := range allMembers {
			if member.UserID != userID {
				targetMemberIDs = append(targetMemberIDs, member.UserID)
			}
		}
	}

	if len(targetMemberIDs) == 0 {
		return nil, code_msg.BadRequest, fmt.Errorf("没有可邀请的成员")
	}

	// 检查WebSocket Hub
	wsHub := getWebSocketHub()
	if wsHub == nil {
		return nil, code_msg.ServerError, fmt.Errorf("WebSocket Hub未初始化")
	}

	// 使用事务创建通话房间和参与者
	var roomID int64
	var roomToken string

	err = h.dao.GetDB().Transaction(func(tx *gorm.DB) error {
		callDao := dao.NewCallDao(tx)

		// 生成房间ID和Token
		roomID = time.Now().UnixNano() / 1e6 // 使用毫秒时间戳作为房间ID
		roomToken = fmt.Sprintf("group_%d_%d_%d", groupID, userID, roomID)

		// 创建通话房间
		room := &models.CallRoom{
			ID:        roomID,
			Type:      callType,
			CallType:  models.CallSceneGroup,
			CreatorID: userID,
			GroupID:   groupID,
			RoomToken: roomToken,
			Status:    models.RoomStatusCalling,
		}

		if err := callDao.CreateRoom(room); err != nil {
			return err
		}

		// 创建参与者（包括创建者）
		participants := make([]*models.CallParticipant, 0, len(targetMemberIDs)+1)
		
		// 添加创建者（状态为joined，因为创建者自动加入）
		participants = append(participants, &models.CallParticipant{
			RoomID:   roomID,
			UserID:   userID,
			Status:   models.ParticipantStatusJoined,
			JoinedAt: time.Now().Unix(),
		})

		// 添加被邀请的成员（状态为invited）
		for _, memberID := range targetMemberIDs {
			participants = append(participants, &models.CallParticipant{
				RoomID: roomID,
				UserID: memberID,
				Status: models.ParticipantStatusInvited,
			})
		}

		// 批量创建参与者
		if err := callDao.CreateParticipants(participants); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 通过WebSocket发送通话邀请给所有被邀请的成员
	inviteMessage := map[string]interface{}{
		"type":      "call_invite",
		"room_id":   roomID,
		"room_token": roomToken,
		"call_type": callType,
		"group_id":  groupID,
		"creator_id": userID,
		"group_name": groupInfo.Name,
	}

	inviteData, err := json.Marshal(inviteMessage)
	if err != nil {
		fmt.Printf("Failed to marshal invite message: %v\n", err)
	} else {
		// 发送给所有被邀请的成员
		for _, memberID := range targetMemberIDs {
			wsHub.SendToUser(memberID, inviteData)
		}
	}

	// 启动通话超时定时器（群聊时，超时时间对所有成员生效）
	if h.timeoutManager != nil {
		h.timeoutManager.StartTimeout(roomID)
	}

	resp = &call.StartGroupCallResponse{
		RoomId:      uint64(roomID),
		RoomToken:   roomToken,
		InvitedCount: int32(len(targetMemberIDs)),
		Success:     true,
	}

	return resp, 0, nil
}

