package loader

import (
	"github.com/hdget/provider-config-viper/param"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"strings"
	"time"
)

type remoteConfigLoader struct {
	viper *viper.Viper
	param *param.Remote
}

func NewRemoteConfigLoader(viper *viper.Viper, param *param.Remote) Loader {
	return &remoteConfigLoader{
		viper: viper,
		param: param,
	}
}

// Load 从环境变量中读取配置信息
func (loader *remoteConfigLoader) Load() error {
	if loader.param == nil || loader.param.Provider == "" || len(loader.param.Endpoints) == 0 {
		return nil
	}

	var err error
	if loader.param.Secret != "" {
		err = loader.viper.AddSecureRemoteProvider(
			loader.param.Provider,
			strings.Join(loader.param.Endpoints, ";"),
			loader.param.WatchPath,
			loader.param.Secret,
		)
	} else {
		err = loader.viper.AddRemoteProvider(
			loader.param.Provider,
			strings.Join(loader.param.Endpoints, ";"),
			loader.param.WatchPath,
		)
	}
	if err != nil {
		return errors.Wrap(err, "add remote Provider")
	}

	loader.viper.SetConfigType(loader.param.RemoteConfigType)

	// 尝试读取，不报错
	_ = loader.viper.ReadRemoteConfig()

	// 自动读取到kvstore
	if err = loader.viper.WatchRemoteConfigOnChannel(); err != nil {
		return err
	}

	// open a goroutine to unmarshal remote config
	go func() {
		for {
			time.Sleep(time.Second * time.Duration(loader.param.WatchInterval))

			if loader.param.WatchCallback != nil {
				loader.param.WatchCallback()
			}
		}
	}()

	return nil
}
