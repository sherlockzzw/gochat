package friend

import (
	"fmt"
	"gochat/api/api/friend"
	"gochat/internal/application/handler/common"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"time"

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
	userID, errCode, err := common.GetUserIDFromContext(ctx)
	if errCode != 0 {
		return nil, errCode, err
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

	// 检查目标用户是否存在
	targetUser, err := h.userDao.GetUserByID(int64(req.GetFriendId()))
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if targetUser == nil {
		return nil, code_msg.BadRequest, nil
	}

	// 检查目标用户的隐私设置：是否需要验证
	// TODO: 从 targetUser.PrivacySettings JSON 中读取 add_friend_need_verify 字段
	// 这里暂时使用默认值（需要验证）

	// 检查是否已有待处理的申请
	existingRequest, err := h.dao.GetFriendRequest(userID, int64(req.GetFriendId()))
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if existingRequest != nil {
		return nil, code_msg.FriendRequestExists, nil
	}

	// 检查目标用户是否允许直接添加（从隐私设置中读取，默认需要验证）
	needVerify := true // 默认需要验证
	// TODO: 从 targetUser.PrivacySettings JSON 中读取 add_friend_need_verify 字段
	// 这里暂时使用默认值，后续可以从隐私设置中读取

	if !needVerify {
		// 直接添加为好友，不需要验证
		now := time.Now().Unix()
		friend1 := &models.Friend{
			UserID:    userID,
			FriendID:  int64(req.GetFriendId()),
			Remark:    "",
			GroupName: req.GetGroupName(),
			IsBlocked: false,
			CreatedAt: now,
			UpdatedAt: now,
		}
		friend2 := &models.Friend{
			UserID:    int64(req.GetFriendId()),
			FriendID:  userID,
			Remark:    "",
			GroupName: "", // 对方的分组由对方自己设置
			IsBlocked: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err = h.dao.CreateFriendsInTransaction(friend1, friend2)
		if err != nil {
			return nil, code_msg.ServerError, err
		}

		// 创建通知：通知对方已添加为好友
		go func() {
			sender, _ := h.userDao.GetUserByID(userID)
			senderName := "用户"
			if sender != nil {
				senderName = sender.Name
			}
			_ = common.CreateSystemNotification(
				int64(req.GetFriendId()),
				"新好友",
				fmt.Sprintf("%s 添加你为好友", senderName),
				0,
			)
		}()

		return &friend.AddFriendResponse{
			Success: true,
			Message: "添加成功",
		}, 0, nil
	}

	// 需要验证：创建好友申请
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

	// 创建通知：给接收者发送好友申请通知
	sender, _ := h.userDao.GetUserByID(userID)
	senderName := "用户"
	if sender != nil {
		senderName = sender.Name
	}

	title := "好友申请"
	content := fmt.Sprintf("%s 申请添加你为好友", senderName)
	if req.GetMessage() != "" {
		content = fmt.Sprintf("%s 申请添加你为好友：%s", senderName, req.GetMessage())
	}

	// 异步创建通知
	go func() {
		if err := common.CreateSystemNotification(
			int64(req.GetFriendId()),
			title,
			content,
			friendRequest.ID,
		); err != nil {
			fmt.Printf("Failed to create friend request notification: %v\n", err)
		}
	}()

	return &friend.AddFriendResponse{
		Success: true,
	}, 0, nil
}
