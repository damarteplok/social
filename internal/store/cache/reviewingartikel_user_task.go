package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/damarteplok/social/internal/store"
	"github.com/go-redis/redis/v8"
)

type ReviewingArtikelStore struct {
	rdb *redis.Client
}

const ReviewingArtikelExpTime = time.Hour * 24 * 7
	
func (s *ReviewingArtikelStore) Get(ctx context.Context, modelID int64) (*store.ReviewingArtikel, error) {
	cacheKey := fmt.Sprintf("ReviewingArtikel-%v", modelID)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var model store.ReviewingArtikel
	if data != "" {
		err := json.Unmarshal([]byte(data), &model)
		if err != nil {
			return nil, err
		}
	}

	return &model, nil
}

func (s *ReviewingArtikelStore) Set(ctx context.Context, model *store.ReviewingArtikel) error {
	cacheKey := fmt.Sprintf("ReviewingArtikel-%v", model.ID)

	json, err := json.Marshal(model)
	if err != nil {
		return err
	}
	return s.rdb.SetEX(ctx, cacheKey, json, ReviewingArtikelExpTime).Err()
}

func (s *ReviewingArtikelStore) Delete(ctx context.Context, modelID int64) {
	cacheKey := fmt.Sprintf("ReviewingArtikel-%v", modelID)
	s.rdb.Del(ctx, cacheKey)
}
