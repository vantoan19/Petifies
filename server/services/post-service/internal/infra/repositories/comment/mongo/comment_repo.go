package mongo_comment

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentRepository struct {
	db *mongo.Client
}

func New(db *mongo.Client) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}
