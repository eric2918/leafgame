package code

import "leafgame/conf"

const (
	// 10101
	InternalServerError = 10100
	TooManyRequests     = 10101
	ParamBindError      = 10102
	AuthorizationError  = 10103
	ResubmitError       = 10104
	SendEmailError      = 10105

	// 登录成功
	//

	IncorrectUsernameOrPassword = 10200
	UsernameExist               = 10201
	UserNotExist                = 10202
	AccountDisabled             = 10203
	UserNameEmpty               = 20203
	NoGameServerAvailable       = 20204

	// 手机号已存在
	// 昵称已存在
	// 邮箱已存在
	// 验证码发送失败
	// 手机号验证失败
	// 验证邮件发送失败
	// 邮箱验证失败

)

func Text(code int64) string {
	lang := conf.Server.Language
	if lang == "zh-cn" {
		return ZhCnText[code]
	} else {
		return EnUsText[code]
	}
	return ""
}
