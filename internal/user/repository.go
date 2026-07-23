package user

import "context"

type Repository interface {
	FindByEmail(ctx context.Context, email string) (Entity, error)
}
