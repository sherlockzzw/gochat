package group

import (
	"gochat/api/api/group"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// RemoveAdmin 移除管理员
func (h *GroupHandler) RemoveAdmin(ctx *gin.Context) {
	req, err := analysis.BindParameter[group.RemoveAdminRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.removeAdminLogic(ctx, &req)
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

func (h *GroupHandler) removeAdminLogic(ctx *gin.Context, req *group.RemoveAdminRequest) (resp *group.RemoveAdminResponse, errCode code_msg.BusinessCode, err error) {
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

	// 检查是否为群主
	isOwner, code, err := h.checkIsOwner(groupID, userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if !isOwner {
		return nil, code, nil
	}

	// 验证目标用户是否在群组中
	isMember, err := h.dao.IsMemberInGroup(groupID, targetUserID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if !isMember {
		return nil, code_msg.UserNotInGroup, nil
	}

	// 获取目标用户当前角色
	targetRole, err := h.dao.GetMemberRole(groupID, targetUserID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 如果已经是普通成员，直接返回成功
	if targetRole == models.GroupRoleMember {
		return &group.RemoveAdminResponse{
			Success: true,
		}, 0, nil
	}

	// 不能移除群主的管理员身份（群主不是管理员）
	if targetRole == models.GroupRoleOwner {
		return nil, code_msg.BadRequest, nil
	}

	// 更新角色为普通成员
	err = h.dao.UpdateMemberRole(groupID, targetUserID, models.GroupRoleMember)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return &group.RemoveAdminResponse{
		Success: true,
	}, 0, nil
}

