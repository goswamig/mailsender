package p

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/iterator"
	"time"
)

var subscribers *firestore.CollectionRef

func MaybeInit() error {
	ctx := context.Background()
	if subscribers == nil {
		firestoreClient, err := firestore.NewClient(ctx, "mailsender-288100")
		if err != nil {
			return err
		}

		subscribers = firestoreClient.Collection("subscribers")
	}
	return nil
}

func AlreadySubscribed(mail string) (bool, error) {
	ctx := context.Background()
	err := MaybeInit()
	if err != nil {
		return false, err
	}
	iter := subscribers.Where("email", "==", mail).Documents(ctx)
	_, err = iter.Next()
	if err == iterator.Done {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, err
}

func AddMailToDB(mail string) error {
	ctx := context.Background()
	err := MaybeInit()
	if err != nil {
		return err
	}
	_, _, err = subscribers.Add(ctx, map[string]interface{}{
		"email":    mail,
		"subscribedDate": time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func RemoveMailFromDB(mail string) error {
	ctx := context.Background()
	MaybeInit()
	subscribed, err := AlreadySubscribed(mail)
	if err != nil {
		return err
	}
	if !subscribed {
		return nil
	}
	iter := subscribers.Where("email", "==", mail).Documents(ctx)
	record, err := iter.Next()
	if err == iterator.Done {
		return nil
	}
	_, err = record.Ref.Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}