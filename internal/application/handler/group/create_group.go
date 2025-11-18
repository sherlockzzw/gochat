package group

import (
	"gochat/api/api/group"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	"gochat/models"
	avatarUtils "gochat/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CreateGroup 创建群组
func (h *GroupHandler) CreateGroup(ctx *gin.Context) {
	req, err := analysis.BindParameter[group.CreateGroupRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.createGroupLogic(ctx, &req)
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

func (h *GroupHandler) createGroupLogic(ctx *gin.Context, req *group.CreateGroupRequest) (resp *group.CreateGroupResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 验证群名称
	if req.GetName() == "" {
		return &group.CreateGroupResponse{
			Success:      false,
			ErrorMessage: "群名称不能为空",
		}, 0, nil
	}

	// 创建群组
	newGroup := &models.Group{
		Name:        req.GetName(),
		Avatar:      req.GetAvatar(),
		OwnerID:     userID,
		MemberCount:  1, // 创建者自己
	}

	if newGroup.Avatar == "" {
		// 如果没有提供头像，使用默认头像
		newGroup.Avatar = avatarUtils.GetDefaultGroupAvatar()
	}

	err = h.dao.CreateGroup(newGroup)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 添加创建者为群主
	err = h.dao.AddMember(newGroup.ID, userID, models.GroupRoleOwner, "")
	if err != nil {
		// 如果添加成员失败，删除群组
		// 这里简化处理，实际应该使用事务
		return nil, code_msg.ServerError, err
	}

	// 如果有初始成员，添加他们
	if len(req.GetMemberIds()) > 0 {
		// 将[]uint32转换为[]uint
		memberIDs := make([]uint, 0, len(req.GetMemberIds()))
		for _, id := range req.GetMemberIds() {
			memberIDs = append(memberIDs, uint(id))
		}
		err = h.dao.AddMembers(newGroup.ID, memberIDs, userID)
		if err != nil {
			// 记录错误但不影响群组创建
			// 实际应该记录日志
		}
	}

	// 重新获取群组信息（包含更新后的成员数量）
	groupInfo, err := h.dao.GetGroupByID(newGroup.ID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 构造响应
	groupResp := &group.GroupInfo{
		Id:          uint32(groupInfo.ID),
		Name:        groupInfo.Name,
		Avatar:      groupInfo.Avatar,
		OwnerId:     uint32(groupInfo.OwnerID),
		Notice:      groupInfo.Notice,
		MemberCount:  int32(groupInfo.MemberCount),
		CreatedAt:    timestamppb.New(groupInfo.CreatedAt),
		UpdatedAt:    timestamppb.New(groupInfo.UpdatedAt),
	}

	return &group.CreateGroupResponse{
		Group:   groupResp,
		Success: true,
	}, 0, nil
}

