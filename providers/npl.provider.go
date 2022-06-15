package providers

import (
	"github.com/spf13/viper"
	"github.com/tulanz/base/nlp"
	"github.com/tulanz/base/nlp/tencent"
	"go.uber.org/zap"
)

// nlp.Lexer{}.Default("百度是一家高科技公司")

func NewNlpProvider(vip *viper.Viper, logger *zap.Logger) (nlp.Summary, error) {
	if !vip.IsSet("nlp") {
		return &nlp.Default{}, nil
	}
	config := vip.Sub("nlp")
	driver := config.GetString("driver")

	switch driver {
	case "tencent":
		region := config.GetString("region")
		secretId := config.GetString("secretId")
		secretKey := config.GetString("secretKey")
		return tencent.NewTencentAI(secretId, secretKey, region)
	default:
		return &nlp.Default{}, nil
	}
}
