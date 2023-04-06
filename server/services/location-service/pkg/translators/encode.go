package translators

import (
	"context"

	locationProtoV1 "github.com/vantoan19/Petifies/proto/location-service/v1"
	"github.com/vantoan19/Petifies/server/services/location-service/pkg/models"
)

func EncodeListNearByLocationsByTypeRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.ListNearByLocationsByTypeReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &locationProtoV1.ListNearByLocationsByTypeRequest{
		LocationType: getLocationType(req.LocationType),
		Longitude:    req.Longitude,
		Latitude:     req.Latitude,
		Radius:       req.Radius,
	}, nil
}

func EncodeListNearByLocationsByTypeResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.ListNearByLocationsByTypeResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	var locations []*locationProtoV1.Location
	for _, l := range resp.Locations {
		location := encodeLocationModel(l)
		locations = append(locations, location)
	}

	return &locationProtoV1.ListNearByLocationsByTypeResponse{
		Locations: locations,
	}, nil
}

func encodeLocationModel(location *models.Location) *locationProtoV1.Location {
	return &locationProtoV1.Location{
		Id:           location.ID.String(),
		EntityId:     location.EntityID.String(),
		LocationType: getLocationType(location.LocationType),
	}
}

func getLocationType(locationType string) locationProtoV1.LocationType {
	switch locationType {
	case "PETIFIES":
		return locationProtoV1.LocationType_LOCATION_TYPE_PETIFIES
	default:
		return locationProtoV1.LocationType_LOCATION_UNKNOWN
	}
}
