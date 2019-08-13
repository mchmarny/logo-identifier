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
		UserID:   id,
		Email:    fmt.Sprintf("id-%s@domain.com", id),
		UserName: "Test User",
		Created:  time.Now(),
		Updated:  time.Now(),
		Picture:  "http://invalid.domain.com/pic1",
	}
}

func getTestEvent(userID, id string) *UserQuery {
	return &UserQuery{
		QueryID:  id,
		UserID:   userID,
		Created:  time.Now(),
		ImageURL: "d1",
		Result:   "v1",
	}
}

func TestUser(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping TestSaveUser")
	}

	ctx := context.Background()
	initStore(ctx)
	defer closeStore(ctx)

	// create
	usr := getTestUserFromID("store-123")
	err := saveUser(ctx, usr)
	assert.Nil(t, err)

	// get
	usr2, err := getUser(ctx, usr.UserID)
	assert.Nil(t, err)
	assert.NotNil(t, usr2)
	assert.Equalf(t, usr.UserID, usr2.UserID, "Users' ID don't equal %s != %s", usr.UserID, usr2.UserID)

	// create events for user
	e1 := getTestEvent(usr2.UserID, "e1")
	err = saveQuery(ctx, e1)
	assert.Nil(t, err)

	e2 := getTestEvent(usr2.UserID, "e2")
	err = saveQuery(ctx, e2)
	assert.Nil(t, err)

}
