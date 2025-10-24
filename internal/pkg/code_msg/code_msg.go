package code_msg

const (
	Success = 200

	// 客户端错误 4xx
	BadRequest          = 400
	Unauthorized        = 401
	Forbidden           = 403
	NotFound            = 404
	MethodNotAllowed    = 405
	Conflict            = 409
	UnprocessableEntity = 422
	TooManyRequests     = 429

	// 服务器错误 5xx
	ServerError        = 500
	NotImplemented     = 501
	BadGateway         = 502
	ServiceUnavailable = 503
	GatewayTimeout     = 504

	// 业务错误码
	UserNotFound      = 1001
	UserAlreadyExists = 1002
	PasswordError     = 1003
	TokenExpired      = 1004
	TokenInvalid      = 1005
	PermissionDenied  = 1006
	ParamError        = 1007
	DatabaseError     = 1008
)

var codeMsgMap = map[int]string{
	Success: "success",

	// 客户端错误
	BadRequest:          "请求参数错误",
	Unauthorized:        "未授权",
	Forbidden:           "禁止访问",
	NotFound:            "资源不存在",
	MethodNotAllowed:    "方法不允许",
	Conflict:            "资源冲突",
	UnprocessableEntity: "请求参数验证失败",
	TooManyRequests:     "请求过于频繁",

	// 服务器错误
	ServerError:        "服务器内部错误",
	NotImplemented:     "功能未实现",
	BadGateway:         "网关错误",
	ServiceUnavailable: "服务不可用",
	GatewayTimeout:     "网关超时",

	// 业务错误码
	UserNotFound:      "用户不存在",
	UserAlreadyExists: "用户已存在",
	PasswordError:     "密码错误",
	TokenExpired:      "Token已过期",
	TokenInvalid:      "Token无效",
	PermissionDenied:  "权限不足",
	ParamError:        "参数错误",
	DatabaseError:     "数据库错误",
}

func GetMsg(code int) string {
	if msg, ok := codeMsgMap[code]; ok {
		return msg
	}
	return "未知错误"
}
