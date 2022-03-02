package grpcclient

import (
	"two_services/assignee-services/config"
)

type GrpcClientI interface {
}

type GrpcClient struct {
	cfg         config.Config
	connections map[string]interface{}
}

func New(cfg config.Config) (*GrpcClient, error) {
	

	return &GrpcClient{
		cfg:         cfg,
		connections: map[string]interface{}{},
	}, nil
}
