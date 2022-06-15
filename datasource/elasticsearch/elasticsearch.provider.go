package elasticsearch

import (
	"errors"

	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
)

func NewElasticSearchClient(vip *viper.Viper) (*elastic.Client, error) {

	if !vip.IsSet("elasticsearch") {
		return nil, errors.New("elasticsearch 没有配置")
	}
	addrs := vip.GetStringSlice("elasticsearch.addrs")
	var user, passwd string
	if vip.IsSet("elasticsearch.user") {
		user = vip.GetString("elasticsearch.user")
	}
	if vip.IsSet("elasticsearch.password") {
		passwd = vip.GetString("elasticsearch.password")
	}

	options := make([]elastic.ClientOptionFunc, 0)
	options = append(options, elastic.SetURL(addrs...))
	options = append(options, elastic.SetSniff(false))
	if user != "" {
		options = append(options, elastic.SetBasicAuth(user, passwd))
	}

	client, err := elastic.NewClient(options...)
	return client, err
}
