package group

import (
	"gochat/api/api/group"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// AddGroupMembers 添加群成员
func (h *GroupHandler) AddGroupMembers(ctx *gin.Context) {
	req, err := analysis.BindParameter[group.AddGroupMembersRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.addGroupMembersLogic(ctx, &req)
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

func (h *GroupHandler) addGroupMembersLogic(ctx *gin.Context, req *group.AddGroupMembersRequest) (resp *group.AddGroupMembersResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	groupID := uint(req.GetGroupId())
	userIDs := req.GetUserIds()

	// 验证群组是否存在
	groupInfo, err := h.dao.GetGroupByID(groupID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if groupInfo == nil {
		return &group.AddGroupMembersResponse{
			Success:      false,
			ErrorMessage: "群组不存在",
		}, 0, nil
	}

	// 检查当前用户是否在群组中
	isMember, err := h.dao.IsMemberInGroup(groupID, userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if !isMember {
		return &group.AddGroupMembersResponse{
			Success:      false,
			ErrorMessage: "您不在该群组中",
		}, 0, nil
	}

	// 验证要添加的用户ID列表
	if len(userIDs) == 0 {
		return &group.AddGroupMembersResponse{
			Success:      false,
			ErrorMessage: "请选择要添加的用户",
		}, 0, nil
	}

	// 检查每个用户是否已在群组中（前端也要验证，但后端也要验证）
	failedUserIDs := make([]uint32, 0)
	validUserIDs := make([]uint, 0)

	for _, uid := range userIDs {
		uidUint := uint(uid)
		isInGroup, err := h.dao.IsMemberInGroup(groupID, uidUint)
		if err != nil {
			// 记录错误但继续处理其他用户
			continue
		}
		if isInGroup {
			failedUserIDs = append(failedUserIDs, uid)
		} else {
			validUserIDs = append(validUserIDs, uidUint)
		}
	}

	// 添加有效的用户
	if len(validUserIDs) > 0 {
		err = h.dao.AddMembers(groupID, validUserIDs, userID)
		if err != nil {
			return nil, code_msg.ServerError, err
		}
	}

	addedCount := int32(len(validUserIDs))

	return &group.AddGroupMembersResponse{
		Success:        true,
		AddedCount:     addedCount,
		FailedUserIds:  failedUserIDs,
		ErrorMessage:   "",
	}, 0, nil
}

