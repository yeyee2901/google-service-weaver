package unitedbpapi

import (
	"context"
	"time"

	"github.com/ServiceWeaver/weaver"
	unitepb "github.com/yeyee2901/unitedb-api-proto/gen/go/unitedb/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// UniteDBService adalah service untuk berkomunikasi dengan API unitedb
type UniteDBService interface {
	// GetBattleItem digunakan untuk mengambil battle item
	GetBattleItem(c context.Context, name, tier string) (*unitepb.GetBattleItemResponse, error)
}

type grpcUniteDBService struct {
	weaver.Implements[UniteDBService]
	weaver.WithConfig[config]

	conn *grpc.ClientConn
}

// TOML Config section : <path package lengkap>/<nama interface>
// `github.com/yeyee2901/service-weaver/service/unitedbapi/UniteDBService`
type config struct {
	GrpcAddress string
}

// Init adalah function pertama yang dijalankan setelah component di
// instantiate
func (s *grpcUniteDBService) Init(_ context.Context) error {
	cc, err := grpc.Dial(
		s.Config().GrpcAddress,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		return err
	}

	s.conn = cc

	return nil
}

func (s *grpcUniteDBService) GetBattleItem(c context.Context, name string, tier string) (*unitepb.GetBattleItemResponse, error) {
	client := unitepb.NewUniteDBServiceClient(s.conn)
	req := &unitepb.GetBattleItemRequest{}

	if name != "" {
		req.Name = &name
	}

	if tier != "" {
		req.Tier = &tier
	}

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	// execute get battle item
	resp, err := client.GetBattleItem(ctx, req)

	// NOTE: error adalah non-serializable, jadi walaupun di wrap seperti
	// ini tidak akan pengaruh
	//  return nil, UniteDBError{reason: err}

	// jadi mending di bagi saja domain nya berdasarkan caller & called
	// function, cukup membedakan asal error dari caller component atau
	// called component
	if err != nil {
		return nil, err
	}

	return resp, nil
}
