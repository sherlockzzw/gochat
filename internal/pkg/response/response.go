package response

import (
	"gochat/internal/pkg/code_msg"
	"gochat/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SvcRequest struct {
	logger    *zap.Logger
	IsDevMode bool
}

func NewSvcRequest() *SvcRequest {
	logger, _ := zap.NewProduction()
	return &SvcRequest{
		logger:    logger,
		IsDevMode: utils.IsDevMode(),
	}
}

type result struct {
	Code    code_msg.BusinessCode `json:"code"`
	Msg     string                `json:"msg"`
	Data    interface{}           `json:"data,omitempty"`
	TraceID string                `json:"trace_id,omitempty"`
	DeBug   string                `json:"debug,omitempty"`
}

// JsonSuccess 成功响应
func (r *SvcRequest) JsonSuccess(c *gin.Context, data any) {
	res := &result{
		Code: code_msg.Success,
		Msg:  r.GetMsg(code_msg.Success),
		Data: data,
	}
	res.TraceID = r.getTraceID(c)
	r.logger.Sugar().Infof("Success Response: %+v", res)
	c.JSON(http.StatusOK, res)
	c.Abort()
}

// JsonError 通用错误响应
func (r *SvcRequest) JsonError(c *gin.Context, err error, msg string) {
	res := &result{
		Code: code_msg.ServerError,
		Msg:  msg,
	}
	res.TraceID = r.getTraceID(c)
	if r.IsDevMode && err != nil {
		res.DeBug = err.Error()
	}
	r.logger.Sugar().Errorf("Error Response: %+v, Error: %v", res, err)
	c.JSON(http.StatusOK, res)
	c.Abort()
}

// JsonErrorAgility 灵活错误响应
func (r *SvcRequest) JsonErrorAgility(c *gin.Context, err error, code code_msg.BusinessCode, msg ...string) {
	res := &result{
		Code: code,
		Msg:  r.GetMsg(code),
	}
	res.TraceID = r.getTraceID(c)

	// 如果提供了自定义消息，使用自定义消息
	if len(msg) > 0 && msg[0] != "" {
		res.Msg = msg[0]
	}

	if r.IsDevMode && err != nil {
		res.DeBug = err.Error()
	}

	r.logger.Sugar().Errorf("Error Response: %+v, Error: %v", res, err)
	c.JSON(http.StatusOK, res)
	c.Abort()
}

// JsonErrorFixation 固定错误码响应
func (r *SvcRequest) JsonErrorFixation(c *gin.Context, code code_msg.BusinessCode) {
	res := &result{
		Code: code,
		Msg:  r.GetMsg(code),
	}
	res.TraceID = r.getTraceID(c)
	if r.IsDevMode {
		res.DeBug = r.GetMsg(code)
	}
	r.logger.Sugar().Errorf("Error Response: %+v", res)
	c.JSON(http.StatusOK, res)
	c.Abort()
}

// JsonBadRequest 400错误
func (r *SvcRequest) JsonBadRequest(c *gin.Context, msg string) {
	r.JsonErrorAgility(c, nil, code_msg.BadRequest, msg)
}

// JsonUnauthorized 401错误
func (r *SvcRequest) JsonUnauthorized(c *gin.Context, msg string) {
	r.JsonErrorAgility(c, nil, code_msg.Unauthorized, msg)
}

// JsonNotFound 404错误
func (r *SvcRequest) JsonNotFound(c *gin.Context, msg string) {
	r.JsonErrorAgility(c, nil, code_msg.NotFound, msg)
}

// JsonParamError 参数错误
func (r *SvcRequest) JsonParamError(c *gin.Context, msg string) {
	r.JsonErrorAgility(c, nil, code_msg.BadRequest, msg)
}

// JsonUserError 用户相关错误
func (r *SvcRequest) JsonUserError(c *gin.Context, code code_msg.BusinessCode, msg ...string) {
	r.JsonErrorAgility(c, nil, code, msg...)
}

// JsonTokenError Token相关错误
func (r *SvcRequest) JsonTokenError(c *gin.Context, code code_msg.BusinessCode, msg ...string) {
	r.JsonErrorAgility(c, nil, code, msg...)
}

// GetMsg 获取错误消息
func (r *SvcRequest) GetMsg(code code_msg.BusinessCode) string {
	return code_msg.GetMsg(code)
}

// getTraceID 获取或生成TraceID
func (r *SvcRequest) getTraceID(c *gin.Context) string {
	// 从请求头获取TraceID
	if traceID := c.GetHeader("X-Trace-ID"); traceID != "" {
		return traceID
	}
	// 从上下文获取TraceID
	if traceID, exists := c.Get("TraceID"); exists {
		if tid, ok := traceID.(string); ok {
			return tid
		}
	}
	// 生成新的TraceID
	return time.Now().Format("20060102150405") + "-" + utils.GenerateRandomString(8)
}

// SetTraceID 设置TraceID到上下文
func (r *SvcRequest) SetTraceID(c *gin.Context, traceID string) {
	c.Set("TraceID", traceID)
}
