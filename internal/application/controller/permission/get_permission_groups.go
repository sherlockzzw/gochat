package permission

import (
	"encoding/json"
	"gochat/api/admin/permission"
	"gochat/internal/pkg/analysis"

	"github.com/gin-gonic/gin"
)

// GetPermissionGroups 获取权限组列表
func (c *PermissionController) GetPermissionGroups(ctx *gin.Context) {
	req, err := analysis.BindQuery[permission.GetPermissionGroupsRequest](ctx, c.response)
	if err != nil {
		return
	}

	page := int(req.GetPage())
	if page <= 0 {
		page = 1
	}
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 10
	}

	groups, total, err := c.permissionDao.GetPermissionGroups(page, pageSize)
	if err != nil {
		c.response.JsonError(ctx, err, err.Error())
		return
	}

	var groupInfos []*permission.PermissionGroupInfo
	for _, g := range groups {
		var permissions []string
		if g.Permissions != "" {
			json.Unmarshal([]byte(g.Permissions), &permissions)
		}
		
		userCount, _ := c.permissionDao.GetGroupUserCount(g.ID)
		
		groupInfos = append(groupInfos, &permission.PermissionGroupInfo{
			Id:          g.ID,
			Name:        g.Name,
			Description: g.Description,
			Permissions: permissions,
			IsDefault:   g.IsDefault,
			UserCount:   int32(userCount),
			CreatedAt:   g.CreatedAt,
			UpdatedAt:   g.UpdatedAt,
		})
	}

	resp := &permission.GetPermissionGroupsResponse{
		Groups: groupInfos,
		Total:  int32(total),
	}

	c.response.JsonSuccess(ctx, resp)
}

