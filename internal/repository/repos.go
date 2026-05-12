package repository

import (
	"context"
	"onlineusersub/internal/domain"
	"time"
)

type SubRepository interface {
	Create(ctx context.Context, subscriptions domain.Subscriptions) (domain.Subscriptions, error)
	GetOneByID(ctx context.Context, uid string) (domain.Subscriptions, error)
	Update(ctx context.Context, sub domain.Subscriptions) error
	Delete(ctx context.Context, id string) error
	ListBySubs(ctx context.Context) ([]domain.Subscriptions, error)
	GetTotalSum(ctx context.Context, userID, serviceName string, from, to time.Time) (int, error)
}

type Repository interface {
	SubRepository
	Close()
}
