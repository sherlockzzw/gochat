package group

import (
	"time"

	"gochat/api/api/group"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
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

	// 创建群组
	now := time.Now().Unix()
	newGroup := &models.Group{
		Name:        req.GetName(),
		Avatar:      req.GetAvatar(),
		OwnerID:     userID,
		MemberCount: 1, // 创建者自己
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if newGroup.Avatar == "" {
		newGroup.Avatar = avatarUtils.GetDefaultGroupAvatar()
	}

	err = h.dao.CreateGroup(newGroup)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	err = h.dao.AddMember(newGroup.ID, userID, models.GroupRoleOwner, "")
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	if len(req.GetMemberIds()) > 0 {
		memberIDs := make([]int64, 0, len(req.GetMemberIds()))
		for _, id := range req.GetMemberIds() {
			memberIDs = append(memberIDs, int64(id))
		}
		err = h.dao.AddMembers(newGroup.ID, memberIDs, userID)
		if err != nil {

		}
	}

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
		MemberCount: int32(groupInfo.MemberCount),
		CreatedAt:   timestamppb.New(time.Unix(groupInfo.CreatedAt, 0)),
		UpdatedAt:   timestamppb.New(time.Unix(groupInfo.UpdatedAt, 0)),
	}

	return &group.CreateGroupResponse{
		Group:   groupResp,
		Success: true,
	}, 0, nil
}
