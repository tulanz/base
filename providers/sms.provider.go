package providers

import (
	"errors"

	"github.com/spf13/viper"
	"github.com/tulanz/base/sms"
	"github.com/tulanz/base/sms/aliyun"
)

func NewSMSProvider(vip *viper.Viper) (sms.Sms, error) {

	if !vip.IsSet("sms") {
		return nil, errors.New("")
	}
	config := vip.Sub("sms")
	driver := config.Get("driver")

	var c sms.Sms
	switch driver {
	case "aliyun":
		regionId := config.GetString("regionId")
		accessKeyId := config.GetString("accessKeyId")
		accessKeySecret := config.GetString("accessKeySecret")
		signName := config.GetString("signName")
		templateCode := config.GetString("templateCode")
		c = aliyun.NewSms(regionId, accessKeyId, accessKeySecret, signName, templateCode)
	default:
		return nil, errors.New("不支持的短信网关")
	}

	return c, nil
}
