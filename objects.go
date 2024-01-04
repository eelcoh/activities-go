package main

import (
	"context"
	"log"
	"time"

	uuid "github.com/gofrs/uuid"
	"google.golang.org/api/iterator"
)

func getActivtities(ctx context.Context) ([]StoredActivity, error) {
	query := client.Collection("Activities")
	iter := query.Documents(ctx)

	var a StoredActivity
	var activities []StoredActivity
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		err = doc.DataTo(&a)
		if err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}
	return activities, nil

}

func getActivity(ctx context.Context, uid string) (*StoredActivity, error) {
	query := client.Collection("Activities").Where("meta.uuid", "==", uid)
	iter := query.Documents(ctx)

	var a StoredActivity
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		err = doc.DataTo(&a)
		if err != nil {
			return nil, err
		}
	}
	return &a, nil
}

func storeActivity(activityType string, activity Activity) (*StoredActivity, error) {
	var storedActivity *StoredActivity
	var act *StoredActivity

	uu, err := uuid.NewV4()

	if err != nil {
		return act, err
	}

	uid := uu.String()

	activity = activity.SetUUID(uid)
	// activityID := bson.NewObjectId()

	storedActivity.ActivityType = activityType
	storedActivity.Activity = activity
	storedActivity.ID = uid

	activities := client.Collection("activities")
	doc := activities.Doc(uid)

	_, err = doc.Create(ctx, storedActivity)

	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
		return act, err
	}

	return storedActivity, nil
}

func updateActivity(storedActivity *StoredActivity) (*StoredActivity, error) {

	var act *StoredActivity

	activities := client.Collection("activities")
	doc := activities.Doc(storedActivity.ID)

	_, err := doc.Set(ctx, storedActivity)
	if err != nil {
		return act, nil
	}
	return storedActivity, nil

}

func freshMeta() (MetaData, error) {
	uu, err := uuid.NewV4()
	if err != nil {
		return MetaData{}, err
	}

	now := time.Now()
	epoch := now.UnixNano() / 1000000

	meta := MetaData{
		Active: true,
		UUID:   uu.String(),
		Date:   epoch,
	}

	return meta, nil

}
