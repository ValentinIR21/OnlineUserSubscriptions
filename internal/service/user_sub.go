package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"onlineusersub/internal/domain"
	"onlineusersub/internal/repository"
	"time"

	"github.com/google/uuid"
)

type SubService interface {
	GetSub(ctx context.Context, id string) (domain.Subscriptions, error)
	GetAllSub(ctx context.Context) ([]domain.Subscriptions, error)
	CreateSub(ctx context.Context, sub domain.Subscriptions) (domain.Subscriptions, error)
	UpdateSub(ctx context.Context, sub domain.Subscriptions) error
	DeleteSub(ctx context.Context, id string) error
	GetSumSub(ctx context.Context, userID, serviceName string, from, to time.Time) (int, error)
}

type subService struct {
	repos repository.SubRepository
}

// Создание нового экземплеря сервиса с переданным репозиторием
func NewSubService(repos repository.SubRepository) SubService {
	return &subService{
		repos: repos,
	}
}

// Ошибки сервисного слоя
var (
	ErrSubNotFound     = errors.New("Подписка не найдена")
	ErrSubSaveFailed   = errors.New("Ошибка сохранения в БД")
	ErrInvalidSub      = errors.New("Невалидный запрос")
	ErrSubUpdateFailed = errors.New("Ошибка обновления подписки")
	ErrSubSumFailed    = errors.New("Ошибка подсчета суммы")
	ErrDeleteSub       = errors.New("Ошибка удаления подписки")
)

// Возвращаение подписки из БД по ID
func (s *subService) GetSub(ctx context.Context, id string) (domain.Subscriptions, error) {

	sub, err := s.repos.GetOneByID(ctx, id)
	if err != nil {
		slog.Info("(service) Подписка не найдена", "id", id)
		return domain.Subscriptions{}, fmt.Errorf("%w, %s", ErrSubNotFound, id)
	}

	return sub, nil
}

// Возвращаение всех подписок из БД
func (s *subService) GetAllSub(ctx context.Context) ([]domain.Subscriptions, error) {

	subs, err := s.repos.ListBySubs(ctx)
	if err != nil {
		slog.Info("(service) Не удалось получить список подписок")
		return nil, fmt.Errorf("%w", ErrSubNotFound)
	}

	return subs, nil
}

// Подписка сохранена в БД
func (s *subService) CreateSub(ctx context.Context, sub domain.Subscriptions) (domain.Subscriptions, error) {

	if err := validateSub(sub); err != nil {
		return domain.Subscriptions{}, fmt.Errorf("%w, %v", ErrInvalidSub, err)
	}

	sub, err := s.repos.Create(ctx, sub)
	if err != nil {
		return domain.Subscriptions{}, fmt.Errorf("%w, %v", ErrSubSaveFailed, err)
	}

	slog.Info("(service) Подписка сохранена в БД", "id", sub.ID)

	return sub, nil
}

// Обновление подписки
func (s *subService) UpdateSub(ctx context.Context, sub domain.Subscriptions) error {

	if err := validateSub(sub); err != nil {
		return fmt.Errorf("%w, %v", ErrInvalidSub, err)
	}

	if err := s.repos.Update(ctx, sub); err != nil {
		return fmt.Errorf("%w, %v", ErrSubUpdateFailed, err)
	}

	slog.Info("(service) Подписка обновлена", "id", sub.ID)

	return nil
}

// Удаление подписки из БД
func (s *subService) DeleteSub(ctx context.Context, id string) error {

	if err := s.repos.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w, %v", ErrDeleteSub, err)
	}

	slog.Info("(service) Подписка удалена", "id", id)

	return nil
}

// Подпсчет суммы подписок
func (s *subService) GetSumSub(ctx context.Context, userID, serviceName string, from, to time.Time) (int, error) {

	sum, err := s.repos.GetTotalSum(ctx, userID, serviceName, from, to)
	if err != nil {
		return 0, fmt.Errorf("%w, %v", ErrSubSumFailed, err)
	}

	slog.Info("(service) Сумма подписок", "id", userID, "sum", sum)

	return sum, nil
}

// Валидация
func validateSub(sub domain.Subscriptions) error {

	if sub.UserID == uuid.Nil {
		slog.Error("UserID незаполнен")
		return errors.New("user_id обязателен")
	}

	if sub.ServiceName == "" {
		slog.Error("ServiceName незаполнен")
		return errors.New("service_name обязателен")
	}

	if sub.Price <= 0 {
		slog.Error("Price меньше нуля")
		return errors.New("price меньше нуля")
	}

	if sub.DateCreated.IsZero() {
		slog.Error("DateCreated незаполнен")
		return errors.New("start_date обязателен")
	}

	return nil
}
