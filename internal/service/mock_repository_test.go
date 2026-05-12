package service

import (
	"context"
	"onlineusersub/internal/domain"
	"time"
)

type mockSubRepository struct {
	createFn      func(ctx context.Context, sub domain.Subscriptions) (domain.Subscriptions, error)
	getOneByIDFn  func(ctx context.Context, id string) (domain.Subscriptions, error)
	updateFn      func(ctx context.Context, sub domain.Subscriptions) error
	deleteFn      func(ctx context.Context, id string) error
	listBySubsFn  func(ctx context.Context) ([]domain.Subscriptions, error)
	getTotalSumFn func(ctx context.Context, userID, serviceName string, from, to time.Time) (int, error)
}

// Реализуем интерфейс SubRepository
func (m *mockSubRepository) Create(ctx context.Context, sub domain.Subscriptions) (domain.Subscriptions, error) {
	return m.createFn(ctx, sub)
}

func (m *mockSubRepository) GetOneByID(ctx context.Context, id string) (domain.Subscriptions, error) {
	return m.getOneByIDFn(ctx, id)
}

func (m *mockSubRepository) Update(ctx context.Context, sub domain.Subscriptions) error {
	return m.updateFn(ctx, sub)
}

func (m *mockSubRepository) Delete(ctx context.Context, id string) error {
	return m.deleteFn(ctx, id)
}

func (m *mockSubRepository) ListBySubs(ctx context.Context) ([]domain.Subscriptions, error) {
	return m.listBySubsFn(ctx)
}

func (m *mockSubRepository) GetTotalSum(ctx context.Context, userID, serviceName string, from, to time.Time) (int, error) {
	return m.getTotalSumFn(ctx, userID, serviceName, from, to)
}
