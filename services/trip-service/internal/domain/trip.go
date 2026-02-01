package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripModel struct {
	ID            primitive.ObjectID
	UserID        string // Customer ID
	Status        string
	RideFareModel *RideFareModel
}
