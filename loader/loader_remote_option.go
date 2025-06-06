package loader

type RemoteConfigLoaderOption struct {
	RemoteProvider      string
	RemoteEndpoints     []string
	RemoteSecret        string
	RemoteConfigType    string
	RemoteWatchInterval int // 单位：秒
}

func NewRemoteConfigLoaderOption() *RemoteConfigLoaderOption {
	return &RemoteConfigLoaderOption{
		RemoteEndpoints:     []string{"http://127.0.0.1:2379"},
		RemoteProvider:      "etcd3",
		RemoteConfigType:    "json",
		RemoteWatchInterval: 5,
	}
}
