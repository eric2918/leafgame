package code

var EnUsText = map[int64]string{
	InternalServerError: "Internal server error",
	TooManyRequests:     "Too many requests",
	ParamBindError:      "Parameter error",
	AuthorizationError:  "Authorization error",
	ResubmitError:       "Please do not submit repeatedly",
	SendEmailError:      "Failed to send mail",

	IncorrectUsernameOrPassword: "Incorrect username and password",
	UsernameExist:               "Username already exists",
	UserNotExist:                "User not exist",
	AccountDisabled:             "Account is disabled",
	NoGameServerAvailable:       "No game server available",
}
