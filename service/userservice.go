package service

import (
	"github.com/gin-gonic/gin"
	"gochat/models"
	"gochat/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// UserList
// @Tags 测试用户列表
// @Success 200 {string} json{"code","data"}
// @Router /user/list [get]
func UserList(c *gin.Context) {
	data := models.GetUserList()

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// CreateUser
// @Tags 创建用户
// @Accept json
// @Produce json
// @Param user body models.UserBasic true "用户信息"
// @Success 200 {string} json{"code","data","message"}
// @Router /user/add [post]
func CreateUser(c *gin.Context) {
	var user = models.UserBasic{}

	utils.DB.Where("name = ?", user.Name).First(&user)
	if user.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "用户已经存在",
		})
		return
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	user.Password = string(hashedPassword)

	if err := models.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    "",
	})
}
