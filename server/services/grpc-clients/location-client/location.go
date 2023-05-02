package locationclient

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"

	locationProtoV1 "github.com/vantoan19/Petifies/proto/location-service/v1"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/location-service/pkg/models"
	"github.com/vantoan19/Petifies/server/services/location-service/pkg/translators"
)

var logger = logging.New("LocationClient")

const locationService = "location_service.v1.LocationService"

type locationClient struct {
	listNearByLocationsByType endpoint.Endpoint
}

type LocationClient interface {
	ListNearByLocationsByType(ctx context.Context, req *models.ListNearByLocationsByTypeReq) (*models.ListNearByLocationsByTypeResp, error)
}

func New(conn *grpc.ClientConn) LocationClient {
	return &locationClient{
		listNearByLocationsByType: grpctransport.NewClient(
			conn,
			locationService,
			"ListNearByLocationsByType",
			translators.EncodeListNearByLocationsByTypeRequest,
			translators.DecodeListNearByLocationByTypeResponse,
			locationProtoV1.ListNearByLocationsByTypeResponse{},
		).Endpoint(),
	}
}

func (l *locationClient) ListNearByLocationsByType(ctx context.Context, req *models.ListNearByLocationsByTypeReq) (*models.ListNearByLocationsByTypeResp, error) {
	logger.Info("Start ListNearByLocationsByType")

	resp, err := l.listNearByLocationsByType(ctx, req)
	if err != nil {
		logger.ErrorData("Finish ListNearByLocationsByType: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListNearByLocationsByType: Successful")
	return resp.(*models.ListNearByLocationsByTypeResp), nil
}
