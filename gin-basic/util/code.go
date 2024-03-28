package util

type Codes struct {
	// #成功
	Success uint // 200
	// #失败
	Fail uint // 400

	//# 认证错误
	AuthError uint // 401
	//MSG_AUTH_ERROR = 'token认证失败, 请重新登录'

	//# 服务器内部错误，状态码500
	ServerError uint // = 500
	//MSG_SERVER_ERROR = '网络操作失败，请稍后重试'

	//# 未发现接口
	NotFoundError uint // 404
	//MSG_NOT_FOUND_ERROR = '服务器没有此接口'

	//# 未知错误
	UnknownError uint // 405
	//MSG_UNKNOWN_ERROR = '未知错误'

	//# 参数错误
	ParamError uint // 420
	//MSG_PARAMETER_ERROR = '参数错误'

	//# 缺少参数
	ParamLack uint // 421
	//MSG_PARAMETER_LACK = '缺少参数'

	//# 查询失败
	QueryError uint //423
	//MSG_QUERY_ERROR = '查询失败'

	//# 查询为空
	EmptyError uint // 424
	//MSG_EMPTY_ERROR = '查询为空'

	//# 整顿期间
	CleanUp uint // 426
	//MSG_CLEAN_UP = '整顿期间'

	//# 未绑定手机号
	PhoneUnbind uint // 430
	//MSG_PHONE_UNBIND = '未绑定手机号'

	//# 手机验证
	PhoneUncheck uint //431
	//MSG_PHONE_UNCHECK = '未验证手机号'

	//# 不支持邮箱登录
	EmailError uint // 432
	//MSG_EMAIL_ERROR = '暂不支持邮箱登录'

	PhoneUsed uint // 433
	//MSG_PHONE_USED = '该手机号已被使用'

	CreateErr   uint

	//# 用户已存在
	UserExistsError uint // 440
	//MSG_BUSSINESS_ERROR = '用户已存在'

	// 用户不存在
	UserNotFont uint // 441
}

var ApiCode = &Codes{
	Success:         200,
	Fail:            400,
	AuthError:       401,
	ServerError:     500,
	NotFoundError:   404,
	UnknownError:    405,
	ParamError:      420,
	ParamLack:       421,
	UserExistsError: 4005,
	QueryError:      4006,
	EmptyError:      4007,
	CleanUp:         4008,
	PhoneUnbind:     4009,
	PhoneUncheck:    4010,
	EmailError:      4011,
	PhoneUsed:       4012,
	CreateErr:       4013,
	UserNotFont:     4014,
}

type Messages struct {
	// #成功
	Success string // 200

	// #失败
	Fail string // 300

	//# 认证错误
	AuthError string // 401
	//MSG_AUTH_ERROR = 'token认证失败, 请重新登录'

	//# 服务器内部错误，状态码500
	ServerError string // = 500
	//MSG_SERVER_ERROR = '网络操作失败，请稍后重试'

	//# 未发现接口
	NotFoundError string // 404
	//MSG_NOT_FOUND_ERROR = '服务器没有此接口'

	//# 未知错误
	UnknownError string // 405
	//MSG_UNKNOWN_ERROR = '未知错误'

	//# 参数错误
	ParamError string // 201
	//MSG_PARAMETER_ERROR = '参数错误'

	//# 拒绝访问
	RejectError string // 202
	//MSG_REJECT_ERROR = '拒绝访问'

	//# 拒绝访问
	MethodError string // 203
	//MSG_METHOD_ERROR = '请求方法错误'

	//# 缺少参数
	ParamLack string // 204
	//MSG_PARAMETER_LACK = '缺少参数'

	//# 业务上的错误
	UserExistsError string // 205
	//MSG_BUSSINESS_ERROR = '用户已存在'

	//# 查询失败
	QueryError string //206
	//MSG_QUERY_ERROR = '查询失败'

	//# 查询为空
	EmptyError string // 207
	//MSG_EMPTY_ERROR = '查询为空'

	//# 整顿期间
	CleanUp string // 208
	//MSG_CLEAN_UP = '整顿期间'

	//# 未绑定手机号
	PhoneUnbind string // 209
	//MSG_PHONE_UNBIND = '未绑定手机号'

	//# 手机验证
	PhoneUncheck string //210
	//MSG_PHONE_UNCHECK = '未验证手机号'

	//# 不支持邮箱登录
	EmailError string // 211
	//MSG_EMAIL_ERROR = '暂不支持邮箱登录'

	PhoneUsed string // 212
	//MSG_PHONE_USED = '该手机号已被使用'
	CreateErr string
	// MSG_CREATE_ERR = '创建失败'
	UserNotFound string
}

var ApiMessage = &Messages{
	Success:         "成功",
	Fail:            "失败",
	AuthError:       "token认证失败, 请重新登录",
	ServerError:     "网络操作失败，请稍后重试",
	NotFoundError:   "服务器没有此接口",
	UnknownError:    "未知错误",
	ParamError:      "参数错误",
	RejectError:     "拒绝访问",
	MethodError:     "请求方法错误",
	ParamLack:       "缺少参数",
	UserExistsError: "用户已存在",
	QueryError:      "查询失败",
	EmptyError:      "查询为空",
	CleanUp:         "整顿期间",
	PhoneUnbind:     "未绑定手机号",
	PhoneUncheck:    "未验证手机号",
	EmailError:      "暂不支持邮箱登录",
	PhoneUsed:       "该手机号已被使用",
	CreateErr:       "创建失败",
	UserNotFound:    "用户不存在",
}
