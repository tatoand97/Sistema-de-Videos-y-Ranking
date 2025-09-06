package repository

import (
	"api/internal/domain"
	"api/internal/domain/entities"
	"api/internal/domain/interfaces"
	"context"
	"errors"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).
		Select("FirstName", "LastName", "Email", "PasswordHash", "CityID").
		Create(user).Error
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetPermissions(ctx context.Context, userID uint) ([]string, error) {
	var perms []string
	err := r.db.WithContext(ctx).
		Table("privilege p").
		Select("DISTINCT p.name").
		Joins("JOIN role_privilege rp ON rp.privilege_id = p.privilege_id").
		Joins("JOIN user_role ur ON ur.role_id = rp.role_id").
		Where("ur.user_id = ?", userID).
		Scan(&perms).Error
	if err != nil {
		return nil, err
	}
	return perms, nil
}
