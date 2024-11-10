package services

import (
	"context"
	"errors"
	"fmt"
	"main/interfaces"
	"main/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrInvalidAssetID = errors.New("invalid asset ID")
	ErrAssetNotFound  = errors.New("asset not found or unauthorized")
)

type AssetService struct {
	collection interfaces.CollectionInterface
}

func NewAssetService(client *mongo.Client) *AssetService {
	return &AssetService{
		collection: client.Database("desafio").Collection("assets"),
	}
}

func (s *AssetService) CreateAsset(ctx context.Context, userID string, novoAsset models.Asset) (models.Asset, error) {
	novoAsset.UserID = userID

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := s.collection.InsertOne(ctx, novoAsset)
	if err != nil {
		return models.Asset{}, err
	}

	novoAsset.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return novoAsset, nil
}

func (s *AssetService) UpdateAsset(ctx context.Context, userID string, assetID string, updatedAsset models.Asset) (models.Asset, error) {
	assetObjectID, err := primitive.ObjectIDFromHex(assetID)
	if err != nil {
		return models.Asset{}, ErrInvalidAssetID
	}
	fmt.Println(userID)

	filter := bson.M{"_id": assetObjectID, "userID": userID}
	update := bson.M{"$set": bson.M{"valor": updatedAsset.Valor, "tipo": updatedAsset.Tipo}}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return models.Asset{}, err
	}

	if result.MatchedCount == 0 {
		return models.Asset{}, ErrAssetNotFound
	}

	updatedAsset.ID = assetID
	updatedAsset.UserID = userID
	return updatedAsset, nil
}

func (s *AssetService) DeleteAsset(ctx context.Context, userID string, assetID string) error {
	assetObjectID, err := primitive.ObjectIDFromHex(assetID)
	if err != nil {
		return ErrInvalidAssetID
	}

	filter := bson.M{"_id": assetObjectID, "userID": userID}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrAssetNotFound
	}

	return nil
}

func (s *AssetService) GetAssets(ctx context.Context, userID string) ([]models.Asset, error) {
	filter := bson.M{"userID": userID}
	var assets []models.Asset

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &assets); err != nil {
		return nil, err
	}

	return assets, nil
}
