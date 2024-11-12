package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/damarteplok/social/internal/store"
	"github.com/go-redis/redis/v8"
)

type SetujuiSomethingStore struct {
	rdb *redis.Client
}

const SetujuiSomethingExpTime = time.Hour * 24 * 7
	
func (s *SetujuiSomethingStore) Get(ctx context.Context, modelID int64) (*store.SetujuiSomething, error) {
	cacheKey := fmt.Sprintf("SetujuiSomething-%v", modelID)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var model store.SetujuiSomething
	if data != "" {
		err := json.Unmarshal([]byte(data), &model)
		if err != nil {
			return nil, err
		}
	}

	return &model, nil
}

func (s *SetujuiSomethingStore) Set(ctx context.Context, model *store.SetujuiSomething) error {
	cacheKey := fmt.Sprintf("SetujuiSomething-%v", model.ID)

	json, err := json.Marshal(model)
	if err != nil {
		return err
	}
	return s.rdb.SetEX(ctx, cacheKey, json, SetujuiSomethingExpTime).Err()
}

func (s *SetujuiSomethingStore) Delete(ctx context.Context, modelID int64) {
	cacheKey := fmt.Sprintf("SetujuiSomething-%v", modelID)
	s.rdb.Del(ctx, cacheKey)
}
