package registry

// 服务元数据信息
type EndpointUnit struct {
	Name string `json:"name"`
	E    Endpoint
}

type Endpoint struct {
	Name        string `json:"name"`
	Addr        string `json:"addr"`
	Port        int    `json:"port"`
	Version     string `json:"version"`
	ClusterName string `json:"cluster_name"`
	ServerName  string `json:"server_name"`
}

// 保存已经注册的服务
var serviceMap map[string]*EndpointUnit
