package code_msg

type BusinessCode uint32

const (
	Success                     BusinessCode = 200 // 成功
	BadRequest                  BusinessCode = 400 //
	Unauthorized                BusinessCode = 401 //
	Forbidden                   BusinessCode = 403 //
	NotFound                    BusinessCode = 404 //
	MethodNotAllowed            BusinessCode = 405 //
	NotAcceptable               BusinessCode = 406 //
	ProxyAuthenticationRequired BusinessCode = 407 //
	RequestTimeout              BusinessCode = 408 //
	RequestConflict             BusinessCode = 409 //
	RequestDelete               BusinessCode = 410 //
	RequestTooLarge             BusinessCode = 413 //
	RequestTooLong              BusinessCode = 414 //
	RequestToMany               BusinessCode = 429 //
	ServerError                 BusinessCode = 500 // 服务器错误
	ServerTimeOut               BusinessCode = 504 // 服务器超时
)

// 身份校验
const (
	UserCheckTokenIllegal          = 4001 // 身份信息不合法
	UserCheckTokenSignatureInvalid = 4002 // 签名无效
	UserCheckTokenStale            = 4003 // 身份信息已过期，请重新登陆

	AdminCheckNotApiPermission = 4101 // 没有该接口权限

	UserRoleHarderNotSet              = 4200
	InvalidMobilePhone   BusinessCode = 4201
	UserNotExists        BusinessCode = 4202
	InvalidSmsCode       BusinessCode = 4203
	InvalidToken         BusinessCode = 4204
	PasswordError        BusinessCode = 4205
)

const (
	UserNameExisted BusinessCode = 4301
)

// GetMsg 获取错误消息
func GetMsg(code BusinessCode) string {
	msgMap := map[BusinessCode]string{
		Success:                     "成功",
		BadRequest:                  "请求参数错误",
		Unauthorized:                "未授权",
		Forbidden:                   "禁止访问",
		NotFound:                    "资源不存在",
		MethodNotAllowed:            "方法不允许",
		NotAcceptable:               "不可接受",
		ProxyAuthenticationRequired: "需要代理认证",
		RequestTimeout:              "请求超时",
		RequestConflict:             "请求冲突",
		RequestDelete:               "请求删除",
		RequestTooLarge:             "请求体过大",
		RequestTooLong:              "请求过长",
		RequestToMany:               "请求过于频繁",
		ServerError:                 "服务器内部错误",
		ServerTimeOut:               "服务器超时",

		// 身份校验
		UserCheckTokenIllegal:          "身份信息不合法",
		UserCheckTokenSignatureInvalid: "签名无效",
		UserCheckTokenStale:            "身份信息已过期，请重新登陆",
		AdminCheckNotApiPermission:     "没有该接口权限",

		// 用户相关
		UserRoleHarderNotSet: "用户角色未设置",
		InvalidMobilePhone:   "手机号格式错误",
		UserNotExists:        "用户不存在",
		InvalidSmsCode:       "短信验证码错误",
		InvalidToken:         "Token无效",
		UserNameExisted:      "用户名已存在",
		PasswordError:        "密码错误",
	}

	if msg, ok := msgMap[code]; ok {
		return msg
	}
	return "未知错误"
}
