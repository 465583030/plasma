package server

import (
	"io"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	"github.com/openfresh/plasma/config"
	"github.com/openfresh/plasma/event"
	"github.com/openfresh/plasma/log"
	"github.com/openfresh/plasma/manager"
	"github.com/openfresh/plasma/metrics"
	"github.com/openfresh/plasma/protobuf"
	"github.com/openfresh/plasma/pubsub"
	"github.com/pkg/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
)

type GRPCServer struct {
	*grpc.Server
	accessLogger *zap.Logger
	errorLogger  *zap.Logger
	config       config.Config
}

func (s *GRPCServer) StreamAccessLogHandler(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	start := time.Now()
	err := handler(srv, ss)
	fields := log.GRPCRequestToLogFields(info, start, err)

	s.accessLogger.Info("grpc", fields...)

	return err
}

func NewGRPCServer(opt Option) (*GRPCServer, error) {
	gs := &GRPCServer{
		accessLogger: opt.AccessLogger,
		errorLogger:  opt.ErrorLogger,
		config:       opt.Config,
	}

	opts := make([]grpc.ServerOption, 0)

	tls := opt.Config.TLS
	if tls.CertFile != "" && tls.KeyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(tls.CertFile, tls.KeyFile)
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.Creds(creds))
	}

	opts = append(opts, grpc.StreamInterceptor(gs.StreamAccessLogHandler))
	gs.Server = grpc.NewServer(opts...)

	proto.RegisterStreamServiceServer(gs.Server, NewStreamServer(opt))

	return gs, nil
}

type refreshEvents struct {
	client *manager.Client
	events []string
}

type StreamServer struct {
	clientManager  *manager.ClientManager
	newClients     chan manager.Client
	removeClients  chan manager.Client
	payloads       chan event.Payload
	resfreshEvents chan refreshEvents
	errChan        chan error
	forceCloseChan chan manager.Client
	pubsub         pubsub.PubSuber
	accessLogger   *zap.Logger
	errorLogger    *zap.Logger
}

func NewStreamServer(opt Option) *StreamServer {
	ss := &StreamServer{
		clientManager:  manager.NewClientManager(),
		newClients:     make(chan manager.Client),
		removeClients:  make(chan manager.Client),
		payloads:       make(chan event.Payload),
		errChan:        make(chan error),
		forceCloseChan: make(chan manager.Client),
		resfreshEvents: make(chan refreshEvents),
		pubsub:         opt.PubSuber,
		accessLogger:   opt.AccessLogger,
		errorLogger:    opt.ErrorLogger,
	}
	ss.pubsub.Subscribe(func(payload event.Payload) {
		ss.payloads <- payload
	})
	ss.Run()

	return ss
}

func (ss *StreamServer) Run() {
	go func() {
		for {
			select {
			case client := <-ss.newClients:
				ss.clientManager.AddClient(client)
				metrics.IncConnection()
			case client := <-ss.removeClients:
				ss.clientManager.RemoveClient(client)
				metrics.DecConnection()
			case payload := <-ss.payloads:
				ss.clientManager.SendPayload(payload)
			case re := <-ss.resfreshEvents:
				ss.clientManager.DeleteEvents(re.client)
				re.client.SetEvents(re.events)
				ss.clientManager.AddClient(*re.client)
			}
		}
	}()
}

func (ss *StreamServer) Events(es proto.StreamService_EventsServer) error {
	client := manager.NewClient([]string{})
	ss.newClients <- client
	go func() {
		for {
			request, err := es.Recv()
			if err == io.EOF {
				<-es.Context().Done()
				return
			}

			if err != nil {
				if grpc.Code(err) != codes.Canceled {
					ss.errChan <- errors.Wrap(err, "Recv error")
					return
				} else {
					<-es.Context().Done()
					return
				}
			}

			if request.ForceClose {
				ss.forceCloseChan <- client
				return
			}

			ss.accessLogger.Info("gRPC",
				zap.Array("request-events", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
					for _, e := range request.Events {
						enc.AppendString(e.Type)
					}
					return nil
				})),
				zap.String("time", time.Now().Format(time.RFC3339)),
			)
			if request.Events == nil {
				ss.errChan <- errors.New("event can't be nil")
				return
			}

			l := len(request.Events)
			events := make([]string, l)
			for i := 0; i < l; i++ {
				events[i] = request.Events[i].Type
			}
			ss.resfreshEvents <- refreshEvents{
				client: &client,
				events: events,
			}
		}
	}()

	for {
		select {
		case err := <-ss.errChan:
			return err
		case pl, open := <-client.ReceivePayload():
			if !open {
				return nil
			}
			eventType := proto.EventType{Type: pl.Meta.Type}
			p := &proto.Payload{
				EventType: &eventType,
				Data:      string(pl.Data),
			}
			if err := es.Send(p); err != nil {
				ss.errorLogger.Error("failed to send message",
					zap.Error(err),
					zap.Object("payload", pl),
				)
				ss.removeClients <- client
				return err
			}
		case <-ss.forceCloseChan:
			ss.removeClients <- client
			return nil
		case <-es.Context().Done():
			ss.removeClients <- client
			return nil
		}

	}
}
