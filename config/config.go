package config

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	GrpcConnMap map[string]*GrpcClientConnection
	m           sync.Mutex
	natsWrapper *NatsConnWrapper
	err         error
)

func Init() {

	natsWrapper, err = NewNatsConnWrapper(GetAppConfig().Natsurl)
	if err != nil {
		logrus.Error(err)
	}
	GrpcConnMap = make(map[string]*GrpcClientConnection)

	if err := natsWrapper.RegisterServices(); err != nil {
		logrus.Error(err)
	}

}

func CloseAllConnections() {
	natsWrapper.CloseConnections()
}
