package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	postclient "github.com/vantoan19/Petifies/server/services/grpc-clients/post-client"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/domain/repositories"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

type redisCommentCacheRepository struct {
	client     *redis.Client
	postClient postclient.PostClient
}

func NewRedisCommentCacheRepository(client *redis.Client, postClient postclient.PostClient) repositories.CommentCacheRepository {
	return &redisCommentCacheRepository{client: client, postClient: postClient}
}

func (r *redisCommentCacheRepository) GetCommentContent(ctx context.Context, commentID uuid.UUID) (*models.Comment, error) {
	key := fmt.Sprintf("comment:%s:content", commentID.String())
	commentContentStr, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		comment, err := r.postClient.GetComment(ctx, &models.GetCommentReq{
			CommentID: commentID,
		})
		if err != nil {
			return nil, err
		}

		go func() {
			r.SetCommentContent(ctx, commentID, *comment)
		}()

		return comment, nil
	} else if err != nil {
		return nil, err
	} else {
		var commentContent models.Comment
		err = json.Unmarshal([]byte(commentContentStr), &commentContent)
		if err != nil {
			return nil, err
		}

		return &commentContent, nil
	}
}

func (r *redisCommentCacheRepository) ListCommentContents(ctx context.Context, commentIDs []uuid.UUID) ([]*models.Comment, error) {
	if len(commentIDs) == 0 {
		return []*models.Comment{}, nil
	}

	var keys []string
	for _, id := range commentIDs {
		keys = append(keys, fmt.Sprintf("comment:%s:content", id.String()))
	}

	commentStrs, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	results := make(map[string]*models.Comment)
	commentsToList := make([]uuid.UUID, 0)
	for i, key := range keys {
		commentStr := commentStrs[i]
		if commentStr == nil {
			commentsToList = append(commentsToList, commentIDs[i])
		} else {
			var commentContent models.Comment
			err = json.Unmarshal([]byte(commentStr.(string)), &commentContent)
			if err != nil {
				return nil, err
			}
			results[key] = &commentContent
		}
	}

	if len(commentsToList) > 0 {
		commentsResp, err := r.postClient.ListComments(ctx, &models.ListCommentsReq{
			CommentIDs: commentsToList,
		})
		if err != nil {
			return nil, err
		}
		for _, c := range commentsResp.Comments {
			key := fmt.Sprintf("comment:%s:content", c.ID.String())
			results[key] = c
		}

		go func() {
			commentsToCache := make(map[string]string)
			for _, c := range commentsResp.Comments {
				key := fmt.Sprintf("comment:%s:content", c.ID.String())
				commentContentStr, err := json.Marshal(c)
				if err != nil {
					return
				}
				commentsToCache[key] = string(commentContentStr)
			}
			tx := r.client.TxPipeline()
			tx.MSet(ctx, commentsToCache)
			_, err = tx.Exec(ctx)
			if err != nil {
				return
			}
		}()
	}

	finalResults := make([]*models.Comment, 0)
	for _, id := range commentIDs {
		key := fmt.Sprintf("comment:%s:content", id.String())
		if val, ok := results[key]; ok {
			finalResults = append(finalResults, val)
		}
	}

	return finalResults, nil
}

func (r *redisCommentCacheRepository) SetCommentContent(ctx context.Context, commentID uuid.UUID, comment models.Comment) error {
	key := fmt.Sprintf("comment:%s:content", commentID.String())

	commentContentStr, err := json.Marshal(comment)
	if err != nil {
		return err
	}
	tx := r.client.TxPipeline()

	tx.Set(ctx, key, commentContentStr, 0)

	_, err = tx.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *redisCommentCacheRepository) ExistsCommentContent(ctx context.Context, commentID uuid.UUID) (bool, error) {
	key := fmt.Sprintf("comment:%s:content", commentID.String())
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}
