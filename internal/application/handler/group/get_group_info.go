package group

import (
	"time"

	"gochat/api/api/group"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetGroupInfo 获取群组信息
func (h *GroupHandler) GetGroupInfo(ctx *gin.Context) {
	req, err := analysis.BindQuery[group.GetGroupInfoRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.getGroupInfoLogic(ctx, &req)
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

func (h *GroupHandler) getGroupInfoLogic(ctx *gin.Context, req *group.GetGroupInfoRequest) (resp *group.GetGroupInfoResponse, errCode code_msg.BusinessCode, err error) {
	groupID := int64(req.GetGroupId())

	// 获取群组信息
	groupInfo, err := h.dao.GetGroupByID(groupID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if groupInfo == nil {
		return &group.GetGroupInfoResponse{}, code_msg.NotFound, nil
	}

	// 获取群成员列表
	members, err := h.dao.GetGroupMembers(groupID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 获取在线用户列表（从WebSocket Hub）
	onlineUserIDs := make(map[int64]bool)
	if utils.WSHub != nil {
		onlineIDs := utils.WSHub.GetOnlineUserIDs()
		for _, uid := range onlineIDs {
			onlineUserIDs[uid] = true
		}
	}

	// 构造群组信息
	groupResp := &group.GroupInfo{
		Id:          uint32(groupInfo.ID),
		Name:        groupInfo.Name,
		Avatar:      groupInfo.Avatar,
		OwnerId:     uint32(groupInfo.OwnerID),
		Notice:      groupInfo.Notice,
		MemberCount:  int32(groupInfo.MemberCount),
		CreatedAt:    timestamppb.New(time.Unix(groupInfo.CreatedAt, 0)),
		UpdatedAt:    timestamppb.New(time.Unix(groupInfo.UpdatedAt, 0)),
	}

	// 构造成员列表
	memberList := make([]*group.GroupMemberInfo, 0, len(members))
	for _, m := range members {
		isOnline := onlineUserIDs[m.UserID]
		memberInfo := &group.GroupMemberInfo{
			Id:         uint32(m.ID),
			GroupId:    uint32(m.GroupID),
			UserId:     uint32(m.UserID),
			Role:       m.Role,
			Nickname:   m.Nickname,
			UserName:   m.UserName,
			UserAvatar: m.UserAvatar,
			UserPhone:  m.UserPhone,
			IsOnline:   isOnline,
			JoinedAt:   timestamppb.New(time.Unix(m.JoinedAt, 0)),
		}
		memberList = append(memberList, memberInfo)
	}

	return &group.GetGroupInfoResponse{
		Group:   groupResp,
		Members: memberList,
	}, 0, nil
}

