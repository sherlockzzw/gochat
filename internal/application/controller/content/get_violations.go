package content

import (
	"gochat/api/admin/content"
	"gochat/internal/pkg/analysis"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetViolations 获取违规内容
func (c *ContentController) GetViolations(ctx *gin.Context) {
	req, err := analysis.BindQuery[content.GetViolationsRequest](ctx, c.response)
	if err != nil {
		return
	}

	page := int(req.GetPage())
	if page <= 0 {
		page = 1
	}
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 20
	}

	violations, total, err := c.violationDao.GetViolations(page, pageSize, req.GetType(), req.GetStatus())
	if err != nil {
		c.response.JsonError(ctx, err, err.Error())
		return
	}

	var violationInfos []*content.ViolationInfo
	for _, v := range violations {
		user, _ := c.userDao.GetUserByID(v.UserID)
		userName := ""
		if user != nil {
			userName = user.Name
		}

		violationInfos = append(violationInfos, &content.ViolationInfo{
			Id:        v.ID,
			Type:      v.Type,
			UserId:    v.UserID,
			UserName:  userName,
			Content:   v.Content,
			Reason:    v.Reason,
			Status:    v.Status,
			CreatedAt: timestamppb.New(time.Unix(v.CreatedAt, 0)),
		})
	}

	resp := &content.GetViolationsResponse{
		Violations: violationInfos,
		Total:      int32(total),
	}

	c.response.JsonSuccess(ctx, resp)
}

