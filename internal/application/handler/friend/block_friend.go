package friend

import (
	"gochat/api/api/friend"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// BlockFriend 屏蔽/取消屏蔽好友
func (h *FriendHandler) BlockFriend(ctx *gin.Context) {
	req, err := analysis.BindParameter[friend.BlockFriendRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.blockFriendLogic(ctx, &req)
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

func (h *FriendHandler) blockFriendLogic(ctx *gin.Context, req *friend.BlockFriendRequest) (resp *friend.BlockFriendResponse, errCode code_msg.BusinessCode, err error) {
	// 获取当前用户ID（从token中解析）
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	friendID := int64(req.GetFriendId())
	isBlock := req.GetIsBlock()

	// 检查是否为好友
	existingFriend, err := h.dao.CheckIsFriend(userID, friendID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if existingFriend == nil {
		return nil, code_msg.NotFriend, nil
	}

	// 屏蔽/取消屏蔽好友
	err = h.dao.BlockFriend(userID, friendID, isBlock)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return &friend.BlockFriendResponse{
		Success: true,
	}, 0, nil
}
