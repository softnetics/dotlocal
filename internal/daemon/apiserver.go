package daemon

import (
	"context"
	"errors"
	"net"
	"os"
	"strconv"

	"github.com/samber/lo"
	"github.com/softnetics/dotlocal/internal"
	api "github.com/softnetics/dotlocal/internal/api/proto"
	"github.com/softnetics/dotlocal/internal/util"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (s *APIServer) Start(ctx context.Context) error {
	err := s.killExistingProcessIfNeeded()
	if err != nil {
		return err
	}

	pid := os.Getpid()
	err = os.WriteFile(util.GetPidPath(), []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		return err
	}

	socketPath := util.GetApiSocketPath()
	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}
	err = os.Chmod(socketPath, 0666)
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

func (s *APIServer) killExistingProcessIfNeeded() error {
	_, err := os.Stat(util.GetApiSocketPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	s.logger.Info("Killing existing process", zap.String("path", util.GetApiSocketPath()))

	_ = killExistingProcess()

	_ = os.Remove(util.GetPidPath())
	_ = os.Remove(util.GetApiSocketPath())

	return nil
}

func killExistingProcess() error {
	pidBytes, err := os.ReadFile(util.GetPidPath())
	if err == nil {
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
	err = process.Kill()
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

func (s *dotLocalServer) CreateMapping(ctx context.Context, req *api.CreateMappingRequest) (*api.Mapping, error) {
	mapping, err := s.dotlocal.CreateMapping(internal.MappingOptions{
		Host:       *req.Host,
		PathPrefix: *req.PathPrefix,
		Target:     *req.Target,
	})
	if err != nil {
		return nil, err
	}
	return mappingToApiMapping(mapping), nil
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

func (s *dotLocalServer) ListMappings(ctx context.Context, _ *emptypb.Empty) (*api.ListMappingsResponse, error) {
	res := &api.ListMappingsResponse{
		Mappings: lo.Map(s.dotlocal.GetMappings(), func(mapping internal.Mapping, _ int) *api.Mapping {
			return mappingToApiMapping(mapping)
		}),
	}
	return res, nil
}

func (s *dotLocalServer) GetSavedState(ctx context.Context, _ *emptypb.Empty) (*api.SavedState, error) {
	mappings := lo.Map(s.dotlocal.GetMappings(), func(mapping internal.Mapping, _ int) *api.Mapping {
		return mappingToApiMapping(mapping)
	})
	return &api.SavedState{
		Mappings:    mappings,
		Preferences: s.dotlocal.GetPreferences(),
	}, nil
}

func (s *dotLocalServer) SetPreferences(ctx context.Context, preferences *api.Preferences) (*emptypb.Empty, error) {
	err := s.dotlocal.SetPreferences(preferences)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *dotLocalServer) GetRootCertificate(ctx context.Context, _ *emptypb.Empty) (*api.GetRootCertificateResponse, error) {
	cert, err := s.dotlocal.caddy.getRootCertificate()
	if err != nil {
		return nil, err
	}
	return &api.GetRootCertificateResponse{
		Certificate: cert.Raw,
		NotBefore:   timestamppb.New(cert.NotBefore),
		NotAfter:    timestamppb.New(cert.NotAfter),
	}, nil
}

func mappingToApiMapping(mapping internal.Mapping) *api.Mapping {
	return &api.Mapping{
		Id:         &mapping.ID,
		Host:       &mapping.Host,
		PathPrefix: &mapping.PathPrefix,
		Target:     &mapping.Target,
		ExpiresAt: &timestamppb.Timestamp{
			Seconds: mapping.ExpresAt.Unix(),
		},
	}
}
