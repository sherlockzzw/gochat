package service

import (
	"github.com/gin-gonic/gin"
	"gochat/models"
	"net/http"
)

func UserList(c *gin.Context) {
	data := models.GetUserList()

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}
