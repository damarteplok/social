package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/damarteplok/social/internal/store"
	"github.com/go-redis/redis/v8"
)

type ApprovingArtikelStore struct {
	rdb *redis.Client
}

const ApprovingArtikelExpTime = time.Hour * 24 * 7
	
func (s *ApprovingArtikelStore) Get(ctx context.Context, modelID int64) (*store.ApprovingArtikel, error) {
	cacheKey := fmt.Sprintf("ApprovingArtikel-%v", modelID)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var model store.ApprovingArtikel
	if data != "" {
		err := json.Unmarshal([]byte(data), &model)
		if err != nil {
			return nil, err
		}
	}

	return &model, nil
}

func (s *ApprovingArtikelStore) Set(ctx context.Context, model *store.ApprovingArtikel) error {
	cacheKey := fmt.Sprintf("ApprovingArtikel-%v", model.ID)

	json, err := json.Marshal(model)
	if err != nil {
		return err
	}
	return s.rdb.SetEX(ctx, cacheKey, json, ApprovingArtikelExpTime).Err()
}

func (s *ApprovingArtikelStore) Delete(ctx context.Context, modelID int64) {
	cacheKey := fmt.Sprintf("ApprovingArtikel-%v", modelID)
	s.rdb.Del(ctx, cacheKey)
}
