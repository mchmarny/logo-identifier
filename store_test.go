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

func TestSession(t *testing.T) {

	ctx := context.Background()
	initStore(ctx)
	defer closeStore(ctx)

	uid := makeID("session-user@domain.com")
	sid := makeDailySessionID(uid)
	assert.NotNil(t, uid)
	assert.NotNil(t, sid)

	c1, err := countSession(ctx, uid, sid)
	assert.Nil(t, err)
	assert.True(t, c1 > -1, "Invalid session count: %d", c1)

	c2, err := countSession(ctx, uid, sid)
	assert.Nil(t, err)
	assert.True(t, c2 == c1+1, "Session count not incremented: %d", c1)

}

func TestUser(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping TestUser")
	}

	ctx := context.Background()
	initStore(ctx)
	defer closeStore(ctx)

	usr := getTestUserFromID("store-123")

	// reset
	err := deleteUser(ctx, usr.UserID)
	assert.Nil(t, err)

	// create
	err = saveUser(ctx, usr)
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
