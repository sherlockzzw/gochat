package friend

import (
	"fmt"
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
		return nil, code_msg.NotFriend, nil
	}

	// 删除双向好友关系（互移列表）
	err = h.dao.DeleteFriend(userID, friendID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 同时删除对方的好友关系
	err = h.dao.DeleteFriend(friendID, userID)
	if err != nil {
		// 记录错误但不影响主流程
		fmt.Printf("Failed to delete reverse friend relation: %v\n", err)
	}

	// 注意：历史消息保留在数据库中，不会被删除
	// 中断新消息接收：通过删除好友关系，后续发送消息时会检查好友关系，从而阻止新消息

	return &friend.DeleteFriendResponse{
		Success: true,
		Message: "删除成功，历史记录已保留",
	}, 0, nil
}
