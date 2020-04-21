package main

import (
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
	s3 := buildBroker()
	if err := s3.Connect(); err != nil {
		log.Fatal().Msg("Failed to connect stan server.")
	}
	_, err := s3.Subscribe(TOPIC, ProcessEvent, stanBroker.SubscribeOption(
		stan.SetManualAckMode(),
		stan.AckWait(time.Second*10),
		stan.MaxInflight(1),
		stan.DurableName("aaa"), ), broker.Queue("bbb"))
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to subscribe topic : %s", TOPIC)
	}
	select{}
}

func ProcessEvent(e broker.Event) error {
	log.Info().Msgf("processing %d times event.", string(e.Message().Body))
	return e.Ack()
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
