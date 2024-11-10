package services

import (
	"context"
	"errors"
	"main/interfaces"
	"main/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrInvalidDebtID = errors.New("invalid debt ID")
	ErrDebtNotFound  = errors.New("debt not found or unauthorized")
)

type DebtService struct {
	collection interfaces.CollectionInterface
}

func NewDebtService(client *mongo.Client) *DebtService {
	return &DebtService{
		collection: client.Database("desafio").Collection("debts"),
	}
}

func (s *DebtService) CreateDebt(ctx context.Context, userID string, novaDebt models.Debt) (models.Debt, error) {
	novaDebt.UserID = userID

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := s.collection.InsertOne(ctx, novaDebt)
	if err != nil {
		return models.Debt{}, err
	}

	novaDebt.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return novaDebt, nil
}

func (s *DebtService) UpdateDebt(ctx context.Context, debtID string, updatedDebt models.Debt) (models.Debt, error) {
	debtObjectID, err := primitive.ObjectIDFromHex(debtID)
	if err != nil {
		return models.Debt{}, ErrInvalidDebtID
	}

	filter := bson.M{"_id": debtObjectID}
	update := bson.M{"$set": bson.M{"valor": updatedDebt.Valor, "tipo": updatedDebt.Tipo}}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return models.Debt{}, err
	}

	if result.MatchedCount == 0 {
		return models.Debt{}, ErrDebtNotFound
	}

	updatedDebt.ID = debtID
	return updatedDebt, nil
}

func (s *DebtService) DeleteDebt(ctx context.Context, debtID string) error {
	debtObjectID, err := primitive.ObjectIDFromHex(debtID)
	if err != nil {
		return ErrInvalidDebtID
	}

	filter := bson.M{"_id": debtObjectID}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrDebtNotFound
	}

	return nil
}

func (s *DebtService) GetDebts(ctx context.Context, userID string) ([]models.Debt, error) {
	filter := bson.M{"userID": userID}
	var userDebts []models.Debt

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &userDebts); err != nil {
		return nil, err
	}

	return userDebts, nil
}
