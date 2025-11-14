package friend

import (
	"fmt"
	"gochat/api/api/friend"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	globalUtils "gochat/utils"

	"github.com/gin-gonic/gin"
)

// GetFriendDetail 获取好友详情
func (h *FriendHandler) GetFriendDetail(ctx *gin.Context) {
	// 从URL参数获取friend_id
	friendIDStr := ctx.Param("friend_id")
	if friendIDStr == "" {
		h.response.JsonError(ctx, nil, "好友ID不能为空")
		return
	}

	// 创建请求对象
	req := &friend.GetFriendDetailRequest{}
	// 这里需要手动解析friend_id，因为它是路径参数
	var friendID uint32
	_, err := fmt.Sscanf(friendIDStr, "%d", &friendID)
	if err != nil {
		h.response.JsonError(ctx, err, "无效的好友ID")
		return
	}
	req.FriendId = friendID

	resp, code, err := h.getFriendDetailLogic(ctx, req)
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

func (h *FriendHandler) getFriendDetailLogic(ctx *gin.Context, req *friend.GetFriendDetailRequest) (resp *friend.GetFriendDetailResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	friendID := uint(req.GetFriendId())

	// 获取好友关系
	friendRelation, err := h.dao.GetFriendDetail(userID, friendID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	if friendRelation == nil {
		return &friend.GetFriendDetailResponse{
			Success: false,
			Message: "好友关系不存在",
		}, 0, nil
	}

	// 获取好友的用户信息（通过UserDao）
	userDao := dao.NewUserDao(globalUtils.DB)
	userInfo, err := userDao.GetUserByID(friendID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	if userInfo == nil {
		return &friend.GetFriendDetailResponse{
			Success: false,
			Message: "用户不存在",
		}, 0, nil
	}

	// 检查在线状态
	wsHub := globalUtils.GetWebSocketHub()
	isOnline := false
	if wsHub != nil {
		isOnline = wsHub.IsUserOnline(friendID)
	}

	// 获取API端口配置
	apiPort := 8080 // 默认端口
	if port := globalUtils.GetApiPort(); port > 0 {
		apiPort = port
	}

	// 构建响应
	friendInfo := &friend.FriendInfo{
		Id:        uint32(friendID),
		Name:      userInfo.Name,
		Phone:     userInfo.Phone,
		Email:     userInfo.Email,
		Avatar:    globalUtils.GetAvatarFullURL(userInfo.Avatar, fmt.Sprintf("http://127.0.0.1:%d", apiPort)),
		IsOnline:  isOnline,
		Remark:    friendRelation.Remark,
		IsBlocked: friendRelation.IsBlocked,
	}

	resp = &friend.GetFriendDetailResponse{
		Friend:  friendInfo,
		Success: true,
		Message: "获取成功",
	}

	return resp, 0, nil
}

