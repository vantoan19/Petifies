package valueobjects

type Visibility string

const (
	PublicVisibility    Visibility = "public"
	FollowersVisibility Visibility = "followers"
	PrivateVisilibity   Visibility = "private"
)
