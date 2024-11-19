package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/damarteplok/social/internal/store"
	"github.com/go-redis/redis/v8"
)

type PembuatanArtikelStore struct {
	rdb *redis.Client
}

const PembuatanArtikelExpTime = time.Hour * 24 * 7
	
func (s *PembuatanArtikelStore) Get(ctx context.Context, modelID int64) (*store.PembuatanArtikel, error) {
	cacheKey := fmt.Sprintf("PembuatanArtikel-%v", modelID)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var model store.PembuatanArtikel
	if data != "" {
		err := json.Unmarshal([]byte(data), &model)
		if err != nil {
			return nil, err
		}
	}

	return &model, nil
}

func (s *PembuatanArtikelStore) Set(ctx context.Context, model *store.PembuatanArtikel) error {
	cacheKey := fmt.Sprintf("PembuatanArtikel-%v", model.ID)

	json, err := json.Marshal(model)
	if err != nil {
		return err
	}
	return s.rdb.SetEX(ctx, cacheKey, json, PembuatanArtikelExpTime).Err()
}

func (s *PembuatanArtikelStore) Delete(ctx context.Context, modelID int64) {
	cacheKey := fmt.Sprintf("PembuatanArtikel-%v", modelID)
	s.rdb.Del(ctx, cacheKey)
}
