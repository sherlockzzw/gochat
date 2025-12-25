package user

import (
	"fmt"
	"time"

	"gochat/api/api/user"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	globalUtils "gochat/utils"

	"github.com/gin-gonic/gin"
)

// UpdateUserProfile 更新用户资料
func (h *UserHandler) UpdateUserProfile(ctx *gin.Context) {
	req, err := analysis.BindParameter[user.UpdateUserProfileRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.updateUserProfileLogic(ctx, &req)
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

func (h *UserHandler) updateUserProfileLogic(ctx *gin.Context, req *user.UpdateUserProfileRequest) (resp *user.UpdateUserProfileResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 获取用户信息
	userInfo, err := h.dao.GetUserByID(userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	if userInfo == nil {
		return nil, code_msg.UserNotExists, nil
	}

	// 更新用户信息（只更新提供的字段）
	updates := make(map[string]interface{})
	
	if req.GetName() != "" {
		updates["name"] = req.GetName()
	}
	if req.GetPhone() != "" {
		updates["phone"] = req.GetPhone()
	}
	if req.GetEmail() != "" {
		updates["email"] = req.GetEmail()
	}
	if req.GetAvatar() != "" {
		updates["avatar"] = req.GetAvatar()
	}
	if req.GetSignature() != "" {
		updates["signature"] = req.GetSignature()
	}

	// 执行更新
	if len(updates) > 0 {
		err = h.dao.UpdateUser(userID, updates)
		if err != nil {
			return nil, code_msg.ServerError, err
		}
	}

	// 重新获取更新后的用户信息
	updatedUser, err := h.dao.GetUserByID(userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 获取API端口配置
	apiPort := 8080
	if port := globalUtils.GetApiPort(); port > 0 {
		apiPort = port
	}
	baseURL := fmt.Sprintf("http://127.0.0.1:%d", apiPort)

	// 构建响应
	userInfoProto := &user.UserInfo{
		Id:        int64(updatedUser.ID),
		Name:      updatedUser.Name,
		Phone:     updatedUser.Phone,
		Email:     updatedUser.Email,
		Avatar:    globalUtils.GetAvatarFullURL(updatedUser.Avatar, baseURL),
		Signature: updatedUser.Signature,
		CreatedAt: time.Unix(updatedUser.CreatedAt, 0).Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Unix(updatedUser.UpdatedAt, 0).Format("2006-01-02 15:04:05"),
	}

	resp = &user.UpdateUserProfileResponse{
		Success: true,
		User:    userInfoProto,
	}

	return resp, 0, nil
}

