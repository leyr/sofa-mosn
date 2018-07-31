package util

import (
	"github.com/alipay/sofa-mosn/internal/api/v2"
	"github.com/alipay/sofa-mosn/pkg/config"
	"github.com/alipay/sofa-mosn/pkg/types"
)

// can set
var (
	MeshLogPath  = "stdout"
	MeshLogLevel = "WARN"
)

// Create Mesh Config
func newFilterChain(name string, downstream, upstream types.Protocol, routers []v2.Router) config.FilterChain {
	proxy := &v2.Proxy{
		DownstreamProtocol: string(downstream),
		UpstreamProtocol:   string(upstream),
		VirtualHosts: []*v2.VirtualHost{
			&v2.VirtualHost{
				Name:    name,
				Domains: []string{"*"},
				Routers: routers,
			},
		},
	}
	chains := make(map[string]interface{})
	b, _ := json.Marshal(proxy)
	json.Unmarshal(b, &chains)
	return config.FilterChain{
		Filters: []config.FilterConfig{
			config.FilterConfig{Type: "proxy", Config: chains},
		},
	}
}
func newTLSFilterChain(name string, downstream, upstream types.Protocol, routers []v2.Router, tls config.TLSConfig) config.FilterChain {
	chain := newFilterChain(name, downstream, upstream, routers)
	chain.TLS = tls
	return chain
}

func newBasicCluster(name string, hosts []string) config.ClusterConfig {
	var vhosts []config.HostConfig
	for _, addr := range hosts {
		vhosts = append(vhosts, config.HostConfig{
			Address: addr,
		})
	}
	return config.ClusterConfig{
		Name:                 name,
		Type:                 "SIMPLE",
		LbType:               "LB_ROUNDROBIN",
		MaxRequestPerConn:    1024,
		ConnBufferLimitBytes: 16 * 1026,
		Hosts:                vhosts,
	}
}

func newBasicTLSCluster(name string, hosts []string, tls config.TLSConfig) config.ClusterConfig {
	cfg := newBasicCluster(name, hosts)
	cfg.TLS = tls
	return cfg
}

func newListener(name, addr string, chains []config.FilterChain) config.ListenerConfig {
	return config.ListenerConfig{
		Name:         name,
		Address:      addr,
		BindToPort:   true,
		LogPath:      MeshLogPath,
		LogLevel:     MeshLogLevel,
		FilterChains: chains,
	}
}

func newMOSNConfig(listeners []config.ListenerConfig, clusterManager config.ClusterManagerConfig) *config.MOSNConfig {
	return &config.MOSNConfig{
		Servers: []config.ServerConfig{
			config.ServerConfig{
				DefaultLogPath:  MeshLogPath,
				DefaultLogLevel: MeshLogLevel,
				Listeners:       listeners,
			},
		},
		ClusterManager: clusterManager,
	}
}

// common case
func newHeaderRouter(cluster string, value string) v2.Router {
	header := v2.HeaderMatcher{Name: "service", Value: value}
	return v2.Router{
		Match: v2.RouterMatch{Headers: []v2.HeaderMatcher{header}},
		Route: v2.RouteAction{ClusterName: cluster},
	}
}
func newPrefixRouter(cluster string, prefix string) v2.Router {
	return v2.Router{
		Match: v2.RouterMatch{Prefix: prefix},
		Route: v2.RouteAction{ClusterName: cluster},
	}
}
func newPathRouter(cluster string, path string) v2.Router {
	return v2.Router{
		Match: v2.RouterMatch{Path: path},
		Route: v2.RouteAction{ClusterName: cluster},
	}
}

//tls config

const (
	cacert string = `-----BEGIN CERTIFICATE-----
MIIBVzCB/qADAgECAhBsIQij0idqnmDVIxbNRxRCMAoGCCqGSM49BAMCMBIxEDAO
BgNVBAoTB0FjbWUgQ28wIBcNNzAwMTAxMDAwMDAwWhgPMjA4NDAxMjkxNjAwMDBa
MBIxEDAOBgNVBAoTB0FjbWUgQ28wWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAARV
DG+YT6LzaR5r0Howj4/XxHtr3tJ+llqg9WtTJn0qMy3OEUZRfHdP9iuJ7Ot6rwGF
i6RXy1PlAurzeFzDqQY8ozQwMjAOBgNVHQ8BAf8EBAMCAqQwDwYDVR0TAQH/BAUw
AwEB/zAPBgNVHREECDAGhwR/AAABMAoGCCqGSM49BAMCA0gAMEUCIQDt9WA96LJq
VvKjvGhhTYI9KtbC0X+EIFGba9lsc6+ubAIgTf7UIuFHwSsxIVL9jI5RkNgvCA92
FoByjq7LS7hLSD8=
-----END CERTIFICATE-----
`
	certchain string = `-----BEGIN CERTIFICATE-----
MIIBdDCCARqgAwIBAgIQbCEIo9Inap5g1SMWzUcUQjAKBggqhkjOPQQDAjASMRAw
DgYDVQQKEwdBY21lIENvMCAXDTcwMDEwMTAwMDAwMFoYDzIwODQwMTI5MTYwMDAw
WjASMRAwDgYDVQQKEwdBY21lIENvMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE
VQxvmE+i82kea9B6MI+P18R7a97SfpZaoPVrUyZ9KjMtzhFGUXx3T/Yriezreq8B
hYukV8tT5QLq83hcw6kGPKNQME4wDgYDVR0PAQH/BAQDAgWgMB0GA1UdJQQWMBQG
CCsGAQUFBwMBBggrBgEFBQcDAjAMBgNVHRMBAf8EAjAAMA8GA1UdEQQIMAaHBH8A
AAEwCgYIKoZIzj0EAwIDSAAwRQIgO9xLIF1AoBsSMU6UgNE7svbelMAdUQgEVIhq
K3gwoeICIQCDC75Fa3XQL+4PZatS3OfG93XNFyno9koyn5mxLlDAAg==
-----END CERTIFICATE-----
`
	privatekey string = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEICWksdaVL6sOu33VeohiDuQ3gP8xlQghdc+2FsWPSkrooAoGCCqGSM49
AwEHoUQDQgAEVQxvmE+i82kea9B6MI+P18R7a97SfpZaoPVrUyZ9KjMtzhFGUXx3
T/Yriezreq8BhYukV8tT5QLq83hcw6kGPA==
-----END EC PRIVATE KEY-----
`
)
