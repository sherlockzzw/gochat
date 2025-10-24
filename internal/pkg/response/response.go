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
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"trace_id,omitempty"`
	DeBug   string      `json:"debug,omitempty"`
}

func (r *SvcRequest) JsonSuccess(c *gin.Context, data any) {
	res := &result{
		Code: code_msg.Success,
		Msg:  r.GetMsg(code_msg.Success),
		Data: data,
	}
	res.TraceID = r.getTraceID(c)
	r.logger.Sugar().Infof("%+v", res)
	c.JSON(http.StatusOK, res)
	c.Abort()
}

func (r *SvcRequest) JsonError(c *gin.Context, err error, msg string) {
	res := &result{
		Code: code_msg.ServerError,
		Msg:  msg,
	}
	res.TraceID = r.getTraceID(c)
	if r.IsDevMode && err != nil {
		res.DeBug = err.Error()
	}
	r.logger.Sugar().Infof("%+v", res)
	c.JSON(http.StatusOK, res)
	c.Abort()
}

func (r *SvcRequest) JsonBadRequest(c *gin.Context, msg string) {
	res := &result{
		Code: code_msg.BadRequest,
		Msg:  msg,
	}
	res.TraceID = r.getTraceID(c)
	r.logger.Sugar().Infof("%+v", res)
	c.JSON(http.StatusOK, res)
	c.Abort()
}

func (r *SvcRequest) JsonUnauthorized(c *gin.Context, msg string) {
	res := &result{
		Code: code_msg.Unauthorized,
		Msg:  msg,
	}
	res.TraceID = r.getTraceID(c)
	r.logger.Sugar().Infof("%+v", res)
	c.JSON(http.StatusOK, res)
	c.Abort()
}

func (r *SvcRequest) JsonNotFound(c *gin.Context, msg string) {
	res := &result{
		Code: code_msg.NotFound,
		Msg:  msg,
	}
	res.TraceID = r.getTraceID(c)
	r.logger.Sugar().Infof("%+v", res)
	c.JSON(http.StatusOK, res)
	c.Abort()
}

func (r *SvcRequest) JsonParamError(c *gin.Context, msg string) {
	res := &result{
		Code: code_msg.ParamError,
		Msg:  msg,
	}
	res.TraceID = r.getTraceID(c)
	r.logger.Sugar().Infof("%+v", res)
	c.JSON(http.StatusOK, res)
	c.Abort()
}

func (r *SvcRequest) GetMsg(code int) string {
	return code_msg.GetMsg(code)
}

func (r *SvcRequest) getTraceID(c *gin.Context) string {
	// 从请求头或上下文获取TraceID
	if traceID := c.GetHeader("X-Trace-ID"); traceID != "" {
		return traceID
	}
	// 如果没有TraceID，生成一个简单的
	return time.Now().Format("20060102150405") + "-" + utils.GenerateRandomString(8)
}
