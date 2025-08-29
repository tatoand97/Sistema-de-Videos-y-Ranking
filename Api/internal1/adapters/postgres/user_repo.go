package postgres

import (
	"database/sql"
	"errors"

	"main_viderk/internal/domain"
	"main_viderk/internal/ports"
)

type userRepo struct{ db *sql.DB }

func NewUserRepo(db *sql.DB) ports.UserRepository { return &userRepo{db: db} }

func (r *userRepo) Create(username, email, passwordHash string, profileImagePath *string) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(`insert into users(username, email, password_hash, profile_image_path)
                        values ($1,$2,$3,$4) returning id, username, email, profile_image_path`,
		username, email, passwordHash, profileImagePath).
		Scan(&u.ID, &u.Username, &u.Email, &u.ProfileImagePath)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetByUsername(username string) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(`select id, username, password_hash, profile_image_path from users where username=$1`, username).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.ProfileImagePath)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &u, err
}

func (r *userRepo) GetByID(id int64) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(`select id, username, password_hash, profile_image_path from users where id=$1`, id).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.ProfileImagePath)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &u, err
}

func (r *userRepo) GetByEmail(email string) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(`select id, username, email, password_hash, profile_image_path from users where email=$1`, email).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.ProfileImagePath)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &u, err
}
