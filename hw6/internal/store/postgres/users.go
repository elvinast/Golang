package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"Go/hw6/internal/models"
	"Go/hw6/internal/store"
)

func (db *DB) Users() store.UserRepository {
	if db.profile == nil {
		db.profile = NewProfileRep(db.conn)
	}
	return db.profile
}

type ProfileRepository struct {
	conn *sqlx.DB
}

func NewProfileRep(conn *sqlx.DB) store.UserRepository {
	return &ProfileRepository{conn: conn}
}

func (c ProfileRepository) Create(ctx context.Context, user *models.User) error {
	_, err := c.conn.Exec("INSERT INTO profile(title) VALUES ($1)", user.Email)
	if err != nil {
		return err
	}
	return nil
}

func (c ProfileRepository) All(ctx context.Context, filter *models.UserFilter) ([]*models.User, error) {
	basicQuery := "SELECT * FROM profile"
	if filter.Query != nil {
		basicQuery += " WHERE email ilike '%" + *filter.Query + "%'"
	}
	users := make([]*models.User, 0)
	if err := c.conn.Select(&users, basicQuery); err != nil {
		return nil, err
	}
	return users, nil
}

func (c ProfileRepository) ByID(ctx context.Context, id int) (*models.User, error) {
	user := new(models.User)
	if err := c.conn.Get(user, "SELECT id, email FROM profile WHERE id=$1", id); err != nil {
		return nil, err
	}
	return user, nil
}

func (c ProfileRepository) Update(ctx context.Context, user *models.User) error {
	_, err := c.conn.Exec("UPDATE profile SET email = $1 WHERE id = $2", user.Email, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (c ProfileRepository) Delete(ctx context.Context, id int) error {
	_, err := c.conn.Exec("DELETE FROM profile WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}