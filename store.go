package main

import (
	"context"
	"errors"

	ev "github.com/mchmarny/gcputil/env"

	"database/sql"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type mysqlDB struct {
	conn            *sql.DB
	getUser         *sql.Stmt
	deleteUser      *sql.Stmt
	saveUser        *sql.Stmt
	saveQuery       *sql.Stmt
	addSession      *sql.Stmt
	getSessionCount *sql.Stmt
}

var (
	db  *mysqlDB
	dsn = ev.MustGetEnvVar("DSN", "")
)

func initStore(ctx context.Context) {

	c, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatalf("Error connecting to DB: %v", err)
	}

	if err := c.Ping(); err != nil {
		c.Close()
		logger.Fatalf("Error connecting to DB: %v", err)
	}

	db = &mysqlDB{
		conn: c,
	}

	if db.getUser, err = c.Prepare(`SELECT
		user_id, email, user_name, created, updated, pic_url
		FROM users WHERE user_id=?`); err != nil {
		logger.Fatalf("Error on selectUser prepare: %v", err)
	}

	if db.deleteUser, err = c.Prepare("DELETE FROM users WHERE user_id=?"); err != nil {
		logger.Fatalf("Error on deleteUser prepare: %v", err)
	}

	if db.saveUser, err = c.Prepare(`INSERT INTO users
		(user_id, email, user_name, created, updated, pic_url)
		VALUES (?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		user_name = ?, updated = ?, pic_url = ?
		`); err != nil {
		logger.Fatalf("Error on insertUser prepare: %v", err)
	}

	if db.saveQuery, err = c.Prepare(`INSERT INTO queries
		(query_id, user_id, created, img_url, result)
		VALUES (?, ?, ?, ?, ?)`); err != nil {
		logger.Fatalf("Error on insertQuery prepare: %v", err)
	}

	if db.addSession, err = c.Prepare(`INSERT INTO sessions
		(session_id, user_id, session_count)
		VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE
		session_count = session_count + 1`); err != nil {
		logger.Fatalf("Error on addSession prepare: %v", err)
	}

	if db.getSessionCount, err = c.Prepare(`SELECT session_count
		FROM sessions WHERE session_id = ?`); err != nil {
		logger.Fatalf("Error on addSession prepare: %v", err)
	}

}

func closeStore(ctx context.Context) {
	if db != nil && db.conn != nil {
		db.conn.Close()
	}
}

func getUser(ctx context.Context, id string) (usr *ServiceUser, err error) {

	if id == "" {
		return nil, errors.New("User ID parameter required")
	}

	rows, err := db.getUser.Query(id)
	if err != nil {
		logger.Printf("Error while quering for user %s: %v", id, err)
		return nil, err
	}
	defer rows.Close()

	var u *ServiceUser
	for rows.Next() {
		u = &ServiceUser{}
		if err := rows.Scan(&u.UserID, &u.Email, &u.UserName,
			&u.Created, &u.Updated, &u.Picture); err != nil {
			return nil, err
		}
	}
	return u, nil
}

func deleteUser(ctx context.Context, id string) error {

	if id == "" {
		logger.Println("nil id on user delete")
		return errors.New("User required")
	}

	_, err := db.deleteUser.Exec(id)
	if err != nil {
		logger.Printf("Error while deleting user %s: %v", id, err)
	}

	return err
}

func saveUser(ctx context.Context, usr *ServiceUser) error {

	if usr == nil || usr.UserID == "" {
		logger.Println("nil id on user save")
		return errors.New("User required")
	}

	_, err := db.saveUser.Exec(
		usr.UserID, usr.Email, usr.UserName, usr.Created, usr.Updated, usr.Picture,
		usr.UserName, usr.Updated, usr.Picture)
	if err != nil {
		logger.Printf("Error while saving user %v: %v", usr, err)
	}

	return err

}

func saveQuery(ctx context.Context, q *UserQuery) error {

	if q == nil || q.QueryID == "" {
		logger.Println("nil id on query save")
		return errors.New("ID required")
	}

	_, err := db.saveQuery.Exec(q.QueryID, q.UserID, q.Created, q.ImageURL, q.Result)
	if err != nil {
		logger.Printf("Error while saving query %v: %v", q, err)
	}

	return err

}

func countSession(ctx context.Context, userID, sessionID string) (c int64, err error) {

	if sessionID == "" || userID == "" {
		logger.Println("Either user or session ID required")
		return 0, errors.New("Both, user and session ID required")
	}

	tx, e := db.conn.Begin()
	if e != nil {
		logger.Printf("Error while creating transaction: %v", e)
	}

	_, e = tx.Stmt(db.addSession).Exec(sessionID, userID, 1)
	if e != nil {
		tx.Rollback()
		logger.Printf("Error while incrementing sessions %s: %v", sessionID, e)
	}

	rows, err := tx.Stmt(db.getSessionCount).Query(sessionID)
	if err != nil {
		tx.Rollback()
		logger.Printf("Error while quering session %s: %v", sessionID, err)
		return 0, err
	}
	defer rows.Close()

	var sessionCount int64
	for rows.Next() {
		if err := rows.Scan(&sessionCount); err != nil {
			tx.Rollback()
			logger.Printf("Error parsing session incrementing results: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logger.Printf("Error committing session incrementing: %v", err)
		return 0, err
	}

	logger.Printf("Session incrementing result: %d", sessionCount)

	return sessionCount, nil

}
