package main

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	uuid "github.com/gofrs/uuid"
	"google.golang.org/api/iterator"
)

func getActivtities(ctx context.Context) ([]Activity, error) {

	activities := client.Collection("activities").OrderBy("meta.date", firestore.Desc)
	iter := activities.Documents(ctx)

	var m map[string]interface{}

	newActivities := make([]Activity, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Get activtities ERR-1 %s", err)
			return nil, err
		}

		m = doc.Data()

		switch m["type"] {
		case "comment":
			var a Comment

			err = doc.DataTo(&a)

			if err != nil {
				log.Printf("Error parsing COMMENT %s : %s", m, err)
			} else {
				newActivities = append(newActivities, a)
			}
		case "blog":
			var a Blog

			err = doc.DataTo(&a)

			if err != nil {
				log.Printf("Error parsing BLOG %s", err)
			} else {
				newActivities = append(newActivities, a)
			}
		case "new bet":
			var a NewRanking

			err = doc.DataTo(&a)

			if err != nil {
				log.Printf("Error parsing NEW BET %s", err)
			} else {
				newActivities = append(newActivities, a)
			}
		case "new ranking":
			var a NewRanking

			err = doc.DataTo(&a)

			if err != nil {
				log.Printf("Error parsing NEW RANKING %s", err)
			} else {
				newActivities = append(newActivities, a)
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

func storeActivity(activity Activity) (*Activity, error) {
	activities := client.Collection("activities")

	doc := activities.Doc(activity.GetUUID())

	_, err := doc.Create(ctx, activity)

	if err != nil {
		log.Fatalf("Failed adding document: %v", err)
		return nil, err
	}

	return &activity, nil
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
