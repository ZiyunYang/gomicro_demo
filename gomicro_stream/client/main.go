package main

import (
	"context"
	"github.com/gomicro_demo/gomicro_stream/protobuf"
	"github.com/micro/go-micro"
	stanBroker "github.com/micro/go-plugins/broker/stan"
	natsRegistry "github.com/micro/go-plugins/registry/nats"
	natsTransport "github.com/micro/go-plugins/transport/nats"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"strings"
)

const (
	NATS_URLS  = "nats://118.31.50.67:4222"
	CLUSTER_ID = "test-cluster"
	NATS_TOKEN = "NATS12345"
	TOPIC      = "greet"
)

func main() {
	service := buildservice("yzyclient")
	client := protobuf.NewLogGatherService("yzyserver", service.Client())

	n := "tcp"
	addr := "127.0.0.1:9094"
	l, err := net.Listen(n, addr)
	if err != nil {
		panic("AAAAH")
	}

	/* HTTP server */
	server := http.Server{
		Handler: http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Disposition", "attachment; filename=log.txt")
			res.Header().Set("Content-Type", http.DetectContentType(nil))
			err := DownloadLogFile(res, client, "/Users/xkahj/logtest.txt", int64(1000), int32(2))
			if err != nil {
				res.WriteHeader(400)
				res.Write([]byte("Failed"))
			}
		}),
	}
	server.Serve(l)
}

func DownloadLogFile(res http.ResponseWriter, client protobuf.LogGatherService, path string, offset int64, whence int32) error {
	downloadReq := &protobuf.DownloadRequest{
		Logfile: &protobuf.Logfile{
			File:  path,
			Topic: "",
		},
		Offset: offset,
		Whence: int64(whence),
	}
	stream, err := client.LogFileStream(context.TODO(), downloadReq)
	rsp, err := stream.Recv()
	if err != nil {
		log.Error().Msgf("Failed to get log file from wise-logger. %s", err.Error())
		return err
	}
	//res.Header().Set("Content-Disposition", "attachment; filename=log.txt")
	//res.Header().Set("Content-Type", http.DetectContentType(nil))
	total := int(rsp.Total)
	last := 1
	res.Write(rsp.DataBytes)
	for i := 0; i < total; i++ {
		log.Info().Msgf("receiving %d ", i)
		rsp, err := stream.Recv()
		if err != nil {
			log.Error().Msg("Failed to get log file from wise-logger.")
			break
		}
		if last+1 != int(rsp.Times) {
			break
		}
		res.Write(rsp.DataBytes)
		last++
	}

	if last < total {
		log.Error().Msgf("Lost some data, received %d sections data", last)
		res.Write([]byte("The downloaded log file is incomplete, please retry it."))
	}
	if err := stream.Close(); err != nil {
		log.Error().Msgf("Failed to close stream. %s", err.Error())
		return err
	}
	log.Info().Msg("finish receiving")
	return nil
}

func buildservice(name string) micro.Service {
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
	client := micro.NewService(
		micro.Name(name),
		micro.Registry(registry),
		micro.Broker(broker),
		micro.Transport(transport),
	)
	return client
}
