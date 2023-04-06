package valueobjects

type PetifiesStatus string

const (
	PetifiesUnavailable PetifiesStatus = "PETIFIES_STATUS_UNAVAILABLE"
	PetifiesOnASession  PetifiesStatus = "PETIFIES_STATUS_ON_A_SESSION"
	PetifiesDeleted     PetifiesStatus = "PETIFIES_STATUS_DELETED"
)
