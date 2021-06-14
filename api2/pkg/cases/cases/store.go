package cases

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/cases/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	collection *mongo.Collection
}

func NewStore(mongoClient *mongo.Client) *Store {
	return &Store{
		collection: mongoClient.Database("core").Collection("cases"),
	}
}

func (s *Store) Get(ctx context.Context, id string) (*api.Case, error) {
	res := s.collection.FindOne(ctx, bson.M{
		"id": id,
	})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var r api.Case
	if err := res.Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *Store) List(ctx context.Context, listOptions ListOptions) (*api.CaseList, error) {

	filter := bson.M{}

	if len(listOptions.Case) != 0 {
		filter["$or"] = bson.M{
			"firstCase":  listOptions.Case,
			"secondCase": listOptions.Case,
		}
	}

	res, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var items []*api.Case
	for {
		if !res.Next(ctx) {
			break
		}
		var r api.Case
		if err := res.Decode(&r); err != nil {
			return nil, err
		}
		items = append(items, &r)
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	if items == nil {
		items = []*api.Case{}
	}
	ret := api.CaseList{
		Items: items,
	}
	return &ret, nil
}

func (s *Store) Update(ctx context.Context, kase *api.Case) error {
	_, err := s.collection.UpdateOne(ctx, bson.M{
		"id": kase.ID,
	}, bson.M{
		"$set": bson.M{
			"caseType":    kase.CaseTypeID,
			"partyId":     kase.PartyID,
			"description": kase.Description,
			"done":        kase.Done,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Create(ctx context.Context, kase *api.Case) error {
	_, err := s.collection.InsertOne(ctx, kase)
	if err != nil {
		return err
	}
	return nil
}
