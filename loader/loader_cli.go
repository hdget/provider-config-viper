package loader

import (
	"bytes"
	"github.com/spf13/viper"
)

type CliConfigLoaderOption struct {
	Content []byte // 如果用WithConfigContent指定了配置内容，则这里不为空
}

type cliConfigLoader struct {
	localViper *viper.Viper
	option     *CliConfigLoaderOption
}

func NewCliConfigLoader(localViper *viper.Viper, option *CliConfigLoaderOption) Loader {
	return &cliConfigLoader{
		localViper: localViper,
		option:     option,
	}
}

func NewCliConfigLoaderOption() *CliConfigLoaderOption {
	return &CliConfigLoaderOption{}
}

// Load 从环境变量中读取配置信息
func (loader *cliConfigLoader) Load() error {
	// 如果指定了配置内容，则合并
	if loader.option.Content != nil {
		_ = loader.localViper.MergeConfig(bytes.NewReader(loader.option.Content))
	}
	return nil
}
