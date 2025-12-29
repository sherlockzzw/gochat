package group

import (
	"gochat/api/api/group"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// MuteMember 禁言成员
func (h *GroupHandler) MuteMember(ctx *gin.Context) {
	req, err := analysis.BindParameter[group.MuteMemberRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.muteMemberLogic(ctx, &req)
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

func (h *GroupHandler) muteMemberLogic(ctx *gin.Context, req *group.MuteMemberRequest) (resp *group.MuteMemberResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	groupID := int64(req.GetGroupId())
	targetUserID := int64(req.GetUserId())
	muteDuration := req.GetMuteDuration()
	reason := req.GetReason()

	// 验证群组是否存在
	groupInfo, err := h.dao.GetGroupByID(groupID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if groupInfo == nil {
		return nil, code_msg.GroupNotFound, nil
	}

	// 检查权限：群主或管理员可以禁言
	isAdminOrOwner, code, err := h.checkIsAdminOrOwner(groupID, userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if !isAdminOrOwner {
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

	// 获取目标用户角色
	targetRole, err := h.dao.GetMemberRole(groupID, targetUserID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 不能禁言群主
	if targetRole == models.GroupRoleOwner {
		return nil, code_msg.BadRequest, nil
	}

	// 管理员不能禁言其他管理员（只有群主可以）
	currentRole, err := h.dao.GetMemberRole(groupID, userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if currentRole == models.GroupRoleAdmin && targetRole == models.GroupRoleAdmin {
		return nil, code_msg.NoPermission, nil
	}

	// 计算禁言到期时间
	now := time.Now().Unix()
	var mutedUntil int64
	if muteDuration > 0 {
		mutedUntil = now + int64(muteDuration)
	} else {
		mutedUntil = 0 // 0表示永久禁言
	}

	// 检查是否已有禁言记录
	existingMute, err := h.muteDao.GetMute(groupID, targetUserID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	if existingMute != nil {
		// 更新禁言记录
		err = h.muteDao.UpdateMute(groupID, targetUserID, mutedUntil, reason)
	} else {
		// 创建新禁言记录
		mute := &models.GroupMute{
			GroupID:   groupID,
			UserID:    targetUserID,
			MutedBy:   userID,
			MutedUntil: mutedUntil,
			Reason:    reason,
		}
		err = h.muteDao.CreateMute(mute)
	}

	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return &group.MuteMemberResponse{
		Success: true,
	}, 0, nil
}

// UnmuteMember 解除禁言
func (h *GroupHandler) UnmuteMember(ctx *gin.Context) {
	req, err := analysis.BindParameter[group.UnmuteMemberRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.unmuteMemberLogic(ctx, &req)
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

func (h *GroupHandler) unmuteMemberLogic(ctx *gin.Context, req *group.UnmuteMemberRequest) (resp *group.UnmuteMemberResponse, errCode code_msg.BusinessCode, err error) {
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

	// 检查权限：群主或管理员可以解除禁言
	isAdminOrOwner, code, err := h.checkIsAdminOrOwner(groupID, userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if !isAdminOrOwner {
		return nil, code, nil
	}

	// 解除禁言
	err = h.muteDao.RemoveMute(groupID, targetUserID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return &group.UnmuteMemberResponse{
		Success: true,
	}, 0, nil
}


