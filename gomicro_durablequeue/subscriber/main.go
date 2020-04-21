package main

import (
	"flag"
	"github.com/micro/go-micro/broker"
	stanBroker "github.com/micro/go-plugins/broker/stan"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog/log"
	log2 "log"
	"strings"
	"time"
)

const (
	TOPIC = "durable.queue.test"
)

var (
	clusterID   string
	clientID    string
	URL         string
	async       bool
	token       string
	qgroup      string
	unsubscribe bool
	durable     string
)

func main() {
	flag.StringVar(&URL, "s", stan.DefaultNatsURL, "The nats server URLs (separated by comma)")
	flag.StringVar(&URL, "server", stan.DefaultNatsURL, "The nats server URLs (separated by comma)")
	flag.StringVar(&clusterID, "c", "test-cluster", "The NATS Streaming cluster ID")
	flag.StringVar(&clusterID, "cluster", "test-cluster", "The NATS Streaming cluster ID")
	flag.StringVar(&clientID, "id", "stan-pub", "The NATS Streaming client ID to connect with")
	flag.StringVar(&clientID, "clientid", "stan-pub", "The NATS Streaming client ID to connect with")
	flag.BoolVar(&async, "a", false, "Publish asynchronously")
	flag.BoolVar(&async, "async", false, "Publish asynchronously")
	flag.StringVar(&token, "cr", "", "Credentials File")
	flag.StringVar(&token, "creds", "", "Credentials File")
	flag.StringVar(&durable, "durable", "", "Durable subscriber name")
	flag.StringVar(&qgroup, "qgroup", "", "Queue group name")
	flag.BoolVar(&unsubscribe, "unsub", false, "Unsubscribe the durable on exit")
	flag.BoolVar(&unsubscribe, "unsubscribe", false, "Unsubscribe the durable on exit")

	log2.SetFlags(0)
	flag.Parse()

	s1 := buildBroker()
	if err := s1.Connect(); err != nil {
		log.Fatal().Msg("Failed to connect stan server.")
	}
	log.Info().Msg("11111")
	_, err := s1.Subscribe(TOPIC, ProcessEvent, stanBroker.SubscribeOption(
		stan.SetManualAckMode(),
		stan.AckWait(time.Second*10),
		stan.MaxInflight(1),
		stan.DurableName(durable)), broker.Queue(qgroup))
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to subscribe topic : %s", TOPIC)
	}
	select {}
}

func ProcessEvent(e broker.Event) error {
	log.Info().Msgf("processing %s times event.", string(e.Message().Body))
	return e.Ack()
}
func buildBroker() broker.Broker {
	options := nats.GetDefaultOptions()
	options.Servers = strings.Split(URL, ",")
	options.Token = token
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
	stanOptions.NatsURL = URL
	var err error
	stanOptions.NatsConn, err = nats.Connect(stanOptions.NatsURL, nats.Token(token))
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
		stanBroker.ClusterID(clusterID),
	)
	return broker
}
