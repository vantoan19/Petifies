package valueobjects

type RelationshipType string

const (
	FollowRelationship RelationshipType = "FOLLOW"
	FriendRelationship RelationshipType = "FRIEND"
)
