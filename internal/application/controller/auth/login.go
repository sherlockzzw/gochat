package auth

import (
	"gochat/api/admin/auth"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Login 管理员登录
func (c *AuthController) Login(ctx *gin.Context) {
	req, err := analysis.BindParameter[auth.LoginRequest](ctx, c.response)
	if err != nil {
		return
	}

	resp, code, err := c.loginLogic(ctx, &req)
	if code != 0 {
		c.response.JsonErrorFixation(ctx, code)
		return
	}
	if err != nil {
		c.response.JsonError(ctx, err, err.Error())
		return
	}

	c.response.JsonSuccess(ctx, resp)
}

func (c *AuthController) loginLogic(ctx *gin.Context, req *auth.LoginRequest) (resp *auth.LoginResponse, errCode code_msg.BusinessCode, err error) {
	// 查找管理员
	admin, err := c.adminDao.GetAdminByName(req.Username)
	if err != nil {
		return nil, code_msg.Unauthorized, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		return nil, code_msg.Unauthorized, err
	}

	// 生成JWT token（这里简化处理，实际应该使用JWT）
	// TODO: 使用JWT生成token
	token := "admin_token_" + req.Username // 临时token

	resp = &auth.LoginResponse{
		Data: &auth.LoginData{
			Token: token,
			Admin: &auth.AdminInfo{
				Id:   admin.ID,
				Name: admin.Name,
			},
		},
	}

	return resp, 0, nil
}
