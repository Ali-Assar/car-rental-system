package db

import (
	"context"
	"os"

	"github.com/Ali-Assar/car-rental-system/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CarStore interface {
	InsertCar(context.Context, *types.Car) (*types.Car, error)
	GetCars(context.Context, Map) ([]*types.Car, error)
}

type MongoCarStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	AgencyStore
}

func NewMongoCarStore(client *mongo.Client, agencyStore AgencyStore) *MongoCarStore {
	DbName := os.Getenv("MONGO_DB_NAME")
	return &MongoCarStore{
		client:      client,
		coll:        client.Database(DbName).Collection("cars"),
		AgencyStore: agencyStore,
	}
}

func (s *MongoCarStore) GetCars(ctx context.Context, filter Map) ([]*types.Car, error) {
	resp, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var cars []*types.Car
	if err := resp.All(ctx, &cars); err != nil {
		return nil, err
	}
	return cars, nil

}

func (s *MongoCarStore) InsertCar(ctx context.Context, car *types.Car) (*types.Car, error) {
	resp, err := s.coll.InsertOne(ctx, car)
	if err != nil {
		return nil, err
	}
	car.ID = resp.InsertedID.(primitive.ObjectID)

	//update the hotel with this room id
	filter := Map{"_id": car.AgencyID}
	update := Map{"$push": bson.M{"cars": car.ID}}
	if err := s.AgencyStore.UpdateAgency(ctx, filter, update); err != nil {
		return nil, err
	}
	return car, nil
}
