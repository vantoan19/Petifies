package models

type ListNearByLocationsByTypeReq struct {
	LocationType string
	Longitude    float64
	Latitude     float64
	Radius       float64
	PageSize     int
	Offset       int
}
