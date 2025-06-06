package loader

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type remoteConfigLoader struct {
	remoteViper *viper.Viper
	option      *RemoteConfigLoaderOption
}

func NewRemoteConfigLoader(remoteViper *viper.Viper, option *RemoteConfigLoaderOption) Loader {
	return &remoteConfigLoader{
		remoteViper: remoteViper,
		option:      option,
	}
}

// Load 从环境变量中读取配置信息
func (loader *remoteConfigLoader) Load() error {
	if loader.option.RemoteProvider == "" || len(loader.option.RemoteEndpoints) == 0 {
		return nil
	}

	var err error
	if loader.option.RemoteSecret != "" {
		err = loader.remoteViper.AddSecureRemoteProvider(
			loader.option.RemoteProvider,
			strings.Join(loader.option.RemoteEndpoints, ";"),
			loader.option.RemotePath,
			loader.option.RemoteSecret,
		)
	} else {
		err = loader.remoteViper.AddRemoteProvider(
			loader.option.RemoteProvider,
			strings.Join(loader.option.RemoteEndpoints, ";"),
			loader.option.RemotePath,
		)
	}
	if err != nil {
		return errors.Wrap(err, "add remote Provider")
	}

	loader.remoteViper.SetConfigType(loader.option.RemoteConfigType)
	err = loader.remoteViper.ReadRemoteConfig()
	if err != nil {
		return errors.Wrap(err, "read remote provider")
	}

	// open a goroutine to watch remote changes forever
	go func() {
		for {
			time.Sleep(time.Second * time.Duration(loader.option.RemoteWatchInterval)) // delay after each request

			// currently, only tested with etcd support
			if err = loader.remoteViper.WatchRemoteConfig(); err != nil {
				continue
			}

			loader.option.RemoteWatchCallback()
		}
	}()

	return nil
}
