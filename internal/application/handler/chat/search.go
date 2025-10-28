package chat

import (
	"gochat/api/api/chat"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"

	"github.com/gin-gonic/gin"
)

// SearchUser 搜索用户
func (h *ChatHandler) SearchUser(ctx *gin.Context) {
	req, err := analysis.BindQuery[chat.SearchUserRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.searchUserLogic(ctx, req)
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

func (h *ChatHandler) searchUserLogic(ctx *gin.Context, req chat.SearchUserRequest) (resp *chat.SearchUserResponse, errCode code_msg.BusinessCode, err error) {
	// 搜索用户
	users, err := h.dao.SearchUsers(req.GetPhone(), req.GetName(), 20)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 转换为响应格式
	var userInfos []*chat.UserInfo
	for _, user := range users {
		userInfo := &chat.UserInfo{
			Id:    uint32(user.ID),
			Name:  user.Name,
			Phone: user.Phone,
			Email: user.Email,
			// TODO: 添加头像和在线状态
		}
		userInfos = append(userInfos, userInfo)
	}

	resp = &chat.SearchUserResponse{
		Users:      userInfos,
		TotalCount: int32(len(userInfos)),
	}

	return resp, 0, nil
}
