package code

var (
	// 系统通用错误，范围：[5000, 5100)
	DBErr        = buildCode(5000, "系统错误01")
	RedisErr     = buildCode(5010, "系统错误02")
	ConvertErr = buildCode(5050, "对象转换错误")
	UnknownErr   = buildCode(5099, "未知错误")

	// 与用户相关的错误，范围：[5100,5200)
	UsernameExistedErr = buildCode(5100, "用户名已被使用")
	AuthErr = buildCode(5101, "账号或密码错误")
	TokenExpiredErr = buildCode(5102, "登录认证已过期")
	TokenNotValidYetErr = buildCode(5103, "账号未激活")
	TokenMalformedErr = buildCode(5104, "非法令牌")
	TokenInvalidErr = buildCode(5105, "无效令牌")
)
