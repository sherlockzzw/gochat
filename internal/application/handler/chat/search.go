package chat

import (
	"gochat/api/api/chat"
	"gochat/internal/application/handler/common"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SearchUser 搜索用户
func (h *ChatHandler) SearchUser(ctx *gin.Context) {
	req, err := analysis.BindQuery[chat.SearchUserRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.searchUserLogic(ctx, &req)
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

func (h *ChatHandler) searchUserLogic(ctx *gin.Context, req *chat.SearchUserRequest) (resp *chat.SearchUserResponse, errCode code_msg.BusinessCode, err error) {
	// 获取在线用户ID列表（用于在线状态筛选）
	var onlineUserIDs []int64
	if req.GetOnlineOnly() {
		onlineUserIDs = common.GetOnlineUserIDs()
	}

	// 搜索用户
	users, err := h.dao.SearchUsers(
		req.GetKeyword(),
		req.GetLoginAccount(),
		req.GetOnlineOnly(),
		onlineUserIDs,
		50, // 增加搜索限制到50
	)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 获取在线用户ID集合（用于判断是否在线）
	onlineUserIDMap := make(map[int64]bool)
	if !req.GetOnlineOnly() {
		onlineUserIDs = common.GetOnlineUserIDs()
	}
	for _, id := range onlineUserIDs {
		onlineUserIDMap[id] = true
	}

	// 转换为响应格式
	var userInfos []*chat.UserInfo
	for _, user := range users {
		isOnline := onlineUserIDMap[user.ID]
		
		var lastSeen *timestamppb.Timestamp
		if user.LastActiveTime > 0 {
			lastSeen = timestamppb.New(time.Unix(user.LastActiveTime, 0))
		}

		userInfo := &chat.UserInfo{
			Id:        uint32(user.ID),
			Name:      user.Name,
			Phone:     user.Phone,
			Email:     user.Email,
			Avatar:    user.Avatar,
			IsOnline:  isOnline,
			LastSeen:  lastSeen,
		}
		userInfos = append(userInfos, userInfo)
	}

	resp = &chat.SearchUserResponse{
		Users:      userInfos,
		TotalCount: int32(len(userInfos)),
	}

	return resp, 0, nil
}
