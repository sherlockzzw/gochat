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
	AccountLocked        BusinessCode = 4206 // 账户已被锁定
)

const (
	UserNameExisted      BusinessCode = 4301 // 用户名已存在
	LoginAccountExists   BusinessCode = 4302 // 登录账号已存在
	PhoneExists          BusinessCode = 4303 // 手机号已被注册
	EmailExists          BusinessCode = 4304 // 邮箱已被注册
	AgreementNotAccepted BusinessCode = 4305 // 未同意用户协议或隐私政策
	PhoneOrEmailRequired BusinessCode = 4306 // 必须提供手机号或邮箱
	CodeTypeMismatch     BusinessCode = 4307 // 验证码类型不匹配
	ParameterError       BusinessCode = 4001 // 参数错误
	CodeError            BusinessCode = 4002 // 验证码错误
	CodeExpired          BusinessCode = 4003 // 验证码已过期

	// 群组相关
	GroupNotFound        BusinessCode = 4401 // 群组不存在
	NotGroupMember       BusinessCode = 4402 // 不在群组中
	NotInGroup           BusinessCode = 4402 // 不在群组中（NotGroupMember的别名，使用相同错误码）
	NoPermission         BusinessCode = 4403 // 没有权限
	CannotRemoveOwner    BusinessCode = 4404 // 不能移除群主
	CannotRemoveAdmin    BusinessCode = 4405 // 管理员不能移除其他管理员
	UserNotInGroup       BusinessCode = 4406 // 用户不在群组中

	// 好友相关
	NotFriend            BusinessCode = 4501 // 不是好友关系
	FriendRelationNotExists BusinessCode = 4502 // 好友关系不存在
	FriendRequestNotExists BusinessCode = 4503 // 好友申请不存在
	AlreadyFriend        BusinessCode = 4504 // 已经是好友
	FriendRequestExists  BusinessCode = 4505 // 已发送过好友申请
	CannotAddSelf        BusinessCode = 4506 // 不能添加自己为好友
	InvalidOperation     BusinessCode = 4507 // 无效的操作类型
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
		AccountLocked:        "账户已被锁定",

		// 群组相关
		GroupNotFound:        "群组不存在",
		NotGroupMember:       "不在群组中",
		// NotInGroup 使用与 NotGroupMember 相同的错误码，不需要单独定义
		NoPermission:         "没有权限",
		CannotRemoveOwner:    "不能移除群主",
		CannotRemoveAdmin:    "管理员不能移除其他管理员",
		UserNotInGroup:       "用户不在群组中",

		// 好友相关
		NotFriend:            "不是好友关系",
		FriendRelationNotExists: "好友关系不存在",
		FriendRequestNotExists: "好友申请不存在",
		AlreadyFriend:        "已经是好友",
		FriendRequestExists:  "已发送过好友申请，请等待对方处理",
		CannotAddSelf:        "不能添加自己为好友",
		InvalidOperation:     "无效的操作类型",
	}

	if msg, ok := msgMap[code]; ok {
		return msg
	}
	return "未知错误"
}
