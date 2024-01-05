package daemon

import (
	"context"
	"errors"
	"net"
	"os"
	"strconv"

	"github.com/softnetics/dotlocal/internal"
	api "github.com/softnetics/dotlocal/internal/api/proto"
	"github.com/softnetics/dotlocal/internal/util"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type APIServer struct {
	logger     *zap.Logger
	grpcServer *grpc.Server
	dotlocal   *DotLocal
	serve      func() error
}

func NewAPIServer(logger *zap.Logger, dotlocal *DotLocal) (*APIServer, error) {
	return &APIServer{
		logger:   logger,
		dotlocal: dotlocal,
	}, nil
}

func (s *APIServer) Start() error {
	err := s.killExistingProcess()
	if err != nil {
		return err
	}

	pid := os.Getpid()
	err = os.WriteFile(util.GetPidPath(), []byte(strconv.Itoa(pid)), 0644)

	socketPath := util.GetApiSocketPath()
	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}
	var opts []grpc.ServerOption
	s.grpcServer = grpc.NewServer(opts...)
	api.RegisterDotLocalServer(s.grpcServer, newDotLocalServer(s.logger, s.dotlocal))

	s.serve = func() error {
		s.logger.Info("API server listening", zap.String("path", socketPath))
		return s.grpcServer.Serve(lis)
	}

	return nil
}

func (s *APIServer) Serve() error {
	return s.serve()
}

func (s *APIServer) Stop() error {
	s.logger.Info("Stopping API server")
	s.grpcServer.Stop()
	os.Remove(util.GetPidPath())
	return nil
}

func (s *APIServer) killExistingProcess() error {
	_, err := os.Stat(util.GetApiSocketPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	s.logger.Info("Killing existing process", zap.String("path", util.GetApiSocketPath()))

	pidBytes, err := os.ReadFile(util.GetPidPath())
	if err != nil {
		return err
	}
	pid, err := strconv.Atoi(string(pidBytes))
	if err != nil {
		return err
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	_ = process.Kill()

	err = os.Remove(util.GetPidPath())
	if err != nil {
		return err
	}
	err = os.Remove(util.GetApiSocketPath())
	if err != nil {
		return err
	}

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

func (s *dotLocalServer) CreateMapping(ctx context.Context, req *api.CreateMappingRequest) (*emptypb.Empty, error) {
	_, err := s.dotlocal.CreateMapping(internal.MappingOptions{
		Host:       *req.Host,
		PathPrefix: *req.PathPrefix,
		Target:     *req.Target,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *dotLocalServer) RemoveMapping(ctx context.Context, key *api.MappingKey) (*emptypb.Empty, error) {
	err := s.dotlocal.RemoveMapping(internal.MappingKey{
		Host:       *key.Host,
		PathPrefix: *key.PathPrefix,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
