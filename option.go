package viper

import (
	"github.com/hdget/provider-config-viper/loader"
)

type Option struct {
	*loader.FileConfigLoaderOption
	*loader.EnvConfigLoaderOption
	*loader.RemoteConfigLoaderOption
	*loader.InputConfigLoaderOption
}

func NewOption() *Option {
	return &Option{
		FileConfigLoaderOption:   loader.NewFileConfigLoaderOption(),
		EnvConfigLoaderOption:    loader.NewEnvConfigLoaderOption(),
		RemoteConfigLoaderOption: loader.NewRemoteConfigLoaderOption(),
	}
}

func (o *Option) UseConfigFile(configFile string) *Option {
	o.FileConfigLoaderOption.File = configFile
	return o
}

func (o *Option) UseEnvPrefix(envPrefix string) *Option {
	o.EnvConfigLoaderOption.EnvPrefix = envPrefix
	return o
}

func (o *Option) UseRemotePath(remotePath string) *Option {
	o.RemoteConfigLoaderOption.RemotePath = remotePath
	return o
}

func (o *Option) UseRemoteWatchCallback(callback func()) *Option {
	o.RemoteConfigLoaderOption.RemoteWatchCallback = callback
	return o
}
