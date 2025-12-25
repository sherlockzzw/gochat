package friend

import (
	"gochat/api/api/friend"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// DeleteFriend 删除好友
func (h *FriendHandler) DeleteFriend(ctx *gin.Context) {
	req, err := analysis.BindParameter[friend.DeleteFriendRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.deleteFriendLogic(ctx, &req)
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

func (h *FriendHandler) deleteFriendLogic(ctx *gin.Context, req *friend.DeleteFriendRequest) (resp *friend.DeleteFriendResponse, errCode code_msg.BusinessCode, err error) {
	// 获取当前用户ID（从token中解析）
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	friendID := int64(req.GetFriendId())

	// 检查是否为好友
	existingFriend, err := h.dao.CheckIsFriend(userID, friendID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if existingFriend == nil {
		return &friend.DeleteFriendResponse{
			Success: false,
			Message: "不是好友关系",
		}, 0, nil
	}

	// 删除好友关系
	err = h.dao.DeleteFriend(userID, friendID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return &friend.DeleteFriendResponse{
		Success: true,
		Message: "删除好友成功",
	}, 0, nil
}
