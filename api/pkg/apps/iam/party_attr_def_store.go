package iam

import (
	"context"
	"github.com/nrc-no/core/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AttributeStore struct {
	getCollection utils.MongoCollectionFn
}

func newAttributeStore(ctx context.Context, mongoClientFn utils.MongoClientFn, database string) (*AttributeStore, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	store := &AttributeStore{
		getCollection: utils.GetCollectionFn(database, "attributes", mongoClientFn),
	}

	collection, err := store.getCollection(ctx)
	if err != nil {
		return nil, err
	}

	if _, err := collection.Indexes().CreateOne(ctx,
		mongo.IndexModel{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		}); err != nil {
		return nil, err
	}

	return store, nil
}

// list returns an PartyAttributeDefinitionList. If PartyAttributeDefinitionListOptions are supplied, list returns a filtered list containing
// only those items whose PartyAttributeDefinition.PartyTypeIDs field contains all the elements given in the query.
func (s *AttributeStore) list(ctx context.Context, listOptions PartyAttributeDefinitionListOptions) (*PartyAttributeDefinitionList, error) {

	filter := bson.M{}

	if len(listOptions.PartyTypeIDs) > 0 {
		filter["partyTypeIds"] = bson.M{
			"$all": listOptions.PartyTypeIDs,
		}
	}

	collection, err := s.getCollection(ctx)
	if err != nil {
		return nil, err
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var list []*PartyAttributeDefinition
	for {
		if !cursor.Next(ctx) {
			break
		}
		var a PartyAttributeDefinition
		if err := cursor.Decode(&a); err != nil {
			return nil, err
		}
		list = append(list, &a)
	}
	if cursor.Err() != nil {
		return nil, cursor.Err()
	}

	return &PartyAttributeDefinitionList{
		Items: list,
	}, nil

}

func (s *AttributeStore) create(ctx context.Context, attribute *PartyAttributeDefinition) error {
	collection, err := s.getCollection(ctx)
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, attribute)
	if err != nil {
		return err
	}
	return nil
}

func (s *AttributeStore) get(ctx context.Context, id string) (*PartyAttributeDefinition, error) {
	collection, err := s.getCollection(ctx)
	if err != nil {
		return nil, err
	}
	result := collection.FindOne(ctx, bson.M{
		"id": id,
	})
	if result.Err() != nil {
		return nil, result.Err()
	}
	var a PartyAttributeDefinition
	if err := result.Decode(&a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *AttributeStore) update(ctx context.Context, attribute *PartyAttributeDefinition) error {
	collection, err := s.getCollection(ctx)
	if err != nil {
		return err
	}
	_, err = collection.UpdateOne(ctx, bson.M{
		"id": attribute.ID,
	}, bson.M{
		"$set": bson.M{
			"countryId":    attribute.CountryID,
			"partyTypeIds": attribute.PartyTypeIDs,
			"isPii":        attribute.IsPersonallyIdentifiableInfo,
			"formControl":  attribute.FormControl,
		},
	})
	if err != nil {
		return err
	}
	return nil
}