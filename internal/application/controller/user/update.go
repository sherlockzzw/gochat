package user

import (
	admin_user "gochat/api/admin/user"
	"gochat/internal/pkg/analysis"

	"github.com/gin-gonic/gin"
)

// UpdateUser 更新用户
func (h *ControllerUser) UpdateUser(ctx *gin.Context) {
	req, err := analysis.BindParameter[admin_user.UpdateUserRequest](ctx, h.response)
	if err != nil {
		return
	}

	// TODO: 实现更新用户逻辑
	// 1. 获取用户
	// 2. 更新字段
	// 3. 保存
	_ = req // 暂时未使用

	resp := &admin_user.UpdateUserResponse{
		Code:    0,
		Message: "更新成功",
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

