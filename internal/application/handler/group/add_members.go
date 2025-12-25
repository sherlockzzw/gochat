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

	groupID := int64(req.GetGroupId())
	userIDs := req.GetUserIds()

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

	// 检查每个用户是否已在群组中（前端也要验证，但后端也要验证）
	failedUserIDs := make([]uint32, 0)
	validUserIDs := make([]int64, 0)

	for _, uid := range userIDs {
		uidInt64 := int64(uid)
		isInGroup, err := h.dao.IsMemberInGroup(groupID, uidInt64)
		if err != nil {
			// 记录错误但继续处理其他用户
			continue
		}
		if isInGroup {
			failedUserIDs = append(failedUserIDs, uid)
		} else {
			validUserIDs = append(validUserIDs, uidInt64)
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
		Success:       true,
		AddedCount:    addedCount,
		FailedUserIds: failedUserIDs,
	}, 0, nil
}

