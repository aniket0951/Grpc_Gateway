package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type SubscribeHandler func(msg *nats.Msg)

type NatsConnWrapper struct {
	Conn             *nats.Conn
	service, address string
}

func NewNatsConnWrapper(url string) (*NatsConnWrapper, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, errors.New("nats connection error : " + err.Error())
	}

	return &NatsConnWrapper{
		Conn: conn,
	}, nil
}

func (n *NatsConnWrapper) SubscribeToTopic(subj string, handler SubscribeHandler) (*nats.Subscription, error) {

	return n.Conn.Subscribe(subj, func(msg *nats.Msg) {
		logrus.Infof("Received message on subject %s: %s", subj, string(msg.Data))
		handler(msg)
	})
}

func (n *NatsConnWrapper) RegisterServices() error {
	// logrus.Info("")
	subscription, err := n.SubscribeToTopic(GetAppConfig().GateWayTopic, n.handleRegisterService)
	logrus.Infof("Subscribed to subject: %s", subscription.Subject)
	return err
}

func (n *NatsConnWrapper) handleRegisterService(msg *nats.Msg) {
	reqMsg := string(msg.Data)
	logrus.Info("Request Message : ", reqMsg)
	requestMsg := strings.Split(string(msg.Data), ":")

	if len(requestMsg) <= 1 {
		n.PublishMsg([]byte("Invalid Request"), msg.Reply)
		return
	}

	service := requestMsg[0]
	address := requestMsg[1] + ":" + requestMsg[len(requestMsg)-1]
	n.address = address
	n.service = service

	// register the service

	err := n.HandleServiceRegister()
	if err != nil {
		n.PublishMsg([]byte(err.Error()), msg.Reply)
		return
	}
	n.PublishMsg([]byte(fmt.Sprintf("Service %s with topic %s has been register", n.service, n.address)), msg.Reply)

}

func (n *NatsConnWrapper) PublishMsg(msg []byte, subject string) error {
	if n.Conn == nil {
		return errors.New("failed to published msg")
	}
	return n.Conn.Publish(subject, msg)
}

func (n *NatsConnWrapper) CloseConnections() {
	n.Conn.Close()
}

// handle the grpc connection instances
func (n *NatsConnWrapper) HandleServiceRegister() error {
	// if _, ok := GrpcConnMap[n.service]; ok {
	// 	logrus.Info("Service Already Register")
	// 	return nil
	// }

	grpcConn, err := NewGrpcConnectionWrapper(n.address)

	if err != nil {
		return nil
	}

	services, err := grpcConn.ListService()

	if err != nil {
		return err
	}

	return grpcConn.RegisterServices(services)
}
