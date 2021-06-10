package code

var (
	// 系统通用错误，范围：[5000, 5100)
	DBErr          = buildCode(5000, "系统错误01")
	RedisErr       = buildCode(5001, "系统错误02")
	ConvertErr     = buildCode(5010, "对象转换错误")
	RequestParamErr = buildCode(5011, "请求参数错误")
	RecordNotFound = buildCode(5044, "数据不存在或已被删除")

	UnknownErr     = buildCode(5099, "未知错误")

	// 与用户相关的错误，范围：[5100,5200)
	UsernameExistedErr  = buildCode(5100, "用户名已被使用")
	AuthErr             = buildCode(5101, "账号或密码错误")
	TokenExpiredErr     = buildCode(5102, "登录认证已过期")
	TokenNotValidYetErr = buildCode(5103, "账号未激活")
	TokenMalformedErr   = buildCode(5104, "非法令牌")
	TokenInvalidErr     = buildCode(5105, "无效令牌")

	// 与商品相关的错误，范围：[5300,5400)
	GoodsSaleOut    = buildCode(5200, "商品已售罄")
	SeckillNotStart = buildCode(5211, "秒杀还未开始")
	SeckillEnded    = buildCode(5222, "秒杀已结束")
)
