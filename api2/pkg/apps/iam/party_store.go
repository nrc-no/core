package iam

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PartyStore struct {
	Collection *mongo.Collection
}

func NewPartyStore(mongoClient *mongo.Client, database string) *PartyStore {
	return &PartyStore{
		Collection: mongoClient.Database(database).Collection("parties"),
	}
}

func (s *PartyStore) Get(ctx context.Context, id string) (*Party, error) {
	res := s.Collection.FindOne(ctx, bson.M{
		"id": id,
	})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var r Party
	if err := res.Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *PartyStore) List(ctx context.Context, listOptions PartyListOptions) (*PartyList, error) {
	filterItems := bson.A{}

	if len(listOptions.PartyTypeID) != 0 {
		filterItems = append(filterItems, bson.M{"partyTypeIds": listOptions.PartyTypeID})
	}

	if len(listOptions.SearchParam) != 0 {
		filterItems = append(filterItems, bson.M{
			"$text": bson.M{
				"$search": listOptions.SearchParam,
			},
		})
	}

	var filter interface{}
	if len(filterItems) == 0 {
		filter = bson.M{}
	} else if len(filterItems) == 1 {
		filter = filterItems[0]
	} else {
		filter = bson.M{"$and": filterItems}
	}

	res, err := s.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var items []*Party
	for {
		if !res.Next(ctx) {
			break
		}
		var r Party
		if err := res.Decode(&r); err != nil {
			return nil, err
		}
		items = append(items, &r)
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	if items == nil {
		items = []*Party{}
	}
	ret := PartyList{
		Items: items,
	}
	return &ret, nil
}

func (s *PartyStore) Update(ctx context.Context, party *Party) error {
	_, err := s.Collection.UpdateOne(ctx, bson.M{
		"id": party.ID,
	}, bson.M{
		"$set": bson.M{
			"attributes":   party.Attributes,
			"partyTypeIds": party.PartyTypeIDs,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *PartyStore) Create(ctx context.Context, party *Party) error {
	_, err := s.Collection.InsertOne(ctx, party)
	if err != nil {
		return err
	}
	return nil
}

type FindOptions struct {
	Attributes map[string]string
}

func (s *PartyStore) Find(ctx context.Context, options FindOptions) (*Party, error) {
	filter := bson.M{}
	for key, value := range options.Attributes {
		filter["attributes."+key] = value
	}
	res := s.Collection.FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}
	var party *Party
	if err := res.Decode(&party); err != nil {
		return nil, err
	}
	return party, nil
}
