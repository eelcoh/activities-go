package main

import (
	"context"
	"log"
	"time"

	uuid "github.com/gofrs/uuid"
	"google.golang.org/api/iterator"
)

func getActivtities(ctx context.Context) ([]StoredActivity, error) {
	log.Println("Get activtities - 1")

	activities := client.Collection("Activities")
	log.Println("Get activtities - 2")
	iter := activities.Documents(ctx)
	log.Println("Get activtities - 3")

	var a StoredActivity
	log.Println("Get activtities - 4")
	var newActivities []StoredActivity
	for {
		doc, err := iter.Next()
		log.Println("Get activtities - 5")
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		err = doc.DataTo(&a)
		log.Println("Get activtities - 6")

		if err != nil {
			return nil, err
		}
		newActivities = append(newActivities, a)
	}
	return newActivities, nil

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
	log.Println("Store activtity - 1")

	var storedActivity *StoredActivity
	log.Println("Store activtity - 2")
	var act *StoredActivity
	log.Println("Store activtity - 3")

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

	log.Println("Store activtity - 4")
	activities := client.Collection("activities")
	log.Println("Store activtity - 5")
	doc := activities.Doc(uid)
	log.Println("Store activtity - 6")
	log.Printf("activity: %s", storedActivity)
	log.Println("Store activtity - 7")

	log.Printf("doc: %s", doc)
	log.Println("Store activtity - 8")
	_, err = doc.Create(ctx, storedActivity)
	log.Println("Store activtity - 9")

	if err != nil {
		log.Fatalf("Failed adding document: %v", err)
		return act, err
	}

	return storedActivity, nil
}

func storeComment(activityType string, comment Comment) (*Comment, error) {
	log.Println("Store comment - 1")

	log.Println("Store comment - 2")
	var act *Comment
	log.Println("Store comment - 3")

	log.Println("Store comment - 4")
	activities := client.Collection("activities")

	log.Println("Store comment - 5")
	doc := activities.Doc(comment.GetUUID())

	log.Println("Store comment - 6")
	log.Printf("comment: %s", comment)
	log.Println("Store comment - 7")

	log.Printf("doc: %s", doc)
	log.Println("Store comment - 8")

	type State struct {
		Capital    string  `firestore:"capital"`
		Population float64 `firestore:"pop"` // in millions
	}
	// _, err := doc.Create(ctx, comment)
	wr, err := doc.Create(ctx, State{
		Capital:    "Albany",
		Population: 19.8,
	})
	log.Println("Store comment - 9")
	log.Printf("comment: %s", wr)

	if err != nil {
		log.Fatalf("Failed adding document: %v", err)
		return act, err
	}

	return &comment, nil
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
