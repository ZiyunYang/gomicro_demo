package main

import (
	"flag"
	"fmt"
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
	clusterID string
	clientID  string
	URL       string
	async     bool
	token     string
)

func main() {
	publisher := buildBroker()
	if err := publisher.Connect(); err != nil {
		log.Error().Msg("Faild to connct stan server.")
		panic(err)
	}
	for i := 0; i < 100; i++ {
		err := publisher.Publish(TOPIC, &broker.Message{Body: []byte(fmt.Sprintf("%d---", i))})
		if err != nil {
			log.Err(err).Msgf("got publish err")
		}
		log.Info().Msgf("Publishing %d", i)
		time.Sleep(time.Second * 5)
	}

}

func buildBroker() broker.Broker {
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

	log2.SetFlags(0)
	//flag.Usage = usage
	flag.Parse()
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
		stanBroker.ClientID(clientID),
	)
	return broker
}
