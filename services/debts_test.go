package services

import (
	"context"
	"main/interfaces"
	"main/models"
	"main/services/test/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestDebtService_CreateDebt(t *testing.T) {
	type fields struct {
		collection interfaces.CollectionInterface
	}
	type args struct {
		ctx      context.Context
		userID   string
		novaDebt models.Debt
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   models.Debt
		want1  error
	}{
		{
			name: "Success - Create Debts",
			fields: fields{
				collection: &mocks.MockCollection{
					InsertOneFunc: func(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
						return &mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}, nil
					},
				},
			},
			args: args{
				ctx:      context.Background(),
				userID:   "user123",
				novaDebt: models.Debt{Valor: 100, Tipo: "Imobiliario"},
			},
			want: func() models.Debt {
				expectedDebt := models.Debt{UserID: "user123", Valor: 100, Tipo: "Imobiliario"}
				return expectedDebt
			}(),
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DebtService{
				collection: tt.fields.collection,
			}
			got, got1 := s.CreateDebt(tt.args.ctx, tt.args.userID, tt.args.novaDebt)
			assert.Equalf(t, tt.want.UserID, got.UserID, "CreateDebt(%v, %v, %v)", tt.args.ctx, tt.args.userID, tt.args.novaDebt)
			assert.Equalf(t, tt.want.Valor, got.Valor, "CreateDebt(%v, %v, %v)", tt.args.ctx, tt.args.userID, tt.args.novaDebt)
			assert.Equalf(t, tt.want.Tipo, got.Tipo, "CreateDebt(%v, %v, %v)", tt.args.ctx, tt.args.userID, tt.args.novaDebt)
			assert.Equalf(t, tt.want1, got1, "CreateDebt(%v, %v, %v)", tt.args.ctx, tt.args.userID, tt.args.novaDebt)
		})
	}
}

func TestDebtService_DeleteDebt(t *testing.T) {
	type fields struct {
		collection interfaces.CollectionInterface
	}
	type args struct {
		ctx    context.Context
		userID string
		debtID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   error
	}{
		{
			name: "Success - Delete Debt",
			fields: fields{
				collection: &mocks.MockCollection{
					DeleteOneFunc: func(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
						return &mongo.DeleteResult{DeletedCount: 1}, nil
					},
				},
			},
			args: args{
				ctx:    context.TODO(),
				userID: "user123",
				debtID: primitive.NewObjectID().Hex(),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DebtService{
				collection: tt.fields.collection,
			}
			assert.Equalf(t, tt.want, s.DeleteDebt(tt.args.ctx, tt.args.debtID), "DeleteDebt(%v, %v)", tt.args.ctx, tt.args.debtID)
		})
	}
}

func TestDebtService_UpdateDebt(t *testing.T) {
	type fields struct {
		collection interfaces.CollectionInterface
	}
	type args struct {
		ctx         context.Context
		userID      string
		debtID      string
		updatedDebt models.Debt
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   models.Debt
		want1  error
	}{
		{
			name: "Success - Update Debt",
			fields: fields{
				collection: &mocks.MockCollection{
					UpdateOneFunc: func(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
						return &mongo.UpdateResult{MatchedCount: 1}, nil
					},
				},
			},
			args: args{
				ctx:         context.TODO(),
				debtID:      primitive.NewObjectID().Hex(),
				updatedDebt: models.Debt{Valor: 2000, Tipo: "Investment"},
			},
			want:  models.Debt{UserID: "validUserID", ID: primitive.NewObjectID().Hex(), Valor: 2000, Tipo: "Investment"},
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DebtService{
				collection: tt.fields.collection,
			}

			got, got1 := s.UpdateDebt(tt.args.ctx, tt.args.debtID, tt.args.updatedDebt)
			assert.Equalf(t, tt.want.Valor, got.Valor, "UpdateDebt(%v, %v, %v)", tt.args.ctx, tt.args.debtID, tt.args.updatedDebt)
			assert.Equalf(t, tt.want.Tipo, got.Tipo, "UpdateDebt(%v, %v, %v)", tt.args.ctx, tt.args.debtID, tt.args.updatedDebt)
			assert.Equalf(t, tt.want1, got1, "UpdateDebt(%v, %v, %v, %v)", tt.args.ctx, tt.args.userID, tt.args.debtID, tt.args.updatedDebt)
		})
	}
}
