package services

import (
	"main/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScoreService_CalculateScore(t *testing.T) {
	type fields struct {
		AssetService *AssetService
		DebtService  *DebtService
	}
	type args struct {
		assets []models.Asset
		debts  []models.Debt
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "Success - Calculate score with assets and debts",
			fields: fields{
				AssetService: &AssetService{},
				DebtService:  &DebtService{},
			},
			args: args{
				assets: []models.Asset{
					{Valor: 100, Tipo: "Imobiliario"},
					{Valor: 200, Tipo: "Financeiro"},
				},
				debts: []models.Debt{
					{Valor: 50},
					{Valor: 100},
				},
			},
			want: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ScoreService{
				AssetService: tt.fields.AssetService,
				DebtService:  tt.fields.DebtService,
			}
			assert.Equalf(t, tt.want, s.CalculateScore(tt.args.assets, tt.args.debts), "CalculateScore(%v, %v)", tt.args.assets, tt.args.debts)
		})
	}
}
