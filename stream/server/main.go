package main

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	"github.com/gomicro_demo/stream/protobuf"
	"github.com/micro/go-micro"
	stanBroker "github.com/micro/go-plugins/broker/stan"
	natsRegistry "github.com/micro/go-plugins/registry/nats"
	natsTransport "github.com/micro/go-plugins/transport/nats"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"strings"
)

const (
	NATS_URLS  = "nats://118.31.50.67:4222"
	CLUSTER_ID = "test-cluster"
	NATS_TOKEN = "NATS12345"
	TOPIC      = "greet"
)

type Streamer struct{}

const (
	natsLimit = 1024 * 1024
)

func (r *Streamer) LogFileStream(ctx context.Context, req *protobuf.DownloadRequest, stream protobuf.LogGather_LogFileStreamStream) error {
	glog.V(4).Infof("downloading log file %s, the offset is %d and the whence is %d", req.Logfile.File, req.Offset, req.Whence)

	f, err := os.Open(req.Logfile.File)
	if err != nil {
		return fmt.Errorf("open file %s error: %v", req.Logfile.File, err)
	}
	defer f.Close()

	// calculate how many streams logger needs to send
	pos, err := f.Seek(req.Offset, int(req.Whence))
	if err != nil {
		return fmt.Errorf("seek file %s error: %v", f.Name(), err)
	}
	log.Info().Msgf("pos:%d\n", pos)
	start := 0
	total := pos / natsLimit
	if pos%natsLimit > 0 {
		total += 1
	}
	log.Info().Msgf("1---pos:%d,total:%d\n", pos, total)
	// send stream
	for i := 1; i <= int(total); i++ {
		buf := make([]byte, natsLimit)
		bytesRead, err := f.ReadAt(buf, int64(start))
		if err != nil && err != io.EOF {
			return err
		}
		log.Info().Msgf("sending data : %s", string(buf))
		if err := stream.Send(&protobuf.DownloadResponse{Total: total, Times: int64(i), DataBytes: buf,}); err != nil {
			return err
		}
		start += bytesRead
		log.Info().Msgf("2---start:%d\n", start)
	}

	return nil
}

func main() {
	service := buildservice("yzyserver")
	protobuf.RegisterLogGatherHandler(service.Server(), &Streamer{})
	if err := service.Run(); err != nil {
		log.Fatal().Err(err).Send()
	}
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
	service := micro.NewService(
		micro.Name(name),
		micro.Registry(registry),
		micro.Broker(broker),
		micro.Transport(transport),
	)
	return service
}
