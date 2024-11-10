package services

import (
	"context"
	"main/models"
	"math"
	"time"
)

type ScoreService struct {
	AssetService *AssetService
	DebtService  *DebtService
}

func NewScoreService(assetService *AssetService, debtService *DebtService) *ScoreService {
	return &ScoreService{
		AssetService: assetService,
		DebtService:  debtService,
	}
}

func (s *ScoreService) CalculateScore(assets []models.Asset, debts []models.Debt) int {
	if len(assets) == 0 && len(debts) == 0 {
		return 500
	}

	totalAssetValue := 0.0
	assetCount := len(assets)
	totalDebtValue := 0.0
	debtCount := len(debts)

	for _, asset := range assets {
		totalAssetValue += asset.Valor
	}

	for _, debt := range debts {
		totalDebtValue += debt.Valor
	}

	assetScore := math.Log(totalAssetValue+1) * math.Log(float64(assetCount)+1)

	debtScore := math.Log(totalDebtValue+1) * math.Log(float64(debtCount)+1) * 1.5

	rawScore := 500 + (assetScore-debtScore)*50

	return int(math.Max(0, math.Min(1000, rawScore)))
}

func (s *ScoreService) GetUserScore(ctx context.Context, userID string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	assets, err := s.AssetService.GetAssets(ctx, userID)
	if err != nil {
		return 0, err
	}

	debts, err := s.DebtService.GetDebts(ctx, userID)
	if err != nil {
		return 0, err
	}

	return s.CalculateScore(assets, debts), nil
}
