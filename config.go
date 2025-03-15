package viper

import (
	"github.com/hdget/common/types"
	"path/filepath"
)

type viperConfig struct {
	app             string
	configEnvPrefix string   // 环境变量前缀
	configContent   []byte   // 如果用WithConfigContent指定了配置内容，则这里不为空
	configRootDirs  []string // 配置文件所在的RootDirs
	configType      string   // 配置内容类型，e,g: toml, json
	configFile      string   // 指定的配置文件
	searchDirs      []string // 未指定配置文件情况下，搜索的目录
	searchFileName  string   // 未指定配置文件情况下，搜索的文件名，不需要文件后缀
}

var (
	defaultConfigRootDirs = []string{
		".",                                      // current dir
		filepath.Join("config", "app"),           // default config root dir1
		filepath.Join("common", "config", "app"), // default config root dir2
	}
)

const (
	defaultEnvPrefix  = "HD"
	defaultConfigType = "toml"
)

func newConfig(sdkConfig *types.SdkConfig) *viperConfig {
	c := &viperConfig{
		app:             sdkConfig.App,
		configEnvPrefix: defaultEnvPrefix,
		configRootDirs:  defaultConfigRootDirs, // 其他环境的BaseDir
		configType:      defaultConfigType,
	}

	if sdkConfig.ConfigFilePath != "" {
		c.configFile = sdkConfig.ConfigFilePath
	}

	if len(sdkConfig.ConfigRootDirs) > 0 {
		c.configRootDirs = sdkConfig.ConfigRootDirs
	} else {
		c.configRootDirs = defaultConfigRootDirs
	}

	return c
}
