package main

import (
	"context"
	"log"
	"time"

	uuid "github.com/gofrs/uuid"
	"google.golang.org/api/iterator"
)

func getActivities(ctx context.Context) ([]Activity, error) {

	activities := client.Collection("activities")
	iter := activities.Documents(ctx)

	var m map[string]interface{}

	newActivities := make([]Activity, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Get activities ERR-1 - some error: %s", err)
			return nil, err
		}

		m = doc.Data()
		if m["type"] == "comment" {
			var c Comment

			err = doc.DataTo(&c)

			if err != nil {
				log.Printf("Get actitities - casting error: ERR-2 %s", err)
			} else {
				newActivities = append(newActivities, c)
			}
		}
	}

	return newActivities, nil

}

func getActivity(ctx context.Context, uid string) (*Activity, error) {
	query := client.Collection("activities").Where("meta.uuid", "==", uid)
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
	query := client.Collection("activities").Where("meta.uuid", "==", uid)
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

func storeComment(comment Comment) (*Comment, error) {
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
		return nil, err
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
