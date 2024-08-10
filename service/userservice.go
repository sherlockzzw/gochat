package service

import (
	"github.com/gin-gonic/gin"
	"gochat/models"
	"net/http"
)

// UserList
// @Tags 测试用户列表
// @Success 200 {string} json{"code","data}
// @Router /user/list [get]
func UserList(c *gin.Context) {
	data := models.GetUserList()

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}
