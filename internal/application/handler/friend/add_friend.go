package friend

import (
	"gochat/api/api/friend"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	"gochat/internal/infrastructure/models"

	"github.com/gin-gonic/gin"
)

// AddFriend 添加好友
func (h *FriendHandler) AddFriend(ctx *gin.Context) {
	req, err := analysis.BindParameter[friend.AddFriendRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.addFriendLogic(ctx, &req)
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

func (h *FriendHandler) addFriendLogic(ctx *gin.Context, req *friend.AddFriendRequest) (resp *friend.AddFriendResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 检查不能添加自己为好友
	if userID == int64(req.GetFriendId()) {
		return nil, code_msg.CannotAddSelf, nil
	}

	// 检查是否已经是好友
	existingFriend, err := h.dao.CheckIsFriend(userID, int64(req.GetFriendId()))
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if existingFriend != nil {
		return nil, code_msg.AlreadyFriend, nil
	}

	// 检查是否已有待处理的申请
	existingRequest, err := h.dao.GetFriendRequest(userID, int64(req.GetFriendId()))
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if existingRequest != nil {
		return nil, code_msg.FriendRequestExists, nil
	}

	// 创建好友申请
	friendRequest := &models.FriendRequest{
		FromUserID: userID,
		ToUserID:   int64(req.GetFriendId()),
		Message:    req.GetMessage(),
		Status:     "pending", // 待处理
	}

	err = h.dao.CreateFriendRequest(friendRequest)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return &friend.AddFriendResponse{
		Success: true,
	}, 0, nil
}
