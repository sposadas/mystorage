package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/sposadas/mystorage/internal/domain"
	"log"
)

type Repository interface {
	Store(ctx context.Context, user *domain.User) (*domain.User, error)
	GetOne(ctx context.Context, uuid uuid.UUID) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
	GetAll(ctx context.Context) ([]*domain.User, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

const (
	InsertUser = "INSERT INTO users(uuid, username, email, active) VALUES (?, ?, ?, ?)"
	GetOneUser = "SELECT u.id, u.uuid, u.username, u.email, u.active FROM users u WHERE u.uuid = ?"
	UpdateUser = "UPDATE users SET uuid = ?, username = ?, email = ?, active = ? WHERE id = ?"
	DeleteUser = "DELETE FROM users WHERE uuid = ?"
	GetUsers   = "SELECT u.id, u.uuid, u.username, u.email, u.active FROM users u"
)

func (r *repository) Store(ctx context.Context, user *domain.User) (*domain.User, error) {
	stmt, err := r.db.Prepare(InsertUser)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()
	var result sql.Result
	result, err = stmt.Exec((*user).UUID, (*user).Username, (*user).Email, (*user).Active)
	if err != nil {
		return nil, err
	}
	insertId, _ := result.LastInsertId()
	(*user).ID = int(insertId)

	return user, nil
}

func (r *repository) GetOne(ctx context.Context, UUID uuid.UUID) (*domain.User, error) {
	result := new(domain.User)

	row := r.db.QueryRowContext(ctx, GetOneUser, UUID)
	err := row.Err()

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	err = row.Scan(&result.ID, &result.UUID, &result.Username, &result.Email, &result.Active)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (r *repository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	result, err := r.db.ExecContext(ctx, UpdateUser, (*user).UUID, (*user).Username, (*user).Email, (*user).Active, (*user).ID)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return user, nil
}

func (r *repository) Delete(ctx context.Context, uuid uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, DeleteUser, uuid)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("%d users deleted", rows)
	}

	return nil
}

func (r *repository) GetAll(ctx context.Context) ([]*domain.User, error) {
	result := make([]*domain.User, 0)
	rows, err := r.db.QueryContext(ctx, GetUsers)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user := new(domain.User)
		err = rows.Scan(&user.ID, &user.UUID, &user.Username, &user.Email, &user.Active)
		if err != nil {
			return nil, err
		}
		result = append(result, user)
	}

	return result, nil
}
