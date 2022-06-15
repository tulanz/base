package aliyun

import (
	"encoding/json"
	"errors"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/tulanz/base/sms"
)

type aliyunSms struct {
	client       *dysmsapi.Client
	RegionId     string
	SignName     string
	TemplateCode string
}

func NewSms(regionId, accessKeyId, accessKeySecret, signName, templateCode string) sms.Sms {
	client, _ := dysmsapi.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)
	return &aliyunSms{
		client:       client,
		SignName:     signName,
		RegionId:     regionId,
		TemplateCode: templateCode,
	}
}

func (a *aliyunSms) SendVerifyCode(phone string, code string) (string, error) {
	request := dysmsapi.CreateSendSmsRequest()
	request.Method = "POST"
	request.Scheme = "https"
	request.Domain = "dysmsapi.aliyuncs.com"
	request.SetVersion("2017-05-25")
	request.RegionId = a.RegionId
	request.SignName = a.SignName
	request.TemplateCode = a.TemplateCode

	param := map[string]string{}
	param["code"] = code
	b, _ := json.Marshal(param)
	request.PhoneNumbers = phone
	request.TemplateParam = string(b)

	response, err := a.client.SendSms(request)
	if err != nil {
		return "", err
	} else if response.Code != "OK" {
		return "", errors.New(response.Message)
	}
	return response.BizId, nil
}
