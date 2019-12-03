package micro

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/server"
	stanBroker "github.com/micro/go-plugins/broker/stan"
	natsRegistry "github.com/micro/go-plugins/registry/nats"
	natsTransport "github.com/micro/go-plugins/transport/nats"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog/log"
	"strings"
)

type Subscribes struct {
	Topic           string
	Handle          broker.Handler
	SubscribeOption broker.SubscribeOption
}

type GoMicroController struct {
	o       nats.Options
	broker  broker.Broker
	service micro.Service
	name    string
	ctx     context.Context
}

func NewGoMicroController(ctx context.Context, NatsUrl []string, Name string, Token, ClusterId string) *GoMicroController {
	g := &GoMicroController{ctx: ctx, o: nats.GetDefaultOptions()}
	g.buildGoMicro(ctx, Token, NatsUrl, ClusterId)
	g.service = micro.NewService(
		micro.Name(Name),
		micro.Registry(natsRegistry.NewRegistry(natsRegistry.Options(g.o))),
		micro.Broker(g.broker),
		micro.Transport(natsTransport.NewTransport(natsTransport.Options(g.o))),
		micro.Context(g.ctx),
	)
	return g
}

func (g *GoMicroController) buildGoMicro(ctx context.Context, Token string, NatsUrl []string, ClusterId string) {

	g.o.ClosedCB = func(conn *nats.Conn) {
		fmt.Println("closed")
	}
	g.o.ReconnectedCB = func(conn *nats.Conn) {
		fmt.Println("reconnected")
	}
	g.o.DisconnectedCB = func(conn *nats.Conn) {
		fmt.Println("disconnected_cb")
	}
	g.o.DiscoveredServersCB = func(conn *nats.Conn) {
		fmt.Println("discoveredServers_cb")
	}
	g.o.Token = Token
	g.o.Servers = NatsUrl
	stanOptions := stan.GetDefaultOptions()
	stanOptions.NatsURL = strings.Join(NatsUrl, ",")
	var err error
	conn, err := nats.Connect(stanOptions.NatsURL, nats.Token(Token))
	stanOptions.NatsConn = conn
	if err != nil {
		log.Error().Err(err).Msg("nats.Connect")
		panic(err)
	}
	stanOptions.ConnectionLostCB = func(conn stan.Conn, e error) {
		if e != nil {
			log.Error().Err(err).Msg("go-stan close!")
		}
	}
	g.broker = stanBroker.NewBroker(
		stanBroker.Options(stanOptions),
		stanBroker.ClusterID(ClusterId),
	)
}

func (g *GoMicroController) Service() micro.Service {
	return g.service
}

func (g *GoMicroController) Client() client.Client {
	return g.service.Client()
}

// 注册服务到nats-stream
func (g *GoMicroController) RegisterHandler(register func(s server.Server) error) error {
	if register == nil {
		panic("register handler is nil")
	}
	return register(g.service.Server())
}

func (g *GoMicroController) Subscribes(subs ...Subscribes) error {
	if err := g.broker.Connect(); err != nil {
		return err
	}

	for _, sub := range subs {
		if _, err := g.broker.Subscribe(sub.Topic, sub.Handle, sub.SubscribeOption); err != nil {
			return err
		}
	}
	return nil
}

func (g *GoMicroController) Publish(topic string, m interface{}) error {
	if len(topic) == 0 {
		return errors.New("publish topic is empty")
	}

	if err := g.broker.Connect(); err != nil {
		return err
	}
	msg := &broker.Message{}
	switch m.(type) {
	case string:
		msg.Body = []byte(m.(string))
	case []byte:
		msg.Body = m.([]byte)
	default:
		if b, err := json.Marshal(m); err != nil {
			return err
		} else {
			msg.Body = b
		}
	}
	err := g.broker.Publish(topic, msg)
	return err
}
