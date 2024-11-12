package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/damarteplok/social/internal/store"
	"github.com/go-redis/redis/v8"
)

type BikinSomethingStore struct {
	rdb *redis.Client
}

const BikinSomethingExpTime = time.Hour * 24 * 7
	
func (s *BikinSomethingStore) Get(ctx context.Context, modelID int64) (*store.BikinSomething, error) {
	cacheKey := fmt.Sprintf("BikinSomething-%v", modelID)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var model store.BikinSomething
	if data != "" {
		err := json.Unmarshal([]byte(data), &model)
		if err != nil {
			return nil, err
		}
	}

	return &model, nil
}

func (s *BikinSomethingStore) Set(ctx context.Context, model *store.BikinSomething) error {
	cacheKey := fmt.Sprintf("BikinSomething-%v", model.ID)

	json, err := json.Marshal(model)
	if err != nil {
		return err
	}
	return s.rdb.SetEX(ctx, cacheKey, json, BikinSomethingExpTime).Err()
}

func (s *BikinSomethingStore) Delete(ctx context.Context, modelID int64) {
	cacheKey := fmt.Sprintf("BikinSomething-%v", modelID)
	s.rdb.Del(ctx, cacheKey)
}
