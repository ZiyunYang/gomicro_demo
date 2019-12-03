package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	stanBroker "github.com/micro/go-plugins/broker/stan"
	natsRegistry "github.com/micro/go-plugins/registry/nats"
	natsTransport "github.com/micro/go-plugins/transport/nats"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog/log"
	micro2 "github.com/upengs/gomicro_demo/pkg/micro"
	"gomicrostan/greeter"
	"strings"
)

const (
	NATS_URLS  = "nats://118.31.50.67:4222"
	CLUSTER_ID = "test-cluster"
	NATS_TOKEN = "NATS12345"
	TOPIC      = "greet"
)

func main() {
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
			log.Error().Msgf("go-stan close! Reason: %v", e)
		}
		log.Info().Msg("go-stan close!")
	}
	broker := stanBroker.NewBroker(
		stanBroker.Options(stanOptions),
		stanBroker.ClusterID(CLUSTER_ID),
		stanBroker.DurableName(TOPIC),
	)
	server := micro.NewService(
		micro.Name("yzysub"),
		micro.Registry(registry),
		micro.Broker(broker),
		micro.Transport(transport),
	)
	micro.RegisterSubscriber(TOPIC, server.Server(), Listen,
		stanBroker.ServerSubscriberOption(stan.DeliverAllAvailable()),
	stanBroker.ServerSubscriberOption(stan.SetManualAckMode()),
	stanBroker.ServerSubscriberOption(stan.MaxInflight(1)))
	err = server.Run()
	if err != nil {
		log.Error().Err(err).Msg("Failed to register subscriber.")
	}
}


func Listen(ctx context.Context, request *greeter.Request) error {
	log.Info().Msg(request.Name)
	fmt.Println(request.Name)
	return nil
}

func sub()  {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := micro2.NewGoMicroController(ctx, []string{NATS_URLS}, "", NATS_TOKEN, CLUSTER_ID).Subscribes([]micro2.Subscribes{{Topic:TOPIC,Handle: func(event broker.Event) error {
		fmt.Println(string(event.Message().Body))
		event.Ack()
		return nil
	},SubscribeOption:stanBroker.SubscribeOption(stan.DeliverAllAvailable())}}...)
	if err!=nil{
		panic(err)
	}
}