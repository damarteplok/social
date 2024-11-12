package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/damarteplok/social/internal/store"
	"github.com/go-redis/redis/v8"
)

type ReviewSomethingStore struct {
	rdb *redis.Client
}

const ReviewSomethingExpTime = time.Hour * 24 * 7
	
func (s *ReviewSomethingStore) Get(ctx context.Context, modelID int64) (*store.ReviewSomething, error) {
	cacheKey := fmt.Sprintf("ReviewSomething-%v", modelID)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var model store.ReviewSomething
	if data != "" {
		err := json.Unmarshal([]byte(data), &model)
		if err != nil {
			return nil, err
		}
	}

	return &model, nil
}

func (s *ReviewSomethingStore) Set(ctx context.Context, model *store.ReviewSomething) error {
	cacheKey := fmt.Sprintf("ReviewSomething-%v", model.ID)

	json, err := json.Marshal(model)
	if err != nil {
		return err
	}
	return s.rdb.SetEX(ctx, cacheKey, json, ReviewSomethingExpTime).Err()
}

func (s *ReviewSomethingStore) Delete(ctx context.Context, modelID int64) {
	cacheKey := fmt.Sprintf("ReviewSomething-%v", modelID)
	s.rdb.Del(ctx, cacheKey)
}
