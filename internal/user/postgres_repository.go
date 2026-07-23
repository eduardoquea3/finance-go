package user

import (
	"context"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) FindByEmail(ctx context.Context, email string) (Entity, error) {
	var entity Entity
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return Entity{}, ErrNotFound
		}
		return Entity{}, err
	}
	return entity, nil
}
