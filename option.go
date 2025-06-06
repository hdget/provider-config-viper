package viper

import (
	"github.com/hdget/provider-config-viper/loader"
)

type Option struct {
	File                      *loader.FileConfigLoaderOption
	Env                       *loader.EnvConfigLoaderOption
	Remote                    *loader.RemoteConfigLoaderOption
	Cli                       *loader.CliConfigLoaderOption
	RemoteGlobalWatchCallback func()
	RemoteAppWatchCallback    func()
}

func NewOption() *Option {
	return &Option{
		File:   loader.NewFileConfigLoaderOption(),
		Env:    loader.NewEnvConfigLoaderOption(),
		Remote: loader.NewRemoteConfigLoaderOption(),
		Cli:    loader.NewCliConfigLoaderOption(),
	}
}

func (o *Option) UseConfigFile(configFile string) *Option {
	o.File.File = configFile
	return o
}

func (o *Option) UseConfigContent(configContent []byte) *Option {
	o.Cli.Content = configContent
	return o
}

func (o *Option) UseEnvPrefix(envPrefix string) *Option {
	o.Env.EnvPrefix = envPrefix
	return o
}

func (o *Option) UseRemoteGlobalWatchCallback(callback func()) *Option {
	o.RemoteGlobalWatchCallback = callback
	return o
}

func (o *Option) UseRemoteAppWatchCallback(callback func()) *Option {
	o.RemoteAppWatchCallback = callback
	return o
}
