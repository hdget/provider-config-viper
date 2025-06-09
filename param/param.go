package param

type Param struct {
	*File
	*Env
	*Remote
	*Cli
}

func GetDefaultParam() *Param {
	return &Param{
		File:   NewFileDefaultParam(),
		Env:    NewEnvDefaultParam(),
		Cli:    NewCliDefaultParam(),
		Remote: nil, // 暂时禁用remote
	}
}
