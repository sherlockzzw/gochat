package analysis

import (
	"gochat/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func BindParameter[T any](ctx *gin.Context, resp *response.SvcRequest) (req *T, err error) {
	req = new(T)

	if err = ctx.ShouldBindJSON(req); err != nil {
		resp.JsonBadRequest(ctx, "参数绑定失败: "+err.Error())
		return nil, err
	}

	// 验证参数
	validate := validator.New()
	if err = validate.Struct(req); err != nil {
		resp.JsonBadRequest(ctx, "参数验证失败: "+err.Error())
		return nil, err
	}

	return req, nil
}

func BindQuery[T any](ctx *gin.Context, resp *response.SvcRequest) (req *T, err error) {
	req = new(T)

	if err = ctx.ShouldBindQuery(req); err != nil {
		resp.JsonBadRequest(ctx, "参数绑定失败: "+err.Error())
		return nil, err
	}

	// 验证参数
	validate := validator.New()
	if err = validate.Struct(req); err != nil {
		resp.JsonBadRequest(ctx, "参数验证失败: "+err.Error())
		return nil, err
	}

	return req, nil
}
