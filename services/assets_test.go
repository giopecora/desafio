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

func TestAssetService_CreateAsset(t *testing.T) {
	type fields struct {
		collection interfaces.CollectionInterface
	}
	type args struct {
		ctx       context.Context
		userID    string
		novoAsset models.Asset
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   models.Asset
		want1  error
	}{
		{
			name: "Success - Create Asset",
			fields: fields{
				collection: &mocks.MockCollection{
					InsertOneFunc: func(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
						return &mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}, nil
					},
				},
			},
			args: args{
				ctx:       context.Background(),
				userID:    "user123",
				novoAsset: models.Asset{Valor: 100, Tipo: "Imobiliario"},
			},
			want: func() models.Asset {
				expectedAsset := models.Asset{UserID: "user123", Valor: 100, Tipo: "Imobiliario"}
				return expectedAsset
			}(),
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AssetService{
				collection: tt.fields.collection,
			}
			got, got1 := s.CreateAsset(tt.args.ctx, tt.args.userID, tt.args.novoAsset)
			assert.Equalf(t, tt.want.UserID, got.UserID, "CreateAsset(%v, %v, %v)", tt.args.ctx, tt.args.userID, tt.args.novoAsset)
			assert.Equalf(t, tt.want.Valor, got.Valor, "CreateAsset(%v, %v, %v)", tt.args.ctx, tt.args.userID, tt.args.novoAsset)
			assert.Equalf(t, tt.want.Tipo, got.Tipo, "CreateAsset(%v, %v, %v)", tt.args.ctx, tt.args.userID, tt.args.novoAsset)
			assert.Equalf(t, tt.want1, got1, "CreateAsset(%v, %v, %v)", tt.args.ctx, tt.args.userID, tt.args.novoAsset)

		})
	}
}

func TestAssetService_DeleteAsset(t *testing.T) {
	type fields struct {
		collection interfaces.CollectionInterface
	}
	type args struct {
		ctx     context.Context
		userID  string
		assetID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   error
	}{
		{
			name: "Success - Delete Asset",
			fields: fields{
				collection: &mocks.MockCollection{
					DeleteOneFunc: func(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
						return &mongo.DeleteResult{DeletedCount: 1}, nil
					},
				},
			},
			args: args{
				ctx:     context.TODO(),
				userID:  "user123",
				assetID: primitive.NewObjectID().Hex(),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AssetService{
				collection: tt.fields.collection,
			}
			assert.Equalf(t, tt.want, s.DeleteAsset(tt.args.ctx, tt.args.userID, tt.args.assetID), "DeleteAsset(%v, %v, %v)", tt.args.ctx, tt.args.userID, tt.args.assetID)
		})
	}
}

func TestAssetService_UpdateAsset(t *testing.T) {
	type fields struct {
		collection interfaces.CollectionInterface
	}
	type args struct {
		ctx          context.Context
		userID       string
		assetID      string
		updatedAsset models.Asset
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   models.Asset
		want1  error
	}{
		{
			name: "Success - Update Asset",
			fields: fields{
				collection: &mocks.MockCollection{
					UpdateOneFunc: func(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
						return &mongo.UpdateResult{MatchedCount: 1}, nil
					},
				},
			},
			args: args{
				ctx:          context.TODO(),
				userID:       "validUserID",
				assetID:      primitive.NewObjectID().Hex(),
				updatedAsset: models.Asset{Valor: 2000, Tipo: "Investment"},
			},
			want:  models.Asset{UserID: "validUserID", ID: primitive.NewObjectID().Hex(), Valor: 2000, Tipo: "Investment"},
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AssetService{
				collection: tt.fields.collection,
			}
			got, got1 := s.UpdateAsset(tt.args.ctx, tt.args.userID, tt.args.assetID, tt.args.updatedAsset)

			assert.Equalf(t, tt.want.UserID, got.UserID, "UpdateAsset(%v, %v, %v, %v)", tt.args.ctx, tt.args.userID, tt.args.assetID, tt.args.updatedAsset)
			assert.Equalf(t, tt.want.Valor, got.Valor, "UpdateAsset(%v, %v, %v, %v)", tt.args.ctx, tt.args.userID, tt.args.assetID, tt.args.updatedAsset)
			assert.Equalf(t, tt.want.Tipo, got.Tipo, "UpdateAsset(%v, %v, %v, %v)", tt.args.ctx, tt.args.userID, tt.args.assetID, tt.args.updatedAsset)
			assert.Equalf(t, tt.want1, got1, "UpdateAsset(%v, %v, %v, %v)", tt.args.ctx, tt.args.userID, tt.args.assetID, tt.args.updatedAsset)
		})
	}
}
