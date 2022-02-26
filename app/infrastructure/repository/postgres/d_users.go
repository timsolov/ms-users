package postgres

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/timsolov/ms-users/app/domain/entity"
	"github.com/timsolov/ms-users/app/infrastructure/repository/postgres/models"
)

// NewUser creates new user record
func (d *DB) NewUser(ctx context.Context, e *entity.User) error {
	var m models.User
	m.FromEntity(e)

	if m.UserID == uuid.Nil {
		m.UserID = uuid.New()
	}

	db := d.db.WithContext(context.TODO())

	if err := db.Create(m).Error; err != nil {
		return E(err)
	}

	return nil
}

// User returns user record by id
func (d *DB) User(ctx context.Context, userID uuid.UUID, columns ...string) (entity.User, error) {
	var m models.User

	db := d.db.WithContext(context.TODO())
	if len(columns) > 0 {
		m.UserID = userID
		db = db.Select(strings.Join(columns, ","))
	}

	if err := db.Take(&m, "user_id = ?", userID).Error; err != nil {
		return m.ToEntity(), E(err)
	}

	return m.ToEntity(), nil
}

// UpdUser changes user record's properties.
func (d *DB) UpdUser(ctx context.Context, e *entity.User, columns ...string) error {
	var m models.User
	m.FromEntity(e)

	db := d.db.WithContext(context.TODO())

	if len(columns) > 0 {
		return E(db.Model(m).Where("user_id = ?", m.UserID).Select(columns).Updates(m).Error)
	}
	return E(db.Model(m).Where("user_id = ?", m.UserID).Updates(m).Error)
}

// DelUser deletes user record record.
func (d *DB) DelUser(ctx context.Context, userID uuid.UUID) error {
	db := d.db.WithContext(context.TODO())

	return E(db.Delete(&entity.User{}, "user_id = ?", userID).Error)
}
