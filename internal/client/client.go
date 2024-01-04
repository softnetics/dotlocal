package client

import (
	api "github.com/softnetics/dotlocal/internal/api/proto"
	"github.com/softnetics/dotlocal/internal/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewApiClient() (api.DotLocalClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial("passthrough:///unix://"+util.GetApiSocketPath(), opts...)
	if err != nil {
		return nil, err
	}
	client := api.NewDotLocalClient(conn)
	return client, nil
}
