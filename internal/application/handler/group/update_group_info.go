package group

import (
	"gochat/api/api/group"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UpdateGroupInfo 更新群组信息（仅群主可操作）
func (h *GroupHandler) UpdateGroupInfo(ctx *gin.Context) {
	req, err := analysis.BindParameter[group.UpdateGroupInfoRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.updateGroupInfoLogic(ctx, &req)
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

func (h *GroupHandler) updateGroupInfoLogic(ctx *gin.Context, req *group.UpdateGroupInfoRequest) (resp *group.UpdateGroupInfoResponse, errCode code_msg.BusinessCode, err error) {
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

	// 检查是否为群主（只有群主可以配置群组信息）
	isOwner, code, err := h.checkIsOwner(groupID, userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if !isOwner {
		return nil, code, nil
	}

	// 更新群组信息
	now := time.Now().Unix()
	updates := make(map[string]interface{})
	updates["updated_at"] = now

	if req.GetName() != "" {
		updates["name"] = req.GetName()
	}
	if req.GetAvatar() != "" {
		updates["avatar"] = req.GetAvatar()
	}
	if req.GetNotice() != "" {
		updates["notice"] = req.GetNotice()
	}

	// 执行更新
	updateGroup := &models.Group{
		ID:        groupID,
		Name:      groupInfo.Name,
		Avatar:    groupInfo.Avatar,
		Notice:    groupInfo.Notice,
		UpdatedAt: now,
	}

	if req.GetName() != "" {
		updateGroup.Name = req.GetName()
	}
	if req.GetAvatar() != "" {
		updateGroup.Avatar = req.GetAvatar()
	}
	if req.GetNotice() != "" {
		updateGroup.Notice = req.GetNotice()
	}

	err = h.dao.UpdateGroup(updateGroup)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 重新获取更新后的群组信息
	updatedGroup, err := h.dao.GetGroupByID(groupID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 构造响应
	groupResp := &group.GroupInfo{
		Id:          uint32(updatedGroup.ID),
		Name:        updatedGroup.Name,
		Avatar:      updatedGroup.Avatar,
		OwnerId:     uint32(updatedGroup.OwnerID),
		Notice:      updatedGroup.Notice,
		MemberCount: int32(updatedGroup.MemberCount),
		CreatedAt:   timestamppb.New(time.Unix(updatedGroup.CreatedAt, 0)),
		UpdatedAt:   timestamppb.New(time.Unix(updatedGroup.UpdatedAt, 0)),
	}

	return &group.UpdateGroupInfoResponse{
		Group:   groupResp,
		Success: true,
	}, 0, nil
}

