package group

import (
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/code_msg"
)

// checkPermission 检查用户权限
// 返回: (是否有权限, 错误码, 错误)
func (h *GroupHandler) checkPermission(groupID, userID int64, requiredRole string) (bool, code_msg.BusinessCode, error) {
	// 获取用户角色
	role, err := h.dao.GetMemberRole(groupID, userID)
	if err != nil {
		return false, code_msg.ServerError, err
	}

	// 检查权限
	switch requiredRole {
	case "owner":
		// 只有群主
		return role == models.GroupRoleOwner, code_msg.NoPermission, nil
	case "admin":
		// 群主或管理员
		return role == models.GroupRoleOwner || role == models.GroupRoleAdmin, code_msg.NoPermission, nil
	case "member":
		// 任何成员都可以
		return role != "", code_msg.NoPermission, nil
	default:
		return false, code_msg.BadRequest, nil
	}
}

// checkIsOwner 检查是否为群主
func (h *GroupHandler) checkIsOwner(groupID, userID int64) (bool, code_msg.BusinessCode, error) {
	return h.checkPermission(groupID, userID, "owner")
}

// checkIsAdminOrOwner 检查是否为管理员或群主
func (h *GroupHandler) checkIsAdminOrOwner(groupID, userID int64) (bool, code_msg.BusinessCode, error) {
	return h.checkPermission(groupID, userID, "admin")
}

// checkIsMember 检查是否为成员
func (h *GroupHandler) checkIsMember(groupID, userID int64) (bool, code_msg.BusinessCode, error) {
	return h.checkPermission(groupID, userID, "member")
}

// getMemberRole 获取成员角色
func (h *GroupHandler) getMemberRole(groupID, userID int64) (string, code_msg.BusinessCode, error) {
	role, err := h.dao.GetMemberRole(groupID, userID)
	if err != nil {
		return "", code_msg.ServerError, err
	}
	if role == "" {
		return "", code_msg.NotGroupMember, nil
	}
	return role, 0, nil
}

