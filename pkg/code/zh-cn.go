package code

var ZhCnText = map[int64]string{
	InternalServerError: "服务器内部错误",
	TooManyRequests:     "请求过多",
	ParamBindError:      "参数信息错误",
	AuthorizationError:  "签名信息错误",
	ResubmitError:       "请勿重复提交",
	SendEmailError:      "发送邮件失败",

	IncorrectUsernameOrPassword: "用户名或密码不正确",
	UsernameExist:               "用户名已存在",
	UserNotExist:                "用户不存在",
	AccountDisabled:             "帐户已禁用",
	UserNameEmpty:               "用户名不可为空",
	NoGameServerAvailable:       "无可用游戏服务器",
}
