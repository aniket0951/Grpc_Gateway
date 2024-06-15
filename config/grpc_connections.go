package config

import (
	"context"
	"errors"
	"strings"

	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClientConnection struct {
	GrpcClientConnection *grpc.ClientConn
	ServicePath          string
}

type GrpcConnectionWrapper struct {
	*grpc.ClientConn
}

func NewGrpcConnectionWrapper(address string) (*GrpcConnectionWrapper, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return &GrpcConnectionWrapper{conn}, err
}

func (gr *GrpcConnectionWrapper) ListService() ([]string, error) {
	if gr == nil {
		return nil, errors.New("grpc connection not found")
	}
	reflectClinet := grpcreflect.NewClientAuto(context.Background(), gr)
	return reflectClinet.ListServices()
}

func (gr *GrpcConnectionWrapper) Close() error {
	return gr.ClientConn.Close()
}

func (gr *GrpcConnectionWrapper) RegisterServices(services []string) error {
	for _, service := range services {
		if service != "grpc.reflection.v1alpha.ServerReflection" && service != "grpc.reflection.v1.ServerReflection" {
			serviceData := strings.Split(service, ".")
			if len(serviceData) > 0 {
				m.Lock()
				if _, ok := GrpcConnMap[serviceData[len(serviceData)-1]]; !ok {
					GrpcConnMap[serviceData[len(serviceData)-1]] = &GrpcClientConnection{
						GrpcClientConnection: gr.ClientConn,
						ServicePath:          service,
					}
				}
				m.Unlock()
			}
		}
	}
	return nil
}
