package response

/*
关键设计原则：
遵循HTTP语义：
	4xx 表示客户端问题（如 400/403/404）
	5xx 表示服务端问题（如 500/503）
	业务冲突用 409（如手机号重复）

避免过度自定义：
	如 EmptyErr（查询为空）不是错误，应返回 200 + 空数据
	UnknownErr 应归为 500

业务错误映射：
	权限问题统一用 403
	数据不存在用 404
	数据冲突用 409

扩展建议：
	go
	// 可补充更多场景
	TooManyRequests uint = 429 // 限流
	PaymentRequired uint = 402 // 支付场景
*/

type Codes struct {
	// #成功
	Success uint // 200 // RFC 9110, 15.3.1

	// #失败（通用客户端错误）
	Fail uint // 400 // RFC 9110, 15.5.1

	// #认证错误
	AuthErr uint // 401 // RFC 9110, 15.5.2

	// #服务器内部错误
	ServerErr uint // 500 // RFC 9110, 15.6.1

	// #未发现接口/资源
	NotFoundErr uint // 404 // RFC 9110, 15.5.5

	// #未知错误（建议用500代替）
	UnknownErr uint // 500 // 非标准，通常映射到500

	// #参数错误（细化400场景）
	ParamErr uint // 400 // RFC 9110, 15.5.1

	// #拒绝访问（权限不足）
	RejectErr uint // 403 // RFC 9110, 15.5.4

	// #请求方法错误
	MethodErr uint // 405 // RFC 9110, 15.5.6

	// #缺少参数（细化400场景）
	ParamLack uint // 400 // RFC 9110, 15.5.1

	// #业务上的错误（如用户已存在）
	UserExistsErr uint // 409 // RFC 9110, 15.5.10

	// #查询失败（数据库错误等）
	QueryErr uint // 500 // RFC 9110, 15.6.1

	// #查询为空（非错误，建议用200+空数据）
	EmptyErr uint // 200 // 非错误场景

	// #整顿期间（临时不可用）
	CleanUp uint // 503 // RFC 9110, 15.6.4

	// #未绑定手机号（业务逻辑错误）
	PhoneUnbind uint // 403 // RFC 9110, 15.5.4

	// #未验证手机号
	PhoneUncheck uint // 403 // RFC 9110, 15.5.4

	// #不支持邮箱登录
	EmailErr uint // 403 // RFC 9110, 15.5.4

	// #手机号已被使用（冲突）
	PhoneUsed uint // 409 // RFC 9110, 15.5.10

	// #创建失败（数据库错误）
	CreateErr uint // 500 // RFC 9110, 15.6.1

	// #用户不存在
	UserNotFound uint // 404 // RFC 9110, 15.5.5

	// #数据不存在
	DataNotExit uint // 404 // RFC 9110, 15.5.5

	// #更新失败
	UpdateErr uint // 500 // RFC 9110, 15.6.1

	// #密码错误
	PasswordErr uint

	// #验证码错误
	CheckCodeErr uint
}

var ApiCode = &Codes{
	Success:       200,
	Fail:          400,
	AuthErr:       401,
	ServerErr:     500,
	NotFoundErr:   404,
	UnknownErr:    500,
	ParamErr:      400,
	RejectErr:     403,
	MethodErr:     405,
	ParamLack:     400,
	UserExistsErr: 409,
	QueryErr:      500,
	EmptyErr:      200,
	CleanUp:       503,
	PhoneUnbind:   403,
	PhoneUncheck:  403,
	EmailErr:      403,
	PhoneUsed:     409,
	CreateErr:     500,
	UserNotFound:  404,
	DataNotExit:   404,
	UpdateErr:     500,
	PasswordErr:   400,
	CheckCodeErr:  400,
}

type Messages struct {
	// #成功
	Success string // 200 // RFC 9110, 15.3.1

	// #失败（通用客户端错误）
	Fail string // 400 // RFC 9110, 15.5.1

	// #认证错误
	AuthErr string // 401 // RFC 9110, 15.5.2

	// #服务器内部错误
	ServerErr string // 500 // RFC 9110, 15.6.1

	// #未发现接口/资源
	NotFoundErr string // 404 // RFC 9110, 15.5.5

	// #未知错误（建议用500代替）
	UnknownErr string // 500 // 非标准，通常映射到500

	// #参数错误（细化400场景）
	ParamErr string // 400 // RFC 9110, 15.5.1

	// #拒绝访问（权限不足）
	RejectErr string // 403 // RFC 9110, 15.5.4

	// #请求方法错误
	MethodErr string // 405 // RFC 9110, 15.5.6

	// #缺少参数（细化400场景）
	ParamLack string // 400 // RFC 9110, 15.5.1

	// #业务上的错误（如用户已存在）
	UserExistsErr string // 409 // RFC 9110, 15.5.10

	// #查询失败（数据库错误等）
	QueryErr string // 500 // RFC 9110, 15.6.1

	// #查询为空（非错误，建议用200+空数据）
	EmptyErr string // 200 // 非错误场景

	// #整顿期间（临时不可用）
	CleanUp string // 503 // RFC 9110, 15.6.4

	// #未绑定手机号（业务逻辑错误）
	PhoneUnbind string // 403 // RFC 9110, 15.5.4

	// #未验证手机号
	PhoneUncheck string // 403 // RFC 9110, 15.5.4

	// #不支持邮箱登录
	EmailErr string // 403 // RFC 9110, 15.5.4

	// #手机号已被使用（冲突）
	PhoneUsed string // 409 // RFC 9110, 15.5.10

	// #创建失败（数据库错误）
	CreateErr string // 500 // RFC 9110, 15.6.1

	// #用户不存在
	UserNotFound string // 404 // RFC 9110, 15.5.5

	// #数据不存在
	DataNotExit string // 404 // RFC 9110, 15.5.5

	// #更新失败
	UpdateErr string // 500 // RFC 9110, 15.6.1

	// #密码错误
	PasswordErr string

	// #验证码错误
	CheckCodeErr string
}

var ApiMsg = &Messages{
	Success:       "Success",
	Fail:          "Fail",
	AuthErr:       "AuthErr",
	ServerErr:     "ServerErr",
	NotFoundErr:   "NotFoundErr",
	UnknownErr:    "UnknownErr",
	ParamErr:      "ParamErr",
	RejectErr:     "RejectErr",
	MethodErr:     "MethodErr",
	ParamLack:     "ParamLack",
	UserExistsErr: "UserExistsErr",
	QueryErr:      "QueryErr",
	EmptyErr:      "EmptyErr",
	CleanUp:       "CleanUp",
	PhoneUnbind:   "PhoneUnbind",
	PhoneUncheck:  "PhoneUncheck",
	EmailErr:      "EmailErr",
	PhoneUsed:     "PhoneUsed",
	CreateErr:     "CreateErr",
	UserNotFound:  "UserNotFound",
	DataNotExit:   "DataNotExit",
	UpdateErr:     "UpdateErr",
	PasswordErr:   "PasswordErr",
	CheckCodeErr:  "CheckCodeErr",
}