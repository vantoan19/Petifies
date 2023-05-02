package translators

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	locationProtoV1 "github.com/vantoan19/Petifies/proto/location-service/v1"
	"github.com/vantoan19/Petifies/server/services/location-service/pkg/models"
)

var (
	MustBeProtoReqErr     = status.Error(codes.InvalidArgument, "must be proto request")
	MustBeProtoRespErr    = status.Error(codes.InvalidArgument, "must be proto response")
	MustBeEndpointReqErr  = status.Error(codes.InvalidArgument, "must be endpoints' request")
	MustBeEndpointRespErr = status.Error(codes.InvalidArgument, "must be endpoints' response")
)

func DecodeListNearByLocationByTypeRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*locationProtoV1.ListNearByLocationsByTypeRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &models.ListNearByLocationsByTypeReq{
		LocationType: req.LocationType.String(),
		Longitude:    req.Longitude,
		Latitude:     req.Latitude,
		Radius:       req.Radius,
		PageSize:     int(req.PageSize),
		Offset:       int(req.Offset),
	}, nil
}

func DecodeListNearByLocationByTypeResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*locationProtoV1.ListNearByLocationsByTypeResponse)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	var locations []*models.Location
	for _, l := range resp.Locations {
		location, err := decodeLocationProtoModel(l)
		if err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}

	return &models.ListNearByLocationsByTypeResp{
		Locations: locations,
	}, nil
}

func decodeLocationProtoModel(location *locationProtoV1.Location) (*models.Location, error) {
	id, err := uuid.Parse(location.GetId())
	if err != nil {
		return nil, err
	}
	entityID, err := uuid.Parse(location.GetEntityId())
	if err != nil {
		return nil, err
	}

	return &models.Location{
		ID:           id,
		EntityID:     entityID,
		LocationType: location.GetLocationType().String(),
	}, nil
}
