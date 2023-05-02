package valueobjects

type PetifiesStatus string

const (
	PetifiesUnavailable PetifiesStatus = "PETIFIES_STATUS_UNAVAILABLE"
	PetifiesAvailable   PetifiesStatus = "PETIFIES_STATUS_AVAILABLE"
	PetifiesDeleted     PetifiesStatus = "PETIFIES_STATUS_DELETED"
)
