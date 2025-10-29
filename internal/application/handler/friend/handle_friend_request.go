package friend

import (
	"gochat/api/api/friend"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	"gochat/models"

	"github.com/gin-gonic/gin"
)

// HandleFriendRequest 处理好友申请
func (h *FriendHandler) HandleFriendRequest(ctx *gin.Context) {
	req, err := analysis.BindParameter[friend.HandleFriendRequestRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.handleFriendRequestLogic(ctx, &req)
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

func (h *FriendHandler) handleFriendRequestLogic(ctx *gin.Context, req *friend.HandleFriendRequestRequest) (resp *friend.HandleFriendRequestResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	requestID := uint(req.GetRequestId())
	action := req.GetAction()

	// 验证操作类型
	if action != "accept" && action != "reject" {
		return &friend.HandleFriendRequestResponse{
			Success: false,
			Message: "无效的操作类型",
		}, 0, nil
	}

	// 更新申请状态
	var status string
	if action == "accept" {
		status = "accepted" // 已同意
	} else {
		status = "rejected" // 已拒绝
	}

	err = h.dao.UpdateFriendRequestStatus(requestID, status)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 如果同意申请，创建好友关系
	if action == "accept" {
		// 查询申请详情获取申请人ID
		friendRequest, err := h.dao.GetFriendRequestByID(requestID)
		if err != nil {
			return nil, code_msg.ServerError, err
		}
		if friendRequest == nil {
			return &friend.HandleFriendRequestResponse{
				Success: false,
				Message: "好友申请不存在",
			}, 0, nil
		}

		fromUserID := friendRequest.FromUserID

		// 创建双向好友关系
		friend1 := &models.Friend{
			UserID:    userID,
			FriendID:  fromUserID,
			IsBlocked: false, // 正常状态
		}
		friend2 := &models.Friend{
			UserID:    fromUserID,
			FriendID:  userID,
			IsBlocked: false, // 正常状态
		}

		// 使用事务确保两条记录都创建成功
		err = h.dao.CreateFriendsInTransaction(friend1, friend2)
		if err != nil {
			return nil, code_msg.ServerError, err
		}
	}

	message := "已拒绝好友申请"
	if action == "accept" {
		message = "已同意好友申请"
	}

	return &friend.HandleFriendRequestResponse{
		Success: true,
		Message: message,
	}, 0, nil
}
