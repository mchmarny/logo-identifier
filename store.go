package main

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/spanner"
	ev "github.com/mchmarny/gcputil/env"
	"google.golang.org/grpc/codes"
)

var (
	dbClient *spanner.Client
	dbID     = ev.MustGetEnvVar("DB_ID", "")
)

// ServiceUser represents service input
type ServiceUser struct {
	UserID   string    `json:"id" spanner:"UserId"`
	Email    string    `json:"email" spanner:"Email"`
	UserName string    `json:"user" spanner:"UserName"`
	Created  time.Time `json:"created" spanner:"Created"`
	Updated  time.Time `json:"updated" spanner:"Updated"`
	Picture  string    `json:"pic" spanner:"Picture"`
}

// UserQuery represents service input
type UserQuery struct {
	UserID   string    `json:"userId" spanner:"UserId"`
	QueryID  string    `json:"queryId" spanner:"QueryId"`
	Created  time.Time `json:"created" spanner:"Created"`
	ImageURL string    `json:"image" spanner:"ImageUrl"`
	Result   string    `json:"result" spanner:"Result"`
}

func initStore(ctx context.Context) {
	c, err := spanner.NewClient(ctx, dbID)
	if err != nil {
		logger.Fatalf("Error while initializing db client: %v", err)
	}
	dbClient = c
}

func closeStore(ctx context.Context) {
	if dbClient != nil {
		dbClient.Close()
	}
}

func getUser(ctx context.Context, id string) (usr *ServiceUser, err error) {

	if id == "" {
		return nil, errors.New("Nil job ID parameter")
	}

	row, err := dbClient.Single().ReadRow(ctx, "Users", spanner.Key{id},
		[]string{"UserId", "Email", "UserName", "Created", "Updated", "Picture"})

	if err != nil {
		if spanner.ErrCode(err) == codes.NotFound {
			logger.Printf("User not found: %s", id)
			return nil, nil
		}

		logger.Printf("Error while quering for user %s: %v", id, err)
		return nil, err
	}

	var u ServiceUser
	if err := row.ToStruct(&u); err != nil {
		logger.Printf("Error while parsing DB user: %v", err)
		return nil, err
	}

	return &u, nil

}

func saveUser(ctx context.Context, usr *ServiceUser) error {

	if usr == nil || usr.UserID == "" {
		logger.Println("nil id on user save")
		return errors.New("User required")
	}

	m, err := spanner.InsertOrUpdateStruct("Users", usr)
	if err != nil {
		logger.Printf("Error while creating Users insert: %v", err)
		return err
	}

	return dbApply(ctx, m)

}

func saveQuery(ctx context.Context, q *UserQuery) error {

	if q == nil || q.QueryID == "" {
		logger.Println("nil id on query save")
		return errors.New("ID required")
	}

	m, err := spanner.InsertStruct("Queries", q)
	if err != nil {
		logger.Printf("Error while creating Users insert: %v", err)
		return err
	}

	return dbApply(ctx, m)

}

func dbApply(ctx context.Context, m *spanner.Mutation) error {

	_, err := dbClient.Apply(ctx, []*spanner.Mutation{m}, spanner.ApplyAtLeastOnce())
	if err != nil {
		logger.Printf("Error while applying user to DB: %v", err)
		return err
	}

	return nil

}
