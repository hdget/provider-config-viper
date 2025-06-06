package loader

type EnvConfigLoaderOption struct {
	EnvPrefix string // 环境变量前缀
}

const (
	defaultEnvPrefix = "HD"
)

func NewEnvConfigLoaderOption() *EnvConfigLoaderOption {
	return &EnvConfigLoaderOption{
		EnvPrefix: defaultEnvPrefix,
	}
}
