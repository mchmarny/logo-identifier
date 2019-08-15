package main

import (
	"time"
)

// ServiceUser represents service input
type ServiceUser struct {
	UserID   string    `json:"id" spanner:"UserId"`
	Email    string    `json:"email" spanner:"Email"`
	UserName string    `json:"name" spanner:"UserName"`
	Created  time.Time `json:"created" spanner:"Created"`
	Updated  time.Time `json:"updated" spanner:"Updated"`
	Picture  string    `json:"pic" spanner:"Picture"`
}

// UserQuery represents service input
type UserQuery struct {
	UserID     string    `json:"userId" spanner:"UserId"`
	QueryID    string    `json:"queryId" spanner:"QueryId"`
	Created    time.Time `json:"created" spanner:"Created"`
	ImageURL   string    `json:"image" spanner:"ImageUrl"`
	Result     string    `json:"result" spanner:"Result"`
	QueryCount int64     `json:"queryCount" spanner:"-"`
	QueryLimit int64     `json:"queryLimit" spanner:"-"`
}
