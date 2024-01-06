package main

import (
	"context"
	"log"
	"time"

	uuid "github.com/gofrs/uuid"
	"google.golang.org/api/iterator"
)

func getActivtities(ctx context.Context) ([]Activity, error) {
	log.Println("Get activtities - 1")

	activities := client.Collection("Activities")
	log.Println("Get activtities - 2")
	iter := activities.Documents(ctx)
	log.Println("Get activtities - 3")

	var a Activity
	log.Println("Get activtities - 4")
	var newActivities []Activity
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

func getActivity(ctx context.Context, uid string) (*Activity, error) {
	query := client.Collection("Activities").Where("meta.uuid", "==", uid)
	iter := query.Documents(ctx)

	var a Activity
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

func getBlog(ctx context.Context, uid string) (*Blog, error) {
	query := client.Collection("Activities").Where("meta.uuid", "==", uid)
	iter := query.Documents(ctx)

	var a Blog
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

// func storeActivity(activityType string, activity Activity) (*StoredActivity, error) {
// 	log.Println("Store activtity - 1")

// 	var storedActivity *StoredActivity
// 	log.Println("Store activtity - 2")
// 	var act *StoredActivity
// 	log.Println("Store activtity - 3")

// 	activity = activity.SetUUID(activity.GetUUID())
// 	// activityID := bson.NewObjectId()

// 	storedActivity.ActivityType = activityType
// 	storedActivity.Activity = activity
// 	storedActivity.ID = uid

// 	log.Println("Store activtity - 4")
// 	activities := client.Collection("activities")
// 	log.Println("Store activtity - 5")
// 	doc := activities.Doc(uid)
// 	log.Println("Store activtity - 6")
// 	log.Printf("activity: %s", storedActivity)
// 	log.Println("Store activtity - 7")

// 	log.Printf("doc: %s", doc)
// 	log.Println("Store activtity - 8")
// 	_, err = doc.Create(ctx, storedActivity)
// 	log.Println("Store activtity - 9")

// 	if err != nil {
// 		log.Fatalf("Failed adding document: %v", err)
// 		return act, err
// 	}

// 	return storedActivity, nil
// }

func storeActivity(activityType string, activity Activity) (*Activity, error) {
	activities := client.Collection("activities")

	doc := activities.Doc(activity.GetUUID())

	_, err := doc.Create(ctx, activity)

	if err != nil {
		log.Fatalf("Failed adding document: %v", err)
		return nil, err
	}

	return &activity, nil
}

func storeComment(activityType string, comment Comment) (*Comment, error) {
	activities := client.Collection("activities")

	doc := activities.Doc(comment.GetUUID())

	_, err := doc.Create(ctx, comment)

	if err != nil {
		log.Fatalf("Failed adding document: %v", err)
		return nil, err
	}

	return &comment, nil
}

func storeBlog(activityType string, blog Blog) (*Blog, error) {
	var act *Blog
	activities := client.Collection("activities")

	doc := activities.Doc(blog.GetUUID())

	_, err := doc.Create(ctx, blog)

	if err != nil {
		log.Fatalf("Failed adding document: %v", err)
		return act, err
	}

	return &blog, nil
}

func storeNewRanking(activityType string, ranking NewRanking) (*NewRanking, error) {
	activities := client.Collection("activities")

	doc := activities.Doc(ranking.GetUUID())

	_, err := doc.Create(ctx, ranking)

	if err != nil {
		log.Fatalf("Failed adding document: %v", err)
		return nil, err
	}

	return &ranking, nil
}

func storeNewBet(activityType string, bet NewBet) (*NewBet, error) {
	activities := client.Collection("activities")

	doc := activities.Doc(bet.GetUUID())

	_, err := doc.Create(ctx, bet)

	if err != nil {
		log.Fatalf("Failed adding document: %v", err)
		return nil, err
	}

	return &bet, nil
}

func updateActivity(activity Activity) (*Activity, error) {

	activities := client.Collection("activities")
	doc := activities.Doc(activity.GetUUID())

	_, err := doc.Set(ctx, activity)
	if err != nil {
		return nil, nil
	}
	return &activity, nil

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
