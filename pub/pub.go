package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	stanBroker "github.com/micro/go-plugins/broker/stan"
	natsRegistry "github.com/micro/go-plugins/registry/nats"
	natsTransport "github.com/micro/go-plugins/transport/nats"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog/log"
	"gomicrostan/greeter"
	"strings"
	"time"
)

const (
	NATS_URLS  = "nats://118.31.50.67:4222"
	CLUSTER_ID = "test-cluster"
	NATS_TOKEN = "NATS12345"
	TOPIC      = "greet"
)

var (
	publisher micro.Publisher
)

//type Greeter struct{}
//
//func (g *Greeter) Hello(ctx context.Context, req *greeter.Request, rsp *greeter.Response) error {
//	rsp.Msg = "Hello 3"
//	log.Info().Msg(rsp.Msg)
//	return nil
//}
func Server() {
	options := nats.GetDefaultOptions()
	options.Servers = strings.Split(NATS_URLS, ",")
	options.Token = NATS_TOKEN
	options.ReconnectedCB = func(*nats.Conn) {
		log.Info().Msg("NATS reconnected!")
	}
	options.ClosedCB = func(conn *nats.Conn) {
		if conn.LastError() != nil {
			//cancel()
			log.Error().Msgf("NATS connection closed! Reason: %v", conn.LastError())
		}
		log.Info().Msg("NATS connection closed!")
	}
	options.DisconnectedErrCB = func(c *nats.Conn, err error) {
		log.Info().Msgf("NATS disconnected.Err:%s", err.Error())
	}
	registry := natsRegistry.NewRegistry(natsRegistry.Options(options))
	transport := natsTransport.NewTransport(natsTransport.Options(options))

	stanOptions := stan.GetDefaultOptions()
	stanOptions.NatsURL = NATS_URLS
	var err error
	stanOptions.NatsConn, err = nats.Connect(stanOptions.NatsURL, nats.Token(NATS_TOKEN))
	if err != nil {
		log.Error().Err(err).Msgf("nats.Connect err: %v", err)
	}
	stanOptions.ConnectionLostCB = func(conn stan.Conn, e error) {
		defer conn.Close()
		if e != nil {
			fmt.Printf("go-stan close! Reason: %v", e)
		}
		log.Info().Msg("go-stan close!")
	}
	broker := stanBroker.NewBroker(
		stanBroker.Options(stanOptions),
		stanBroker.ClusterID(CLUSTER_ID),
		stanBroker.ClientID("client-123"),
		stanBroker.DurableName("sayhi"),
	)
	client := micro.NewService(
		micro.Name("yzyserpub"),
		micro.Registry(registry),
		micro.Broker(broker),
		micro.Transport(transport),
	)
	//err := greeter.RegisterGreeterHandler(server.Server(), &Greeter{})
	//if err != nil {
	//	panic(err)
	//}

	publisher = micro.NewPublisher(TOPIC, client.Client())

	err = client.Run()
	if err != nil {
		//panic(err)
	}
}

func main() {
	go Server()
	time.Sleep(time.Second * 5)
	i := 1
	for {
		log.Info().Msgf("%d times greeting.", i)
		publisher.Publish(context.Background(), &greeter.Request{Name:fmt.Sprintf("%d times greeting.", i)})
		i++
		time.Sleep(time.Second * 5)
	}
}
