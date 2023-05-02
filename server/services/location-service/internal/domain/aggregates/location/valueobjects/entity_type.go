package valueobjects

type EntityType string

const (
	PetifiesDogWalking  EntityType = "LOCATION_TYPE_PETIFIES_DOG_WALKING"
	PetifiesCatPlaying  EntityType = "LOCATION_TYPE_PETIFIES_CAT_PLAYING"
	PetifiesDogSitting  EntityType = "LOCATION_TYPE_PETIFIES_DOG_SITTING"
	PetifiesCatSitting  EntityType = "LOCATION_TYPE_PETIFIES_CAT_SITTING"
	PetifiesDogAdoption EntityType = "LOCATION_TYPE_PETIFIES_DOG_ADOPTION"
	PetifiesCatAdoption EntityType = "LOCATION_TYPE_PETIFIES_CAT_ADOPTION"
	UnknownType         EntityType = "UNKNOWN"
)
