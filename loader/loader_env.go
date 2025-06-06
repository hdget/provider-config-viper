package loader

import "github.com/spf13/viper"

type envLoader struct {
	localViper *viper.Viper
	option     *EnvConfigLoaderOption
}

func NewEnvConfigLoader(localViper *viper.Viper, option *EnvConfigLoaderOption) Loader {
	return &envLoader{
		localViper: localViper,
		option:     option,
	}
}

// Load 从环境变量中读取配置信息
func (loader *envLoader) Load() error {
	envPrefix := loader.option.EnvPrefix
	// 如果设置了环境变量前缀，则尝试自动获取环境变量中的配置
	if envPrefix == "" {
		envPrefix = defaultEnvPrefix
	}

	loader.localViper.SetEnvPrefix(envPrefix)
	loader.localViper.AutomaticEnv()
	return nil
}
