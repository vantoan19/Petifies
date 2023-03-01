package mongo_post

import "go.mongodb.org/mongo-driver/mongo"

type PostRepository struct {
	db *mongo.Client
}
