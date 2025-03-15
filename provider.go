// Package config
// the default setting hierarchy looks like below:
//
//	...
//	setting/app/<app>/<app>.test.toml
//	setting/dapr/*
//	...
package viper

import (
	"bytes"
	"fmt"
	"github.com/hdget/common/intf"
	"github.com/hdget/common/types"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

// viperConfigProvider 命令行配置
type viperConfigProvider struct {
	env    string
	local  *viper.Viper
	config *viperConfig
	// 配置项

}

func (p *viperConfigProvider) GetCapability() types.Capability {
	return Capability
}

const (
	// 最小化的配置,保证日志工作正常
	tplMinimalConfigContent = `
[sdk]
  [sdk.logger]
	   level = "debug"
	   filename = "%s.log"
	   [sdk.logger.rotate]
		   max_age = 7`
)

// New 初始化config provider
func New(sdkConfig *types.SdkConfig) (intf.ConfigProvider, error) {
	env, err := getCurrentEnv()
	if err != nil {
		return nil, err
	}

	provider := &viperConfigProvider{
		local:  viper.New(),
		env:    env,
		config: newConfig(sdkConfig),
	}

	err = provider.loadLocal()
	if err != nil {
		return nil, errors.Wrap(err, "load local config")
	}

	return provider, nil
}

func (p *viperConfigProvider) Unmarshal(configVar any, args ...string) error {
	if len(args) > 0 {
		return p.local.UnmarshalKey(args[0], configVar)
	}
	return p.local.Unmarshal(configVar)
}

func (p *viperConfigProvider) GetApp() string {
	return p.config.app
}

// ///////////////////////////////////////////////////////////////
// private functions
// //////////////////////////////////////////////////////////////

// Load 从各个配置源获取配置数据, 并加载到configVar中， 同名变量配置高的覆盖低的
// - remote: kvstore配置（低）
// - configFile: 文件配置(中）
// - env: 环境变量配置(高)
func (p *viperConfigProvider) loadLocal() error {
	// 必须设置config的类型
	p.local.SetConfigType(p.config.configType)

	// 如果指定了配置内容，则合并
	if p.config.configContent != nil {
		_ = p.local.MergeConfig(bytes.NewReader(p.config.configContent))
	}

	// 如果环境变量为空，则加载最小基本配置
	if p.env == "" {
		return p.loadMinimal()
	}

	// 尝试从环境变量中获取配置信息
	p.loadFromEnv()

	// 尝试从配置文件中获取配置信息
	return p.loadFromFile()
}

// loadFromEnv 从环境文件中读取配置信息
func (p *viperConfigProvider) loadFromEnv() {
	// 如果设置了环境变量前缀，则尝试自动获取环境变量中的配置
	if p.config.configEnvPrefix != "" {
		p.local.SetEnvPrefix(p.config.configEnvPrefix)
		p.local.AutomaticEnv()
	}
}

func (p *viperConfigProvider) loadMinimal() error {
	minimalConfig := fmt.Sprintf(tplMinimalConfigContent, p.config.app)
	return p.local.MergeConfig(bytes.NewReader([]byte(minimalConfig)))
}

func (p *viperConfigProvider) loadFromFile() error {
	// 找配置文件
	err := p.setupConfigFile()
	if err != nil {
		return errors.Wrapf(err, "setup config file")
	}

	// 读取配置文件
	err = p.local.ReadInConfig()
	if err != nil {
		return errors.Wrapf(err, "read config file: %s", p.local.ConfigFileUsed())
	}

	return nil
}

func (p *viperConfigProvider) setupConfigFile() error {
	// 如果指定了配置文件
	if p.config.configFile != "" {
		p.local.SetConfigFile(p.config.configFile)
		return nil
	}

	// 未指定配置文件
	// 获取config filename
	searchConfigFileName := p.config.searchFileName
	if searchConfigFileName == "" {
		searchConfigFileName = p.getDefaultConfigFilename()
	}

	// 获取config dirs
	searchConfigDirs := p.config.searchDirs
	if len(searchConfigDirs) == 0 {
		foundDir := p.findConfigDir()
		if foundDir == "" {
			return fmt.Errorf("config dir not found, app: %s, env: %s", p.config.app, p.env)
		}
		searchConfigDirs = append(searchConfigDirs, foundDir)
	}

	// 设置搜索选项
	for _, dir := range searchConfigDirs {
		p.local.AddConfigPath(dir) // 指定目录
	}
	p.local.SetConfigName(searchConfigFileName)

	return nil
}

// getDefaultConfigFilename 缺省的配置文件名: <app>.<env>
func (p *viperConfigProvider) getDefaultConfigFilename() string {
	return strings.Join([]string{p.config.app, p.env}, ".")
}

// findConfigDirs 缺省的配置文件名: <app>.<env>
func (p *viperConfigProvider) findConfigDir() string {
	// iter to root directory
	absStartPath, err := filepath.Abs(".")
	if err != nil {
		return ""
	}

	var found string
	matchFile := fmt.Sprintf("%s.%s.%s", p.config.app, p.env, p.config.configType)
	currPath := absStartPath
LOOP:
	for {
		for _, rootDir := range p.config.configRootDirs {
			// possible parent dir name
			dirName := filepath.Join(rootDir, p.config.app)
			checkDir := filepath.Join(currPath, dirName, matchFile)
			matches, err := filepath.Glob(checkDir)
			if err == nil && len(matches) > 0 {
				found = filepath.Join(currPath, dirName)
				break LOOP
			}
		}

		// If we're already at the root, stop finding
		// windows has the driver name, so it need use TrimRight to test
		abs, _ := filepath.Abs(currPath)
		if abs == string(filepath.Separator) || len(strings.TrimRight(currPath, string(filepath.Separator))) <= 3 {
			break
		}

		// else, get parent dir
		currPath = filepath.Dir(currPath)
	}

	return found
}

func getCurrentEnv() (string, error) {
	if v, exists := os.LookupEnv("HD_ENV"); exists {
		return v, nil
	}
	return "", errors.New("env not found")
}
