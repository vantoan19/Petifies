package valueobjects

type LocationStatus string

const (
	LocationAvailable   LocationStatus = "AVAILABLE"
	LocationDeleted     LocationStatus = "DELETED"
	LocationUnavailable LocationStatus = "UNAVAILABLE"
)
