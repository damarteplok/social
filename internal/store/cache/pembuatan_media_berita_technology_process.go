package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/damarteplok/social/internal/store"
	"github.com/go-redis/redis/v8"
)

type PembuatanMediaBeritaTechnologyStore struct {
	rdb *redis.Client
}

const PembuatanMediaBeritaTechnologyExpTime = time.Hour * 24 * 7
	
func (s *PembuatanMediaBeritaTechnologyStore) Get(ctx context.Context, modelID int64) (*store.PembuatanMediaBeritaTechnology, error) {
	cacheKey := fmt.Sprintf("PembuatanMediaBeritaTechnology-%v", modelID)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var model store.PembuatanMediaBeritaTechnology
	if data != "" {
		err := json.Unmarshal([]byte(data), &model)
		if err != nil {
			return nil, err
		}
	}

	return &model, nil
}

func (s *PembuatanMediaBeritaTechnologyStore) Set(ctx context.Context, model *store.PembuatanMediaBeritaTechnology) error {
	cacheKey := fmt.Sprintf("PembuatanMediaBeritaTechnology-%v", model.ID)

	json, err := json.Marshal(model)
	if err != nil {
		return err
	}
	return s.rdb.SetEX(ctx, cacheKey, json, PembuatanMediaBeritaTechnologyExpTime).Err()
}

func (s *PembuatanMediaBeritaTechnologyStore) Delete(ctx context.Context, modelID int64) {
	cacheKey := fmt.Sprintf("PembuatanMediaBeritaTechnology-%v", modelID)
	s.rdb.Del(ctx, cacheKey)
}
