package content

import (
	"gochat/api/admin/content"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"time"

	"github.com/gin-gonic/gin"
)

// BlockAbnormalOperation 拦截异常操作
func (c *ContentController) BlockAbnormalOperation(ctx *gin.Context) {
	req, err := analysis.BindParameter[content.BlockAbnormalOperationRequest](ctx, c.response)
	if err != nil {
		return
	}

	// TODO: 实现拦截逻辑
	// 1. 记录拦截操作
	// 2. 可能需要禁用用户或限制操作

	// 记录日志
	c.adminLogDao.CreateLog(&models.AdminLog{
		AdminID:     0, // TODO
		ActionType:  models.ActionTypeContentBlock,
		Description: "拦截异常操作: " + req.GetReason(),
		TargetType:  "content",
		TargetID:    req.GetId(),
		CreatedAt:   time.Now().Unix(),
	})

	c.response.JsonSuccess(ctx, &content.BlockAbnormalOperationResponse{
		Code:    0,
		Message: "拦截成功",
	})
}

