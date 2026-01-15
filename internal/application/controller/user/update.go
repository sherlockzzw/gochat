package user

import (
	admin_user "gochat/api/admin/user"
	"gochat/internal/pkg/analysis"
	"time"

	"github.com/gin-gonic/gin"
)

// UpdateUser 更新用户
func (h *ControllerUser) UpdateUser(ctx *gin.Context) {
	req, err := analysis.BindParameter[admin_user.UpdateUserRequest](ctx, h.response)
	if err != nil {
		return
	}

	if req.GetId() <= 0 {
		h.response.JsonError(ctx, gin.Error{Type: gin.ErrorTypePublic, Meta: "invalid id"}, "invalid id")
		return
	}

	updates := map[string]interface{}{}
	if req.GetName() != "" {
		updates["name"] = req.GetName()
		updates["login_account"] = req.GetName()
	}
	if req.GetPhone() != "" {
		updates["phone"] = req.GetPhone()
	}
	if req.GetEmail() != "" {
		updates["email"] = req.GetEmail()
	}
	if req.GetStatus() == 0 || req.GetStatus() == 1 {
		updates["status"] = req.GetStatus()
	}

	if len(updates) == 0 {
		resp := &admin_user.UpdateUserResponse{Code: 0, Message: "更新成功"}
		h.response.JsonSuccess(ctx, resp)
		return
	}
	updates["updated_at"] = time.Now().Unix()

	user, err := h.dao.GetUserByID(req.GetId())
	if err != nil {
		h.response.JsonError(ctx, err, "获取用户失败")
		return
	}
	if user == nil {
		h.response.JsonError(ctx, gin.Error{Type: gin.ErrorTypePublic, Meta: "用户不存在"}, "用户不存在")
		return
	}

	if err := h.dao.UpdateUser(req.GetId(), updates); err != nil {
		h.response.JsonError(ctx, err, "更新失败")
		return
	}

	resp := &admin_user.UpdateUserResponse{Code: 0, Message: "更新成功"}
	if req.GetStatus() == 0 {
		//h.logger.Info("admin update user", nil)
	}

	h.response.JsonSuccess(ctx, resp)
}

// DisableUser 禁用用户
func (h *ControllerUser) DisableUser(ctx *gin.Context) {
	req, err := analysis.BindParameter[admin_user.DisableUserRequest](ctx, h.response)
	if err != nil {
		return
	}

	// TODO: 实现禁用用户逻辑
	// 1. 更新用户状态
	// 2. 记录日志
	_ = req // 暂时未使用

	resp := &admin_user.DisableUserResponse{
		Code:    0,
		Message: "操作成功",
	}

	h.response.JsonSuccess(ctx, resp)
}
