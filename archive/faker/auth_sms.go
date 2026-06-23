package faker

import (
	"time"
)

type authSms struct {
	Message             string              `json:"message"`
	SmsVerificationCode smsVerificationCode `json:"smsVerificationCode"`
}

type smsVerificationCode struct {
	Sended  bool      `json:"sended"`
	Expires time.Time `json:"expires"`
}

func GetAuthSms() *authSms {
	return &authSms{
		Message: "sms.verification_required",
		SmsVerificationCode: smsVerificationCode{
			Sended:  true,
			Expires: time.Now().Add(time.Minute * 1),
		},
	}
}
