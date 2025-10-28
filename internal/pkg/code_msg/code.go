package code_msg

var CodeMsg = map[BusinessCode]string{
	Success:                     "请求成功",
	BadRequest:                  "请求参数有误",
	Unauthorized:                "请登录后访问",
	Forbidden:                   "拒绝访问",
	NotFound:                    "服务器无法找到所请求的资源",
	MethodNotAllowed:            "禁止了使用当前 HTTP 方法的请求",
	NotAcceptable:               "服务器端无法提供与 Accept-Charset 以及 Accept-Language 消息头指定的值相匹配的响应",
	ProxyAuthenticationRequired: "缺乏位于浏览器与可以访问所请求资源的服务器之间的代理服务器",
	RequestTimeout:              "服务器想要将没有在使用的连接关闭",
	RequestConflict:             "服务器在完成请求时发生冲突",
	RequestDelete:               "如果请求的资源已永久删除，服务器就会返回此响应",
	RequestTooLarge:             "服务器无法处理请求，因为请求实体过大，超出服务器的处理能力",
	RequestTooLong:              "请求的URI过长，服务器无法处理",
	RequestToMany:               "请求过于频繁，请稍后再试",
	ServerError:                 "服务器开小差～",
	ServerTimeOut:               "请求超时",

	UserNameExisted: "用户名已存在",
}
