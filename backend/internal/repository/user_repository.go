package repository

import (
	"context"

	"github.com/seymourrisey/staredesk/internal/entity"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}
