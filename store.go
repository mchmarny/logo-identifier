package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

const (
	userCollectionName  = "logo-identifier-user"
	eventCollectionName = "logo-identifier-event"
)

var (
	fsClient     *firestore.Client
	userColl     *firestore.CollectionRef
	eventColl    *firestore.CollectionRef
	errNilDocRef = errors.New("firestore: nil DocumentRef")
)

// ServiceUser represents service input
type ServiceUser struct {
	ID      string    `json:"id" firestore:"id"`
	Email   string    `json:"email" firestore:"email"`
	Name    string    `json:"name" firestore:"name"`
	Created time.Time `json:"created" firestore:"created"`
	Updated time.Time `json:"updated" firestore:"updated"`
	Picture string    `json:"pic" firestore:"pic"`
}

// UserEvent represents service input
type UserEvent struct {
	ID     string    `json:"id" firestore:"id"`
	UserID string    `json:"user" firestore:"user"`
	On     time.Time `json:"ts" firestore:"ts"`
	Image  string    `json:"image" firestore:"image"`
	Result string    `json:"result" firestore:"result"`
}

func initStore(ctx context.Context) {

	// in case called multiple times during test
	if eventColl != nil && userColl != nil && fsClient != nil {
		return
	}

	// Assumes GOOGLE_APPLICATION_CREDENTIALS is set
	c, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		logger.Fatalf("Error while creating Firestore client: %v", err)
	}
	fsClient = c
	userColl = c.Collection(userCollectionName)
	eventColl = c.Collection(eventCollectionName)
}

func getUser(ctx context.Context, id string) (usr *ServiceUser, err error) {

	if id == "" {
		return nil, errors.New("Nil job ID parameter")
	}

	d, err := userColl.Doc(id).Get(ctx)
	if err != nil {
		if err == errNilDocRef {
			return nil, fmt.Errorf("No user for ID: %s", id)
		}
		return nil, err
	}

	var u ServiceUser
	if err := d.DataTo(&u); err != nil {
		return nil, fmt.Errorf("Stored data not user: %v", err)
	}

	return &u, nil

}

func saveUser(ctx context.Context, usr *ServiceUser) error {

	if usr == nil || usr.ID == "" {
		logger.Println("nil id on user save")
		return errors.New("Nil ID")
	}

	_, err := userColl.Doc(usr.ID).Set(ctx, usr)
	if err != nil {
		logger.Printf("error on save: %v", err)
		return fmt.Errorf("Error on save: %v", err)
	}

	return nil

}

func saveEvent(ctx context.Context, event *UserEvent) error {

	if event == nil || event.ID == "" {
		logger.Println("nil id on event save")
		return errors.New("Nil ID")
	}

	_, err := eventColl.Doc(event.ID).Set(ctx, event)
	if err != nil {
		logger.Printf("error on save: %v", err)
		return fmt.Errorf("Error on save: %v", err)
	}

	return nil

}

func deleteUser(ctx context.Context, id string) error {

	if id == "" {
		return errors.New("Nil job ID parameter")
	}

	batch := fsClient.Batch()
	doc, err := userColl.Doc(id).Get(ctx)
	if err != nil {
		logger.Printf("Error on doc get: %v", err)
		return err
	}
	batch.Delete(doc.Ref)

	iter := eventColl.Where("userId", "==", id).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Printf("Error on iter next: %v", err)
			return err
		}
		batch.Delete(doc.Ref)
	}

	_, err = batch.Commit(ctx)
	return err

}
