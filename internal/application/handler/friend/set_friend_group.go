package friend

import (
	"gochat/api/api/friend"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// SetFriendGroup 设置好友分组
func (h *FriendHandler) SetFriendGroup(ctx *gin.Context) {
	req, err := analysis.BindParameter[friend.SetFriendGroupRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.setFriendGroupLogic(ctx, &req)
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

func (h *FriendHandler) setFriendGroupLogic(ctx *gin.Context, req *friend.SetFriendGroupRequest) (resp *friend.SetFriendGroupResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	friendID := int64(req.GetFriendId())
	groupName := req.GetGroupName()

	// 验证是否为好友
	existingFriend, err := h.dao.CheckIsFriend(userID, friendID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if existingFriend == nil {
		return nil, code_msg.NotFriend, nil
	}

	// 设置分组（空字符串表示移除分组）
	err = h.dao.SetFriendGroup(userID, friendID, groupName)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return &friend.SetFriendGroupResponse{
		Success: true,
		Message: "设置成功",
	}, 0, nil
}

