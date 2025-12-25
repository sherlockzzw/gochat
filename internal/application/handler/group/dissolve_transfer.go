package group

import (
	"gochat/api/api/group"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// DissolveGroup 解散群组
func (h *GroupHandler) DissolveGroup(ctx *gin.Context) {
	req, err := analysis.BindParameter[group.DissolveGroupRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.dissolveGroupLogic(ctx, &req)
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

func (h *GroupHandler) dissolveGroupLogic(ctx *gin.Context, req *group.DissolveGroupRequest) (resp *group.DissolveGroupResponse, errCode code_msg.BusinessCode, err error) {
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

	// 检查是否为群主
	if groupInfo.OwnerID != userID {
		return nil, code_msg.NoPermission, nil
	}

	// 解散群组
	err = h.dao.DissolveGroup(groupID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return &group.DissolveGroupResponse{
		Success: true,
	}, 0, nil
}

// TransferGroup 转让群组
func (h *GroupHandler) TransferGroup(ctx *gin.Context) {
	req, err := analysis.BindParameter[group.TransferGroupRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.transferGroupLogic(ctx, &req)
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

func (h *GroupHandler) transferGroupLogic(ctx *gin.Context, req *group.TransferGroupRequest) (resp *group.TransferGroupResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	groupID := int64(req.GetGroupId())
	newOwnerID := int64(req.GetNewOwnerId())

	// 验证群组是否存在
	groupInfo, err := h.dao.GetGroupByID(groupID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if groupInfo == nil {
		return nil, code_msg.GroupNotFound, nil
	}

	// 检查是否为群主
	if groupInfo.OwnerID != userID {
		return nil, code_msg.NoPermission, nil
	}

	// 不能转让给自己
	if newOwnerID == userID {
		return nil, code_msg.BadRequest, nil
	}

	// 验证新群主是否在群组中
	isMember, err := h.dao.IsMemberInGroup(groupID, newOwnerID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if !isMember {
		return nil, code_msg.UserNotInGroup, nil
	}

	// 转让群组
	err = h.dao.TransferGroup(groupID, userID, newOwnerID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return &group.TransferGroupResponse{
		Success: true,
	}, 0, nil
}

