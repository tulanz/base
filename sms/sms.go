package sms

type Sms interface {
	SendVerifyCode(phone string, code string) (string, error)
}
