// Package viper

package viper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hdget/common/intf"
	loader2 "github.com/hdget/provider-config-viper/loader"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// viperConfigProvider 命令行配置
type viperConfigProvider struct {
	app    string
	env    string
	local  *viper.Viper
	remote *viper.Viper
	option *Option
}

const (
	remoteLocalKey  = "local.%s.%s"
	remoteGlobalKey = "global.%s"
)

// New 初始化config provider
func New(app, env string, option *Option) (intf.ConfigProvider, error) {
	provider := &viperConfigProvider{
		app:    app,
		env:    env,
		local:  viper.New(),
		remote: viper.New(),
		option: option,
	}

	err := provider.Load()
	if err != nil {
		return nil, errors.Wrap(err, "load local option")
	}

	return provider, nil
}

// Load 从各个配置源获取配置数据, 并加载到configVar中，同名变量优先级高的覆盖低的
// - remote: kvstore配置（低）
// - configFile: 文件配置(中)
// - env: 环境变量配置(高)
// - input: 命令行参数配置(最高)
func (p *viperConfigProvider) Load() error {
	// 如果环境变量为空，则加载最小基本配置
	if p.env == "" {
		return loader2.NewMinimalConfigLoader(p.app, p.local).Load()
	}

	if err := loader2.NewRemoteConfigLoader(p.remote, p.option.RemoteConfigLoaderOption).Load(); err != nil {
		return errors.Wrap(err, "load config from remote")
	}

	// 尝试从环境变量中获取配置信息
	if err := loader2.NewFileConfigLoader(p.local, p.option.FileConfigLoaderOption).Load(); err != nil {
		return errors.Wrap(err, "load config from file")
	}

	// 尝试从环境变量中获取配置信息
	if err := loader2.NewEnvConfigLoader(p.local, p.option.EnvConfigLoaderOption).Load(); err != nil {
		return errors.Wrap(err, "load config from env")
	}

	// 尝试从环境变量中获取配置信息
	if err := loader2.NewInputConfigLoader(p.local, p.option.InputConfigLoaderOption).Load(); err != nil {
		return errors.Wrap(err, "load config from cli")
	}

	return nil
}

// Unmarshal local
func (p *viperConfigProvider) Unmarshal(configVar any, args ...string) error {
	if err := p.mergeRemoteLocalConfig(); err != nil {
		return errors.Wrap(err, "merge remote local config")
	}

	if len(args) > 0 {
		return p.local.UnmarshalKey(args[0], configVar)
	}
	return p.local.Unmarshal(configVar)
}

func (p *viperConfigProvider) mergeRemoteLocalConfig() error {
	localPath := fmt.Sprintf(remoteLocalKey, p.env, p.app)
	if p.remote.Get(localPath) != nil {
		data, err := json.Marshal(p.remote.Sub(localPath).AllSettings())
		if err != nil {
			return err
		}

		err = p.local.MergeConfig(bytes.NewReader(data))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *viperConfigProvider) UnmarshalGlobal(configVar any, args ...string) error {
	globalPath := fmt.Sprintf(remoteGlobalKey, p.env)
	if len(args) > 0 {
		return p.remote.UnmarshalKey(globalPath, configVar)
	}
	return p.local.Unmarshal(configVar)
}
