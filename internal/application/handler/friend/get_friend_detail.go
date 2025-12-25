package friend

import (
	"fmt"
	"time"

	"gochat/api/api/friend"
	"gochat/internal/application/handler/common"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	globalUtils "gochat/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetFriendDetail 获取好友详情
func (h *FriendHandler) GetFriendDetail(ctx *gin.Context) {
	req, err := analysis.BindQuery[friend.GetFriendDetailRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.getFriendDetailLogic(ctx, &req)
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

func (h *FriendHandler) getFriendDetailLogic(ctx *gin.Context, req *friend.GetFriendDetailRequest) (resp *friend.GetFriendDetailResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	friendID := int64(req.GetFriendId())

	// 获取好友关系
	friendRelation, err := h.dao.GetFriendDetail(userID, friendID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	if friendRelation == nil {
		return nil, code_msg.FriendRelationNotExists, nil
	}

	// 获取好友的用户信息（通过UserDao）
	userDao := dao.NewUserDao(globalUtils.DB)
	userInfo, err := userDao.GetUserByID(friendID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	if userInfo == nil {
		return nil, code_msg.UserNotExists, nil
	}

	// 检查在线状态
	isOnline := common.IsUserOnline(friendID)

	// 获取API端口配置
	apiPort := 8080 // 默认端口
	if port := globalUtils.GetApiPort(); port > 0 {
		apiPort = port
	}

	// 构建响应
	friendInfo := &friend.FriendInfo{
		Id:        uint32(friendID),
		Name:      userInfo.Name,
		Phone:     userInfo.Phone,
		Email:     userInfo.Email,
		Avatar:    globalUtils.GetAvatarFullURL(userInfo.Avatar, fmt.Sprintf("http://127.0.0.1:%d", apiPort)),
		IsOnline:  isOnline,
		Remark:    friendRelation.Remark,
		IsBlocked: friendRelation.IsBlocked,
		CreatedAt: timestamppb.New(time.Unix(friendRelation.CreatedAt, 0)),
	}

	resp = &friend.GetFriendDetailResponse{
		Friend:  friendInfo,
		Success: true,
	}

	return resp, 0, nil
}
