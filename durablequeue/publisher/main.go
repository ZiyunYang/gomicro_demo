package main

import (
	"fmt"
	"github.com/micro/go-micro/broker"
	stanBroker "github.com/micro/go-plugins/broker/stan"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

const (
	NATS_URLS  = "nats://118.31.50.67:4222"
	CLUSTER_ID = "test-cluster"
	NATS_TOKEN = "NATS12345"
	TOPIC      = "durable.queue.test"
)

func main() {
	publisher := buildBroker()
	if err := publisher.Connect(); err != nil {
		log.Error().Msg("Faild to connct stan server.")
		panic(err)
	}
	for i := 0; i < 100; i++ {
		err:=publisher.Publish(TOPIC, &broker.Message{Body: []byte(fmt.Sprintf("%d---", i))})
		if err!=nil{
			log.Err(err).Msgf("got publish err")
		}
		log.Info().Msgf("Publishing %d",i)
		time.Sleep(time.Second * 5)
	}

}

func buildBroker() broker.Broker {
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
	//registry := natsRegistry.NewRegistry(natsRegistry.Options(options))
	//transport := natsTransport.NewTransport(natsTransport.Options(options))

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
	)
	return broker
}
