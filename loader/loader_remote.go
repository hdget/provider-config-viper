package loader

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"strings"
	"time"
)

type remoteConfigLoader struct {
	remoteViper         *viper.Viper
	option              *RemoteConfigLoaderOption
	remotePath          string
	remoteWatchCallback func()
}

func NewRemoteConfigLoader(remoteViper *viper.Viper, option *RemoteConfigLoaderOption, remotePath string, remoteWatchCallback func()) Loader {
	return &remoteConfigLoader{
		remoteViper:         remoteViper,
		option:              option,
		remotePath:          remotePath,
		remoteWatchCallback: remoteWatchCallback,
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
			loader.remotePath,
			loader.option.RemoteSecret,
		)
	} else {
		err = loader.remoteViper.AddRemoteProvider(
			loader.option.RemoteProvider,
			strings.Join(loader.option.RemoteEndpoints, ";"),
			loader.remotePath,
		)
	}
	if err != nil {
		return errors.Wrap(err, "add remote Provider")
	}

	loader.remoteViper.SetConfigType(loader.option.RemoteConfigType)

	// 尝试读取，不报错
	_ = loader.remoteViper.ReadRemoteConfig()

	if err = loader.remoteViper.WatchRemoteConfigOnChannel(); err != nil {
		return err
	}

	// open a goroutine to unmarshal remote config
	go func() {
		for {
			time.Sleep(time.Second * time.Duration(loader.option.RemoteWatchInterval))

			if loader.remoteWatchCallback != nil {
				loader.remoteWatchCallback()
			}
		}
	}()

	return nil
}
