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

type redisPostCacheRepository struct {
	client     *redis.Client
	postClient postclient.PostClient
}

func NewRedisPostCacheRepository(client *redis.Client, postClient postclient.PostClient) repositories.PostCacheRepository {
	return &redisPostCacheRepository{client: client, postClient: postClient}
}

func (r *redisPostCacheRepository) GetPostContent(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	key := fmt.Sprintf("post:%s:content", postID.String())
	postContentStr, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		post, err := r.postClient.GetPost(ctx, &models.GetPostReq{
			PostID: postID,
		})
		if err != nil {
			return nil, err
		}

		go func() {
			r.SetPostContent(ctx, postID, *post)
		}()

		return post, nil
	} else if err != nil {
		return nil, err
	} else {
		var postContent models.Post
		err = json.Unmarshal([]byte(postContentStr), &postContent)
		if err != nil {
			return nil, err
		}

		return &postContent, nil
	}

}

func (r *redisPostCacheRepository) ListPostContents(ctx context.Context, postIds []uuid.UUID) ([]*models.Post, error) {
	if len(postIds) == 0 {
		return []*models.Post{}, nil
	}

	var keys []string
	for _, id := range postIds {
		keys = append(keys, fmt.Sprintf("post:%s:content", id.String()))
	}

	postStrs, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	results := make(map[string]*models.Post)
	postsToList := make([]uuid.UUID, 0)
	for i, key := range keys {
		postStr := postStrs[i]
		if postStr == nil {
			postsToList = append(postsToList, postIds[i])
		} else {
			var postContent models.Post
			err = json.Unmarshal([]byte(postStr.(string)), &postContent)
			if err != nil {
				return nil, err
			}
			results[key] = &postContent
		}
	}

	if len(postsToList) > 0 {
		postsResp, err := r.postClient.ListPosts(ctx, &models.ListPostsReq{
			PostIDs: postsToList,
		})
		if err != nil {
			return nil, err
		}

		for _, c := range postsResp.Posts {
			key := fmt.Sprintf("post:%s:content", c.ID.String())
			results[key] = c
		}

		go func() {
			postsToCache := make(map[string]string)
			for _, p := range postsResp.Posts {
				key := fmt.Sprintf("post:%s:content", p.ID.String())
				postContentStr, err := json.Marshal(p)
				if err != nil {
					return
				}
				postsToCache[key] = string(postContentStr)
			}
			tx := r.client.TxPipeline()
			tx.MSet(ctx, postsToCache)
			_, err = tx.Exec(ctx)
			if err != nil {
				return
			}
		}()
	}

	finalResults := make([]*models.Post, 0)
	for _, id := range postIds {
		key := fmt.Sprintf("post:%s:content", id.String())
		if val, ok := results[key]; ok {
			finalResults = append(finalResults, val)
		}
	}

	return finalResults, nil
}

func (r *redisPostCacheRepository) SetPostContent(ctx context.Context, postID uuid.UUID, post models.Post) error {
	key := fmt.Sprintf("post:%s:content", postID.String())

	postContentStr, err := json.Marshal(post)
	if err != nil {
		return err
	}

	tx := r.client.TxPipeline()

	tx.Set(ctx, key, postContentStr, 0)

	_, err = tx.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *redisPostCacheRepository) ExistsPostContent(ctx context.Context, postID uuid.UUID) (bool, error) {
	key := fmt.Sprintf("post:%s:content", postID.String())
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}
