package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/timsolov/ms-users/app/domain/entity"
)

// DB describes
type DB struct {
	db *sqlx.DB
}

func New(ctx context.Context, dsn string, maxConns, maxIdle int, connLifeTime time.Duration) (*DB, error) {
	db, err := sqlx.ConnectContext(ctx, "pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxConns)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(connLifeTime)
	// db.SetConnMaxIdleTime(connIdleTime)

	for i := 0; i < 5; i++ {
		if err := db.Ping(); err == nil {
			return &DB{db: db}, nil
		} else {
			fmt.Println("can't connect to DB retry after 2 seconds")
			time.Sleep(2 * time.Second)
		}
	}

	return nil, fmt.Errorf("can't connect to DB: %s", dsn)
}

// CreateUser creates new user record
func (d *DB) CreateUser(ctx context.Context, m *entity.User) error {
	if m.UserID == uuid.Nil {
		m.UserID = uuid.New()
	}
	if m.CreatedAt.IsZero() || m.UpdatedAt.IsZero() {
		now := time.Now()
		m.CreatedAt = now
		m.UpdatedAt = now
	}

	err := d.execr(ctx, 1,
		`INSERT 
			INTO "users" (user_id, email, first_name, last_name, created_at, updated_at) 
			VALUES (?,?,?,?,?,?)`,
		m.UserID, m.Email, m.FirstName, m.LastName, m.CreatedAt, m.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Profile returns user record
func (d *DB) Profile(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	var user entity.User
	return user, E(d.get(ctx, &user, "SELECT user_id, email, first_name, last_name, created_at, updated_at FROM users WHERE user_id = ?", userID))
}
