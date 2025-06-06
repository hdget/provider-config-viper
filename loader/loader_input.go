package loader

import (
	"bytes"
	"github.com/spf13/viper"
)

type InputConfigLoaderOption struct {
	Content []byte // 如果用WithConfigContent指定了配置内容，则这里不为空
}

type inputLoader struct {
	localViper *viper.Viper
	option     *InputConfigLoaderOption
}

func NewInputConfigLoader(localViper *viper.Viper, option *InputConfigLoaderOption) Loader {
	return &inputLoader{
		localViper: localViper,
		option:     option,
	}
}

// Load 从环境变量中读取配置信息
func (loader *inputLoader) Load() error {
	// 如果指定了配置内容，则合并
	if loader.option.Content != nil {
		_ = loader.localViper.MergeConfig(bytes.NewReader(loader.option.Content))
	}
	return nil
}
