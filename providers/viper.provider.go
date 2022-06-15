// viper.Viper

package providers

import (
	"github.com/spf13/viper"
	"github.com/tulanz/pkg/errors"
)

func NewViperProvider() (*viper.Viper, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("./conf")
	viper.AddConfigPath("/etc/share/")

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "fatal error config")
	}
	return viper.GetViper(), nil
}
