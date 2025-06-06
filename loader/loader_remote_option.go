package loader

type RemoteConfigLoaderOption struct {
	RemoteProvider      string
	RemoteEndpoints     []string
	RemoteSecret        string
	RemotePath          string
	RemoteConfigType    string
	RemoteWatchInterval int // 单位：秒
	RemoteWatchCallback func()
}

const (
	defaultRemotePath = "/config"
)

func NewRemoteConfigLoaderOption() *RemoteConfigLoaderOption {
	return &RemoteConfigLoaderOption{
		RemoteEndpoints:  make([]string, 0),
		RemotePath:       defaultRemotePath,
		RemoteConfigType: "json",
	}
}
