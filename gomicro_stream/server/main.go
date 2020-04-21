package main

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	"github.com/gomicro_demo/gomicro_stream/protobuf"
	"github.com/micro/go-micro"
	stanBroker "github.com/micro/go-plugins/broker/stan"
	natsRegistry "github.com/micro/go-plugins/registry/nats"
	natsTransport "github.com/micro/go-plugins/transport/nats"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog/log"
	"io/ioutil"
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
	natsLimit = 32 * 1024
)

func init() {
	var str string
	for i := 0; i < 1000; i++ {
		str = fmt.Sprintf("%s\n%s--%d", str, "log", i)
	}
	err := ioutil.WriteFile("/Users/xkahj/logtest.txt", []byte(str), 0777)
	if err != nil {
		panic(err)
	}
}

func (r *Streamer) LogFileStream(ctx context.Context, req *protobuf.DownloadRequest, stream protobuf.LogGather_LogFileStreamStream) error {
	glog.V(4).Infof("downloading log file %s, the offset is %d and the whence is %d", req.Logfile.File, req.Offset, req.Whence)

	f, err := os.Open(req.Logfile.File)
	if err != nil {
		return fmt.Errorf("open file %s error: %v", req.Logfile.File, err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("open file %s error: %v", req.Logfile.File, err)
	}

	//
	//// calculate how many streams logger needs to send
	////pos, err := f.Seek(req.Offset, int(req.Whence))
	//if err != nil {
	//	return fmt.Errorf("seek file %s error: %v", f.Name(), err)
	//}
	//log.Info().Msgf("pos:%d\n", pos)
	//start := 0
	total := (fi.Size() + natsLimit - 1) / natsLimit
	//log.Info().Msgf("1---pos:%d,total:%d\n", pos, total)
	// send stream
	buf := make([]byte, natsLimit)
	var n int
	for i := 1; i <= int(total); i++ {
		n, err = f.Read(buf)
		if err != nil {
			log.Error().Msgf("err--", err.Error())
			return err
		}

		//bytesRead, err := f.ReadAt(buf, int64(start))

		log.Info().Msgf("sending data : %d", n)
		if err := stream.Send(&protobuf.DownloadResponse{Total: total, Times: int64(i), DataBytes: buf[:n]}); err != nil {
			return err
		}
		//start += bytesRead
		//log.Info().Msgf("2---start:%d\n", start)
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
