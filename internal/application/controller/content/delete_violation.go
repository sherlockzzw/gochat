package content

import (
	"gochat/api/admin/content"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"time"

	"github.com/gin-gonic/gin"
)

// DeleteViolation 删除违规内容
func (c *ContentController) DeleteViolation(ctx *gin.Context) {
	req, err := analysis.BindParameter[content.DeleteViolationRequest](ctx, c.response)
	if err != nil {
		return
	}

	if err := c.violationDao.DeleteViolation(req.GetId()); err != nil {
		c.response.JsonError(ctx, err, err.Error())
		return
	}

	// 记录日志
	c.adminLogDao.CreateLog(&models.AdminLog{
		AdminID:     0, // TODO
		ActionType:  models.ActionTypeContentDelete,
		Description: "删除违规内容",
		TargetType:  "content",
		TargetID:    req.GetId(),
		CreatedAt:   time.Now().Unix(),
	})

	c.response.JsonSuccess(ctx, &content.DeleteViolationResponse{
		Code:    0,
		Message: "删除成功",
	})
}

