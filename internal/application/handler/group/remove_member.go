package group

import (
	"gochat/api/api/group"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	"gochat/internal/infrastructure/models"

	"github.com/gin-gonic/gin"
)

// RemoveGroupMember 移除群成员（踢人）
func (h *GroupHandler) RemoveGroupMember(ctx *gin.Context) {
	req, err := analysis.BindParameter[group.RemoveGroupMemberRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.removeGroupMemberLogic(ctx, &req)
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

func (h *GroupHandler) removeGroupMemberLogic(ctx *gin.Context, req *group.RemoveGroupMemberRequest) (resp *group.RemoveGroupMemberResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	groupID := int64(req.GetGroupId())
	targetUserID := int64(req.GetUserId())

	// 验证群组是否存在
	groupInfo, err := h.dao.GetGroupByID(groupID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if groupInfo == nil {
		return nil, code_msg.GroupNotFound, nil
	}

	// 检查当前用户是否在群组中
	isMember, err := h.dao.IsMemberInGroup(groupID, userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if !isMember {
		return nil, code_msg.NotGroupMember, nil
	}

	// 获取当前用户的角色
	currentUserRole, err := h.dao.GetMemberRole(groupID, userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 获取目标用户的角色
	targetUserRole, err := h.dao.GetMemberRole(groupID, targetUserID)
	if err != nil {
		return nil, code_msg.NotGroupMember, nil
	}

	// 权限检查：只有群主和管理员可以踢人，且不能踢群主
	if currentUserRole != models.GroupRoleOwner && currentUserRole != models.GroupRoleAdmin {
		return nil, code_msg.NoPermission, nil
	}

	if targetUserRole == models.GroupRoleOwner {
		return nil, code_msg.CannotRemoveOwner, nil
	}

	// 管理员不能踢其他管理员（只有群主可以）
	if currentUserRole == models.GroupRoleAdmin && targetUserRole == models.GroupRoleAdmin {
		return nil, code_msg.CannotRemoveAdmin, nil
	}

	// 移除成员
	err = h.dao.RemoveMember(groupID, targetUserID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return &group.RemoveGroupMemberResponse{
		Success: true,
	}, 0, nil
}

