package group

import (
	"gochat/api/api/group"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"

	"github.com/gin-gonic/gin"
)

// CheckMemberInGroup 检查用户是否在群组中
func (h *GroupHandler) CheckMemberInGroup(ctx *gin.Context) {
	req, err := analysis.BindQuery[group.CheckMemberInGroupRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.checkMemberInGroupLogic(ctx, &req)
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

func (h *GroupHandler) checkMemberInGroupLogic(ctx *gin.Context, req *group.CheckMemberInGroupRequest) (resp *group.CheckMemberInGroupResponse, errCode code_msg.BusinessCode, err error) {
	groupID := uint(req.GetGroupId())
	userID := uint(req.GetUserId())

	// 检查用户是否在群组中
	isMember, err := h.dao.IsMemberInGroup(groupID, userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	if !isMember {
		return &group.CheckMemberInGroupResponse{
			IsMember: false,
			Role:     "",
		}, 0, nil
	}

	// 获取用户角色
	role, err := h.dao.GetMemberRole(groupID, userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return &group.CheckMemberInGroupResponse{
		IsMember: true,
		Role:     role,
	}, 0, nil
}

