package friend

import (
	"gochat/api/api/friend"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"

	"github.com/gin-gonic/gin"
)

// CheckFriend 检查好友关系
func (h *FriendHandler) CheckFriend(ctx *gin.Context) {
	req, err := analysis.BindParameter[friend.CheckFriendRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.checkFriendLogic(ctx, &req)
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

func (h *FriendHandler) checkFriendLogic(ctx *gin.Context, req *friend.CheckFriendRequest) (resp *friend.CheckFriendResponse, errCode code_msg.BusinessCode, err error) {
	// 获取当前用户ID（从token中解析）
	userID := uint(1) // TODO: 从JWT token中获取真实用户ID

	friendID := uint(req.GetFriendId())

	// 检查好友关系
	existingFriend, err := h.dao.CheckIsFriend(userID, friendID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	isFriend := existingFriend != nil

	return &friend.CheckFriendResponse{
		IsFriend:  isFriend,
		IsBlocked: false, // TODO: 需要查询屏蔽状态
		Remark:    "",    // TODO: 需要查询备注
	}, 0, nil
}
