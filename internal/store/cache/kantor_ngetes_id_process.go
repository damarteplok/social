package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/damarteplok/social/internal/store"
	"github.com/go-redis/redis/v8"
)

type KantorNgetesIdStore struct {
	rdb *redis.Client
}

const KantorNgetesIdExpTime = time.Hour * 24 * 7
	
func (s *KantorNgetesIdStore) Get(ctx context.Context, modelID int64) (*store.KantorNgetesId, error) {
	cacheKey := fmt.Sprintf("KantorNgetesId-%v", modelID)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var model store.KantorNgetesId
	if data != "" {
		err := json.Unmarshal([]byte(data), &model)
		if err != nil {
			return nil, err
		}
	}

	return &model, nil
}

func (s *KantorNgetesIdStore) Set(ctx context.Context, model *store.KantorNgetesId) error {
	cacheKey := fmt.Sprintf("KantorNgetesId-%v", model.ID)

	json, err := json.Marshal(model)
	if err != nil {
		return err
	}
	return s.rdb.SetEX(ctx, cacheKey, json, KantorNgetesIdExpTime).Err()
}

func (s *KantorNgetesIdStore) Delete(ctx context.Context, modelID int64) {
	cacheKey := fmt.Sprintf("KantorNgetesId-%v", modelID)
	s.rdb.Del(ctx, cacheKey)
}
