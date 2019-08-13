package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getTestUserFromID(id string) *ServiceUser {
	return &ServiceUser{
		ID:      id,
		Email:   fmt.Sprintf("id-%s@domain.com", id),
		Name:    "Test User",
		Created: time.Now(),
		Updated: time.Now(),
		Picture: "http://invalid.domain.com/pic1",
	}
}

func getTestEvent(userID, id string) *UserEvent {
	return &UserEvent{
		ID:     id,
		UserID: userID,
		On:     time.Now(),
		Image:  "d1",
		Result: "v1",
	}
}

func TestUser(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping TestSaveUser")
	}

	ctx := context.Background()
	initStore(ctx)

	// create
	usr := getTestUserFromID("store-123")
	err := saveUser(ctx, usr)
	assert.Nil(t, err)

	// get
	usr2, err := getUser(ctx, usr.ID)
	assert.Nil(t, err)
	assert.NotNil(t, usr2)
	assert.Equalf(t, usr.ID, usr2.ID, "Users' ID don't equal %s != %s", usr.ID, usr2.ID)

	// create events for user
	event1 := getTestEvent(usr2.ID, "e1")
	err = saveEvent(ctx, event1)
	assert.Nil(t, err)

	event2 := getTestEvent(usr2.ID, "e2")
	err = saveEvent(ctx, event2)
	assert.Nil(t, err)

	// delete user and its events
	err = deleteUser(ctx, usr2.ID)
	assert.Nil(t, err)

}
