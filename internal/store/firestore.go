package store

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/eelcoh/activities/internal/activity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const defaultCollection = "activities"

type Firestore struct {
	client     *firestore.Client
	collection string
}

func NewFirestore(client *firestore.Client) *Firestore {
	return &Firestore{
		client:     client,
		collection: defaultCollection,
	}
}

func (s *Firestore) List(ctx context.Context) ([]activity.Activity, error) {
	docs, err := s.client.Collection(s.collection).
		OrderBy("meta.date", firestore.Desc).
		Documents(ctx).
		GetAll()
	if err != nil {
		return nil, err
	}

	items := make([]activity.Activity, 0, len(docs))
	for _, doc := range docs {
		item, err := decodeActivity(doc)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (s *Firestore) Get(ctx context.Context, id string) (activity.Activity, error) {
	doc, err := s.client.Collection(s.collection).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, activity.ErrNotFound
		}
		return nil, err
	}

	return decodeActivity(doc)
}

func (s *Firestore) Create(ctx context.Context, item activity.Activity) error {
	_, err := s.client.Collection(s.collection).Doc(item.GetMeta().UUID).Create(ctx, item)
	return err
}

func (s *Firestore) Update(ctx context.Context, item activity.Activity) error {
	_, err := s.client.Collection(s.collection).Doc(item.GetMeta().UUID).Set(ctx, item)
	return err
}

func decodeActivity(doc *firestore.DocumentSnapshot) (activity.Activity, error) {
	rawType, ok := doc.Data()["type"].(string)
	if !ok {
		return nil, errors.New("activity document missing type")
	}

	switch rawType {
	case activity.TypeComment:
		var item activity.Comment
		if err := doc.DataTo(&item); err != nil {
			return nil, err
		}
		return item, nil
	case activity.TypeBlog:
		var item activity.Blog
		if err := doc.DataTo(&item); err != nil {
			return nil, err
		}
		return item, nil
	case activity.TypeNewBet:
		var item activity.NewBet
		if err := doc.DataTo(&item); err != nil {
			return nil, err
		}
		return item, nil
	case activity.TypeNewRanking:
		var item activity.NewRanking
		if err := doc.DataTo(&item); err != nil {
			return nil, err
		}
		return item, nil
	default:
		return nil, fmt.Errorf("unsupported activity type %q", rawType)
	}
}
