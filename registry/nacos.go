package registry

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type Nacos struct {
	client        naming_client.INamingClient
	configClient  config_client.IConfigClient
	Namespace     string
	Address       string
	NacosUser     string
	NacosPassword string
}

func (n *Nacos) NewNamingService() error {
	addresses := strings.Split(n.Address, ";")
	serverConfigs := []constant.ServerConfig{}
	for _, a := range addresses {
		address := strings.Split(a, ":")
		port, err := strconv.ParseUint(address[1], 10, 64)
		if err != nil {
			return err
		}
		serverConfigs = append(serverConfigs, constant.ServerConfig{
			IpAddr: address[0],
			Port:   port,
		})
	}
	clientConfig := constant.NewClientConfig(
		constant.WithNamespaceId(n.Namespace),
		constant.WithUsername(n.NacosUser),
		constant.WithPassword(n.NacosPassword),
		// constant.WithTimeoutMs(5000),
		constant.WithTimeoutMs(30*1000),
		constant.WithBeatInterval(10*1000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("error"),
		// constant.WithAppName("yiyantest"),
		// constant.WithEndpoint("jmenv.tbsite.net:8080"),
		// constant.WithClusterName("serverlist"),
		// constant.WithEndpointQueryParams("nofix=1"),
	)

	configClient, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: serverConfigs,
	})
	if err != nil {
		return err
	}

	client, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: serverConfigs,
	})
	if err != nil {
		return err
	}
	n.client = client
	n.configClient = configClient
	serviceMap = make(map[string]*EndpointUnit)
	return nil
}

// 添加/注册新的服务
func (n *Nacos) AddEndpoint(e Endpoint) error {
	instance := vo.RegisterInstanceParam{
		Ip:          e.Addr,
		Port:        uint64(e.Port),
		Enable:      true,
		Healthy:     true,
		ClusterName: e.ClusterName,
		ServiceName: e.Name,
		GroupName:   e.ServerName,
	}
	_, err := n.client.RegisterInstance(instance)
	if err != nil {
		return err
	}
	fmt.Printf("success registry cluster:%s group: %s service:%s address: %s:%d !!! \n", e.ClusterName, e.ServerName, e.Name, e.Addr, e.Port)
	serviceMap[e.Name] = &EndpointUnit{
		Name: e.Name,
		E:    e,
	}
	return nil
}

// 移除一个服务
func (n *Nacos) DelEndpoint(e Endpoint) error {
	instance := vo.DeregisterInstanceParam{
		Ip:          e.Addr,
		Port:        uint64(e.Port),
		Cluster:     e.ClusterName,
		GroupName:   e.ServerName,
		ServiceName: e.Name,
	}
	_, err := n.client.DeregisterInstance(instance)
	if err != nil {
		return err
	}
	return nil
}

// 移除一个服务
func (n *Nacos) DelAllEndpoint() error {
	for _, v := range serviceMap {
		err := n.DelEndpoint(v.E)
		if err != nil {
			log.Fatalln("Ignore Failure Continue...")
		}
	}
	return nil
}

// 选择一个服务
func (n *Nacos) SelectEndpoint(regionName, serverName, serviceName string) ([]model.Instance, error) {
	endpoints, err := n.client.SelectAllInstances(vo.SelectAllInstancesParam{
		Clusters:    []string{regionName},
		ServiceName: serviceName,
		GroupName:   serverName,
	})
	if err != nil {
		return nil, err
	}
	return endpoints, nil
}

func (n *Nacos) ConfigPush(dataId, group string) error {
	b, err := n.configClient.PublishConfig(vo.ConfigParam{
		DataId:  dataId,
		Group:   group,
		Content: "{}",
	})
	if err != nil {
		return err
	}
	if !b {
		return fmt.Errorf("error")
	}
	return nil
}

func (n *Nacos) ListenConfig(dataId, group string) error {
	return n.configClient.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
		},
	})
}

func (n *Nacos) GetConfig(dataId, group string) (string, error) {
	return n.configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
}

func (n *Nacos) ConfigClient() config_client.IConfigClient {
	return n.configClient
}
