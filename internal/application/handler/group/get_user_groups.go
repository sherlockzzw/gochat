package group

import (
	"time"

	"gochat/api/api/group"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetUserGroups 获取用户群组列表
func (h *GroupHandler) GetUserGroups(ctx *gin.Context) {
	resp, code, err := h.getUserGroupsLogic(ctx)
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

func (h *GroupHandler) getUserGroupsLogic(ctx *gin.Context) (resp *group.GetUserGroupsResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 获取用户加入的群组列表
	groups, err := h.dao.GetUserGroups(userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 构造响应
	groupList := make([]*group.GroupInfo, 0, len(groups))
	for _, g := range groups {
		groupInfo := &group.GroupInfo{
			Id:          uint32(g.ID),
			Name:        g.Name,
			Avatar:      g.Avatar,
			OwnerId:     uint32(g.OwnerID),
			Notice:      g.Notice,
			MemberCount:  int32(g.MemberCount),
			CreatedAt:    timestamppb.New(time.Unix(g.CreatedAt, 0)),
			UpdatedAt:    timestamppb.New(time.Unix(g.UpdatedAt, 0)),
		}
		groupList = append(groupList, groupInfo)
	}

	return &group.GetUserGroupsResponse{
		Groups:     groupList,
		TotalCount: int32(len(groupList)),
	}, 0, nil
}

