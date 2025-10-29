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

// GetFriendRequests 获取好友申请列表
func (h *FriendHandler) GetFriendRequests(ctx *gin.Context) {
	req, err := analysis.BindParameter[friend.GetFriendRequestsRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.getFriendRequestsLogic(ctx, &req)
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

func (h *FriendHandler) getFriendRequestsLogic(ctx *gin.Context, req *friend.GetFriendRequestsRequest) (resp *friend.GetFriendRequestsResponse, errCode code_msg.BusinessCode, err error) {
	// 获取当前用户ID（从token中解析）
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

	// 获取好友申请列表（包含申请人信息）
	requests, err := h.dao.GetFriendRequestsWithUserInfo(userID, page, pageSize)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 转换为响应格式
	var requestInfos []*friend.FriendRequestInfo
	for _, r := range requests {
		// 获取API端口配置
		apiPort := 8080 // 默认端口
		if port := globalUtils.GetApiPort(); port > 0 {
			apiPort = port
		}

		requestInfo := &friend.FriendRequestInfo{
			Id:             uint32(r.ID),
			FromUserId:     uint32(r.FromUserID),
			FromUserName:   r.FromUserName,
			FromUserAvatar: globalUtils.GetAvatarFullURL(r.FromUserAvatar, fmt.Sprintf("http://127.0.0.1:%d", apiPort)), // 返回完整头像URL
			Message:        r.Message,
			Status:         r.Status,
			CreatedAt:      r.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		requestInfos = append(requestInfos, requestInfo)
	}

	resp = &friend.GetFriendRequestsResponse{
		Requests:   requestInfos,
		TotalCount: int32(len(requestInfos)),
	}

	return resp, 0, nil
}
