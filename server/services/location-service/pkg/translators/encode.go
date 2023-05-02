package translators

import (
	"context"

	commonProto "github.com/vantoan19/Petifies/proto/common"
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
		PageSize:     int32(req.PageSize),
		Offset:       int32(req.Offset),
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

func getLocationType(locationType string) commonProto.LocationType {
	switch locationType {
	case "LOCATION_TYPE_PETIFIES_DOG_WALKING":
		return commonProto.LocationType_LOCATION_TYPE_PETIFIES_DOG_WALKING
	case "LOCATION_TYPE_PETIFIES_CAT_PLAYING":
		return commonProto.LocationType_LOCATION_TYPE_PETIFIES_CAT_PLAYING
	case "LOCATION_TYPE_PETIFIES_DOG_SITTING":
		return commonProto.LocationType_LOCATION_TYPE_PETIFIES_DOG_SITTING
	case "LOCATION_TYPE_PETIFIES_CAT_SITTING":
		return commonProto.LocationType_LOCATION_TYPE_PETIFIES_CAT_SITTING
	case "LOCATION_TYPE_PETIFIES_DOG_ADOPTION":
		return commonProto.LocationType_LOCATION_TYPE_PETIFIES_DOG_ADOPTION
	case "LOCATION_TYPE_PETIFIES_CAT_ADOPTION":
		return commonProto.LocationType_LOCATION_TYPE_PETIFIES_CAT_ADOPTION
	default:
		return commonProto.LocationType_LOCATION_UNKNOWN
	}
}
