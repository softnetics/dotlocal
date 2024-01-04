package daemon

import (
	"log"
	"net"

	"github.com/samber/lo"
	"github.com/softnetics/dotlocal/internal"
	api "github.com/softnetics/dotlocal/internal/api/proto"
	"github.com/softnetics/dotlocal/internal/util"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type APIServer struct {
	logger     *zap.Logger
	grpcServer *grpc.Server
	dotlocal   *DotLocal
}

func NewAPIServer(logger *zap.Logger, dotlocal *DotLocal) (*APIServer, error) {
	return &APIServer{
		logger:   logger,
		dotlocal: dotlocal,
	}, nil
}

func (s *APIServer) Start() error {
	lis, err := net.Listen("unix", util.GetApiSocketPath())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	s.grpcServer = grpc.NewServer(opts...)
	api.RegisterDotLocalServer(s.grpcServer, newDotLocalServer(s.logger, s.dotlocal))

	s.logger.Info("API server listening", zap.String("path", util.GetApiSocketPath()))
	s.grpcServer.Serve(lis)

	return nil
}

func (s *APIServer) Stop() error {
	s.grpcServer.Stop()
	return nil
}

type dotLocalServer struct {
	api.UnimplementedDotLocalServer

	logger   *zap.Logger
	dotlocal *DotLocal
}

func newDotLocalServer(logger *zap.Logger, dotlocal *DotLocal) *dotLocalServer {
	return &dotLocalServer{
		logger:   logger,
		dotlocal: dotlocal,
	}
}

func (s *dotLocalServer) CreateMappingWhileConnected(stream api.DotLocal_CreateMappingWhileConnectedServer) error {
	var createdMappings []internal.Mapping

	for {
		req, err := stream.Recv()
		if err != nil {
			if err != stream.Context().Err() {
				break
			}
			return err
		}

		s.logger.Info("CreateMappingWhileConnected", zap.Any("req", req))
		mapping, err := s.dotlocal.CreateMapping(internal.MappingOptions{
			Host:       *req.Host,
			PathPrefix: *req.PathPrefix,
			Target:     *req.Target,
		})
		if err != nil {
			return err
		}
		createdMappings = append(createdMappings, mapping)
	}

	return s.dotlocal.RemoveMapping(lo.Map(createdMappings, func(mapping internal.Mapping, _ int) string {
		return mapping.ID
	})...)
}
