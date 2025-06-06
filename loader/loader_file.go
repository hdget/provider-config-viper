package loader

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

type fileConfigLoader struct {
	app        string
	env        string
	localViper *viper.Viper
	option     *FileConfigLoaderOption
}

func NewFileConfigLoader(localViper *viper.Viper, option *FileConfigLoaderOption) Loader {
	return &fileConfigLoader{
		localViper: localViper,
		option:     option,
	}
}

func (f *fileConfigLoader) Load() error {
	// 找配置文件
	err := f.setupConfigFileFindOptions()
	if err != nil {
		return errors.Wrapf(err, "setup config file")
	}

	// 读取配置文件
	err = f.localViper.ReadInConfig()
	if err != nil {
		return errors.Wrapf(err, "read config file: %s", f.localViper.ConfigFileUsed())
	}

	return nil
}

func (f *fileConfigLoader) setupConfigFileFindOptions() error {
	// 设置文件配置的类型
	f.localViper.SetConfigType(f.option.FileConfigType)

	// 如果指定了配置文件
	if f.option.File != "" {
		f.localViper.SetConfigFile(f.option.File)
		return nil
	}

	// 未指定配置文件
	// 获取config filename
	searchConfigFileName := f.option.SearchFileName
	if searchConfigFileName == "" {
		searchConfigFileName = f.getDefaultConfigFilename()
	}

	// 获取config dirs
	searchConfigDirs := f.option.SearchDirs
	if len(searchConfigDirs) == 0 {
		foundDir := f.findConfigDir()
		if foundDir == "" {
			return fmt.Errorf("config dir not found, app: %s, env: %s", f.app, f.env)
		}
		searchConfigDirs = append(searchConfigDirs, foundDir)
	}

	// 设置搜索选项
	for _, dir := range searchConfigDirs {
		f.localViper.AddConfigPath(dir) // 指定目录
	}
	f.localViper.SetConfigName(searchConfigFileName)

	return nil
}

// getDefaultConfigFilename 缺省的配置文件名: <app>.<env>
func (f *fileConfigLoader) getDefaultConfigFilename() string {
	return strings.Join([]string{f.app, f.env}, ".")
}

// findConfigDirs 缺省的配置文件名: <app>.<env>
func (f *fileConfigLoader) findConfigDir() string {
	// iter to root directory
	absStartPath, err := filepath.Abs(".")
	if err != nil {
		return ""
	}

	var found string
	matchFile := fmt.Sprintf("%s.%s.%s", f.app, f.env, f.option.FileConfigType)
	currPath := absStartPath
LOOP:
	for {
		for _, rootDir := range f.option.RootDirs {
			// possible parent dir name
			dirName := filepath.Join(rootDir, f.app)
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
