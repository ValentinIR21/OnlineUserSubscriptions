package repository

import (
	"context"
	"onlineusersub/internal/domain"
)

type SubRepository interface {
	GetOneByID(ctx context.Context, uid string) (domain.Subscriptions, error)
	Update(ctx context.Context, sub domain.Subscriptions) error
	Delete(ctx context.Context, id string) error
	ListByUserID(ctx context.Context, uid string) ([]domain.Subscriptions, error)
	GetTotalSum(ctx context.Context, userID, serviceName, from, to string) (int, error)
}

type Repository interface {
	SubRepository
	Close()
}
