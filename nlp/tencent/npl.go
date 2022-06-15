package tencent

import (
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tnlp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/nlp/v20190408"
	"github.com/tulanz/base/nlp"
	"github.com/tulanz/base/utils/simple"
)

type TencentAI struct {
	client *tnlp.Client
}

func NewTencentAI(SecretId, SecretKey, region string) (nlp.Summary, error) {
	credential := common.NewCredential(SecretId, SecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "nlp.tencentcloudapi.com"

	client, err := tnlp.NewClient(credential, region, cpf)
	if err != nil {
		return nil, err
	}
	return &TencentAI{client: client}, nil
}

func (t *TencentAI) Default(title, text string, length uint64) (string, error) {

	text = simple.Substr(text, 0, 2000) // 腾讯最多不超过2000字

	request := tnlp.NewAutoSummarizationRequest()
	request.Text = common.StringPtr(text)
	request.Length = common.Uint64Ptr(length)
	response, err := t.client.AutoSummarization(request)

	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return "", err
	}
	if err != nil {
		return "", err
	}
	return *response.Response.Summary, nil
}
