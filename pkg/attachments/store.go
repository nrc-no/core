package attachments

import (
	"context"
	"github.com/nrc-no/core/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AttachmentStore struct {
	getCollection func() (*mongo.Collection , error)
}

func NewAttachmentStore(ctx context.Context, mongoClientSrc storage.MongoClientSrc, database string, ) (*AttachmentStore, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	store := &AttachmentStore{
		getCollection: storage.GetCollectionFn(mongoClientSrc, database, "attachments"),
	}

	collection, err := store.getCollection()
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

func (s *AttachmentStore) Get(ctx context.Context, id string) (*Attachment, error) {
	collection, err := s.getCollection()
	if err != nil {
		return nil, err
	}

	res := collection.FindOne(ctx, bson.M{
		"id": id,
	})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var r Attachment
	if err := res.Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *AttachmentStore) List(ctx context.Context, options AttachmentListOptions) (*AttachmentList, error) {
	var filter = bson.M{}

	if len(options.AttachedToID) != 0 {
		filter = bson.M{
			"attachedToId": options.AttachedToID,
		}
	}

	collection, err := s.getCollection()
	if err != nil {
		return nil, err
	}

	res, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var items []*Attachment
	for {
		if !res.Next(ctx) {
			break
		}
		var r Attachment
		if err := res.Decode(&r); err != nil {
			return nil, err
		}
		items = append(items, &r)
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	if items == nil {
		items = []*Attachment{}
	}
	ret := AttachmentList{
		Items: items,
	}
	return &ret, nil
}

func (s *AttachmentStore) Update(ctx context.Context, attachment *Attachment) error {
	collection, err := s.getCollection()
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(ctx, bson.M{
		"id": attachment.ID,
	}, bson.M{
		"$set": bson.M{
			"attachedToId": attachment.AttachedToID,
			"body":         attachment.Body,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *AttachmentStore) Create(ctx context.Context, attachment *Attachment) error {
	collection, err := s.getCollection()
	if err != nil {
		return err
	}

	_, err = collection.InsertOne(ctx, attachment)
	if err != nil {
		return err
	}
	return nil
}
