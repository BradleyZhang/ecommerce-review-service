package registory

import (
	"context"
	"fmt"

	"github.com/BradleyZhang/ecommerce-review-service/common"
	"github.com/hashicorp/consul/api"
)

type Registry struct {
	consul *api.Client
	config Config
}
type Config struct {
	ConsulAddress string
	ServiceName   string
	ServiceID     string
	ServicePort   int
	ServiceIP     string
	// HealthPath    string
	Tags []string
}

func NewClient(cfg Config) (*Registry, error) {
	consulConfig := api.DefaultConfig()
	if cfg.ConsulAddress != "" {
		consulConfig.Address = cfg.ConsulAddress
	}
	consul, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}
	return &Registry{
		consul: consul,
		config: cfg,
	}, nil
}
func DefaultConfig() Config {
	return Config{
		ConsulAddress: common.ConsulAddress,
		ServiceName:   common.Name,
		ServiceIP:     common.IP,
		ServicePort:   common.PortGRPC,
		ServiceID:     fmt.Sprintf("%s-%s-%d", common.Name, common.IP, common.PortGRPC),
		// HealthPath:    "/health",
	}
}

func (c *Registry) RegisterWithContext(ctx context.Context) error {
	err := c.Register()
	if err != nil {
		return err
	}
	defer func() {
		go func() {
			<-ctx.Done()
			if err := c.Deregister(); err != nil {
				common.ErrorLog(err)
			}
		}()
	}()
	return nil
}
func (c *Registry) Register() error {
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", c.config.ServiceIP, c.config.ServicePort),
		Interval:                       "10s",
		Timeout:                        "5s",
		DeregisterCriticalServiceAfter: "1m",
	}

	registration := &api.AgentServiceRegistration{
		ID:      c.config.ServiceID,
		Name:    c.config.ServiceName,
		Address: c.config.ServiceIP,
		Port:    c.config.ServicePort,
		Tags:    c.config.Tags,
		Check:   check,
	}

	err := c.consul.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}
	common.SysLogf("Registered serviced %s with ID %s", c.config.ServiceName, c.config.ServiceID)
	return nil
}

func (c *Registry) Deregister() error {
	err := c.consul.Agent().ServiceDeregister(c.config.ServiceID)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}
	common.SysLogf("Deregistered service with ID %s", c.config.ServiceID)
	return nil
}
