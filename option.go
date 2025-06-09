package viper

import "github.com/hdget/provider-config-viper/param"

type Option func(*param.Param)

func WithConfigFile(configFile string) Option {
	return func(param *param.Param) {
		param.File.File = configFile
	}
}

func WithConfigContent(configContent []byte) Option {
	return func(param *param.Param) {
		param.Cli.Content = configContent
	}
}

func WithEnableRemote() Option {
	return func(p *param.Param) {
		p.Remote = param.NewRemoteDefaultParam()
	}
}

func WithRemote(remoteParam *param.Remote) Option {
	return func(p *param.Param) {
		p.Remote = remoteParam
	}
}
