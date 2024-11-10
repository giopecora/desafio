package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"main/cache"
	"main/services"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type ScoreHandler struct {
	scoreService *services.ScoreService
	redisClient  cache.RedisClient
	cacheTTL     time.Duration
}

func NewScoreHandler(scoreService *services.ScoreService, redisClient cache.RedisClient) *ScoreHandler {
	return &ScoreHandler{
		scoreService: scoreService,
		redisClient:  redisClient,
		cacheTTL:     time.Minute * 15,
	}
}

func (h *ScoreHandler) GetUserScoreHandler(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	cacheKey := fmt.Sprintf("user_score:%s", userID)

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	scoreChan := make(chan int)
	errChan := make(chan error)

	go func() {
		cachedScore, err := h.redisClient.Get(ctx, cacheKey)
		if err != nil {
			cachedScore = ""
			return
		}

		if cachedScore != "" {
			var score int
			if err := json.Unmarshal([]byte(cachedScore), &score); err == nil {
				scoreChan <- score
			} else {
				errChan <- fmt.Errorf("cache unmarshal error: %v", err)
			}
		} else {
			errChan <- fmt.Errorf("score not found in cache")
		}
	}()

	go func() {
		score, err := h.scoreService.GetUserScore(ctx, userID)
		if err != nil {
			errChan <- fmt.Errorf("error calculating score: %v", err)
			return
		}
		scoreChan <- score
	}()

	select {
	case score := <-scoreChan:
		responseJSON(w, http.StatusOK, map[string]int{"score": score})
	case err := <-errChan:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	case <-ctx.Done():
		http.Error(w, "Request timed out", http.StatusRequestTimeout)
	}
}
