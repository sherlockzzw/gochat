package friend

import (
	"fmt"
	"gochat/api/api/friend"
	"gochat/internal/application/handler/common"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	globalUtils "gochat/utils"
	"time"

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

	requestID := int64(req.GetRequestId())
	action := req.GetAction()

	// 验证操作类型
	if action != "accept" && action != "reject" {
		return nil, code_msg.InvalidOperation, nil
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

	// 查询申请详情获取申请人ID（用于创建好友关系和通知）
	var fromUserID int64
	friendRequest, err := h.dao.GetFriendRequestByID(requestID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if friendRequest == nil {
		return nil, code_msg.FriendRequestNotExists, nil
	}
	fromUserID = friendRequest.FromUserID

	// 如果同意申请，创建好友关系
	if action == "accept" {
		now := time.Now().Unix()
		// 创建双向好友关系
		friend1 := &models.Friend{
			UserID:    userID,
			FriendID:  fromUserID,
			Remark:    "",
			GroupName: "", // 接收者可以后续设置分组
			IsBlocked: false,
			CreatedAt: now,
			UpdatedAt: now,
		}
		friend2 := &models.Friend{
			UserID:    fromUserID,
			FriendID:  userID,
			Remark:    "",
			GroupName: "", // 申请者可以后续设置分组
			IsBlocked: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		// 使用事务确保两条记录都创建成功
		err = h.dao.CreateFriendsInTransaction(friend1, friend2)
		if err != nil {
			return nil, code_msg.ServerError, err
		}
	}

	// 创建通知：给申请者发送好友申请处理结果通知
	userDao := dao.NewUserDao(globalUtils.DB)
	handler, _ := userDao.GetUserByID(userID)
	handlerName := "用户"
	if handler != nil {
		handlerName = handler.Name
	}

	var title, content string
	if action == "accept" {
		title = "好友申请已通过"
		content = fmt.Sprintf("%s 已同意你的好友申请", handlerName)
	} else {
		title = "好友申请被拒绝"
		content = fmt.Sprintf("%s 拒绝了你的好友申请", handlerName)
	}

	// 异步创建通知
	go func() {
		if err := common.CreateSystemNotification(
			fromUserID,
			title,
			content,
			requestID,
		); err != nil {
			fmt.Printf("Failed to create friend request handle notification: %v\n", err)
		}
	}()

	return &friend.HandleFriendRequestResponse{
		Success: true,
	}, 0, nil
}
