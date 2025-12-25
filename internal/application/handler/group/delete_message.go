package group

import (
	"gochat/api/api/group"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	globalUtils "gochat/utils"

	"github.com/gin-gonic/gin"
)

// DeleteGroupMessage 管理员删除群组消息
func (h *GroupHandler) DeleteGroupMessage(ctx *gin.Context) {
	req, err := analysis.BindParameter[group.DeleteGroupMessageRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.deleteGroupMessageLogic(ctx, &req)
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

func (h *GroupHandler) deleteGroupMessageLogic(ctx *gin.Context, req *group.DeleteGroupMessageRequest) (resp *group.DeleteGroupMessageResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	groupID := int64(req.GetGroupId())
	messageID := req.GetMessageId()

	// 验证群组是否存在
	groupInfo, err := h.dao.GetGroupByID(groupID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if groupInfo == nil {
		return nil, code_msg.GroupNotFound, nil
	}

	// 检查权限：群主或管理员可以删除消息
	isAdminOrOwner, code, err := h.checkIsAdminOrOwner(groupID, userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if !isAdminOrOwner {
		return nil, code, nil
	}

		// 验证消息是否存在且属于该群组
		// 注意：这里需要创建ChatDao，因为GroupHandler中没有ChatDao
		chatDao := dao.NewChatDao(globalUtils.DB, globalUtils.MongoDB)
		messages, err := chatDao.GetMessagesByIDs([]string{messageID})
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if len(messages) == 0 {
		return nil, code_msg.BadRequest, nil
	}

	msg := messages[0]
	if msg.GroupID != groupID {
		return nil, code_msg.BadRequest, nil
	}

	// 删除消息（管理员删除是双向删除，所有成员都看不到）
	err = chatDao.DeleteMessage(messageID, userID, true)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return &group.DeleteGroupMessageResponse{
		Success: true,
	}, 0, nil
}

