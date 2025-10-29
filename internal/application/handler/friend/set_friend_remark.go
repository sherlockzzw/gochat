package friend

import (
	"gochat/api/api/friend"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"

	"github.com/gin-gonic/gin"
)

// SetFriendRemark 设置好友备注
func (h *FriendHandler) SetFriendRemark(ctx *gin.Context) {
	req, err := analysis.BindParameter[friend.SetFriendRemarkRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.setFriendRemarkLogic(ctx, &req)
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

func (h *FriendHandler) setFriendRemarkLogic(ctx *gin.Context, req *friend.SetFriendRemarkRequest) (resp *friend.SetFriendRemarkResponse, errCode code_msg.BusinessCode, err error) {
	// 获取当前用户ID（从token中解析）
	userID := uint(1) // TODO: 从JWT token中获取真实用户ID

	friendID := uint(req.GetFriendId())
	remark := req.GetRemark()

	// 检查是否为好友
	existingFriend, err := h.dao.CheckIsFriend(userID, friendID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if existingFriend == nil {
		return &friend.SetFriendRemarkResponse{
			Success: false,
			Message: "不是好友关系",
		}, 0, nil
	}

	// 设置好友备注
	err = h.dao.SetFriendRemark(userID, friendID, remark)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return &friend.SetFriendRemarkResponse{
		Success: true,
		Message: "设置备注成功",
	}, 0, nil
}
