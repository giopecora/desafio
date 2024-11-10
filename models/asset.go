package models

type Asset struct {
	ID     string  `bson:"_id,omitempty"`
	UserID string  `bson:"userID"`
	Valor  float64 `bson:"valor"`
	Tipo   string  `bson:"tipo"`
}
