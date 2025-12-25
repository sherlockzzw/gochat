package group

import (
	"gochat/api/api/group"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// LeaveGroup 退出群组
func (h *GroupHandler) LeaveGroup(ctx *gin.Context) {
	req, err := analysis.BindParameter[group.LeaveGroupRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.leaveGroupLogic(ctx, &req)
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

func (h *GroupHandler) leaveGroupLogic(ctx *gin.Context, req *group.LeaveGroupRequest) (resp *group.LeaveGroupResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	groupID := int64(req.GetGroupId())

	// 验证群组是否存在
	groupInfo, err := h.dao.GetGroupByID(groupID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if groupInfo == nil {
		return nil, code_msg.GroupNotFound, nil
	}

	// 验证用户是否在群组中
	isMember, err := h.dao.IsMemberInGroup(groupID, userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if !isMember {
		return nil, code_msg.NotGroupMember, nil
	}

	// 获取用户角色
	role, err := h.dao.GetMemberRole(groupID, userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 群主不能直接退出，需要先转让群组或解散群组
	if role == models.GroupRoleOwner {
		return nil, code_msg.BadRequest, nil
	}

	// 移除成员（退出群组）
	err = h.dao.RemoveMember(groupID, userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return &group.LeaveGroupResponse{
		Success: true,
	}, 0, nil
}

