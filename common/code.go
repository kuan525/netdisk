package common

type ErrorCode int32

// 状态码从10000开始
const (
	_ int32 = iota + 9999
	// StatusOk 正常
	StatusOk
	// StatusParamInvalid 请求参数无效
	StatusParamInvalid
	// StatusServerError 服务出错
	StatusServerError
	// StatusRegisterFailed 注册失败
	StatusRegisterFailed
	// StatusLoginFailed 登陆失败
	StatusLoginFailed
	// StatusTokenInvalid token无效
	StatusTokenInvalid
	// StatusUserNotExists 用户不存在
	StatusUserNotExists
)
