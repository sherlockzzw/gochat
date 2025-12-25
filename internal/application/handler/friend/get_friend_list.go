package friend

import (
	"fmt"
	"gochat/api/api/friend"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	globalUtils "gochat/utils"

	"github.com/gin-gonic/gin"
)

// GetFriendList 获取好友列表
func (h *FriendHandler) GetFriendList(ctx *gin.Context) {
	req, err := analysis.BindParameter[friend.GetFriendListRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.getFriendListLogic(ctx, &req)
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

func (h *FriendHandler) getFriendListLogic(ctx *gin.Context, req *friend.GetFriendListRequest) (resp *friend.GetFriendListResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	// 获取在线用户ID列表
	wsHub := globalUtils.GetWebSocketHub()
	var onlineUserIDs []int64
	if wsHub != nil {
		onlineUserIDs = wsHub.GetOnlineUserIDs()
	} else {
		// WebSocket Hub未初始化，返回空列表
		onlineUserIDs = []int64{}
	}

	// 获取好友列表（包含用户信息和在线状态）
	friends, err := h.dao.GetFriendListWithUserInfoAndOnlineStatus(userID, page, pageSize, onlineUserIDs)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 转换为响应格式
	var friendInfos []*friend.FriendInfo
	for _, f := range friends {
		// 获取API端口配置
		apiPort := 8080 // 默认端口
		if port := globalUtils.GetApiPort(); port > 0 {
			apiPort = port
		}

		friendInfo := &friend.FriendInfo{
			Id:        uint32(f.FriendID),
			Name:      f.FriendName,
			Phone:     f.FriendPhone,
			Email:     f.FriendEmail,
			Avatar:    globalUtils.GetAvatarFullURL(f.FriendAvatar, fmt.Sprintf("http://127.0.0.1:%d", apiPort)), // 返回完整头像URL
			IsOnline:  f.IsOnline,                                                                                // 使用WebSocket Hub查询的在线状态
			Remark:    f.Remark,
			IsBlocked: f.IsBlocked,
		}
		friendInfos = append(friendInfos, friendInfo)
	}

	resp = &friend.GetFriendListResponse{
		Friends:    friendInfos,
		TotalCount: int32(len(friendInfos)),
	}

	return resp, 0, nil
}
