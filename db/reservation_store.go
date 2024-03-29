package db

import (
	"context"
	"os"

	"github.com/Ali-Assar/car-rental-system/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReservationStore interface {
	InsertReservation(context.Context, *types.Reservation) (*types.Reservation, error)
	GetReservation(context.Context, bson.M) ([]*types.Reservation, error)
	GetReservationByID(context.Context, string) (*types.Reservation, error)
	UpdateReservation(context.Context, string, bson.M) error
}

type MongoReservationStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	ReservationStore
}

func NewMongoReservationStore(client *mongo.Client) *MongoReservationStore {
	DbName := os.Getenv("MONGO_DB_NAME")
	return &MongoReservationStore{
		client: client,
		coll:   client.Database(DbName).Collection("reservation"),
	}
}

func (s *MongoReservationStore) UpdateReservation(ctx context.Context, id string, update bson.M) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	m := bson.M{"$set": update}
	_, err = s.coll.UpdateByID(ctx, oid, m)
	return err
}

func (s *MongoReservationStore) GetReservation(ctx context.Context, filter bson.M) ([]*types.Reservation, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var reservation []*types.Reservation
	if err := cur.All(ctx, &reservation); err != nil {
		return nil, err
	}
	return reservation, nil
}

func (s *MongoReservationStore) InsertReservation(ctx context.Context, reservation *types.Reservation) (*types.Reservation, error) {
	resp, err := s.coll.InsertOne(ctx, reservation)
	if err != nil {
		return nil, err
	}
	reservation.ID = resp.InsertedID.(primitive.ObjectID)
	return reservation, nil
}

func (s *MongoReservationStore) GetReservationByID(ctx context.Context, id string) (*types.Reservation, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var reservation types.Reservation
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&reservation); err != nil {
		return nil, err
	}
	return &reservation, nil
}
