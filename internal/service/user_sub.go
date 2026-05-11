package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"onlineusersub/internal/domain"
	"onlineusersub/internal/repository"
	"time"

	"github.com/google/uuid"
)

type SubService interface {
	GetSub(ctx context.Context, id string) (domain.Subscriptions, error)
	GetAllSub(ctx context.Context, uid string) ([]domain.Subscriptions, error)
	CreateSub(ctx context.Context, sub domain.Subscriptions) error
	UpdateSub(ctx context.Context, sub domain.Subscriptions) error
	DeleteSub(ctx context.Context, id string) error
	GetSumSub(ctx context.Context, userID, serviceName string, from, to time.Time) (int, error)
}

type subService struct {
	repos repository.SubRepository
}

func NewSubService(repos repository.SubRepository) SubService {
	return &subService{
		repos: repos,
	}
}

var (
	ErrSubNotFound     = errors.New("Подписка не найдена")
	ErrSubSaveFailed   = errors.New("Ошибка сохранения в БД")
	ErrInvalidSub      = errors.New("Невалидный запрос")
	ErrSubUpdateFailed = errors.New("Ошибка обновления подписки")
	ErrSubSumFailed    = errors.New("Ошибка подсчета суммы")
	ErrDeleteSub       = errors.New("Ошибка удаления подписки")
)

// Возвращаение подписки из БД
func (s *subService) GetSub(ctx context.Context, id string) (domain.Subscriptions, error) {

	order, err := s.repos.GetOneByID(ctx, id)
	if err != nil {
		log.Printf("(service) Подписка не найдена: %s", id)
		return domain.Subscriptions{}, fmt.Errorf("%w, %s", ErrSubNotFound, id)
	}

	return order, nil
}

// Возвращаение всех подписок из БД
func (s *subService) GetAllSub(ctx context.Context, uid string) ([]domain.Subscriptions, error) {

	order, err := s.repos.ListByUserID(ctx, uid)
	if err != nil {
		log.Printf("(service) Не удалось получить список подписок id: %s", uid)
		return nil, fmt.Errorf("%w, %s", ErrSubNotFound, uid)
	}

	return order, nil
}

// Подписка сохранена в БД
func (s *subService) CreateSub(ctx context.Context, sub domain.Subscriptions) error {

	if err := validateSub(sub); err != nil {
		return fmt.Errorf("%w, %v", ErrInvalidSub, err)
	}

	if err := s.repos.Create(ctx, sub); err != nil {
		return fmt.Errorf("%w, %v", ErrSubSaveFailed, err)
	}

	log.Printf("(service) Подписка %s сохранена в БД", sub.ID)

	return nil
}

// Обновление подписки
func (s *subService) UpdateSub(ctx context.Context, sub domain.Subscriptions) error {

	if err := validateSub(sub); err != nil {
		return fmt.Errorf("%w, %v", ErrInvalidSub, err)
	}

	if err := s.repos.Update(ctx, sub); err != nil {
		return fmt.Errorf("%w, %v", ErrSubUpdateFailed, err)
	}

	log.Printf("(service) Подписка %s обновлена", sub.ID)

	return nil
}

// Удаление подписки из БД
func (s *subService) DeleteSub(ctx context.Context, id string) error {

	if err := s.repos.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w, %v", ErrDeleteSub, err)
	}

	log.Printf("(service) Подписка %s удалена", id)

	return nil
}

// Подпсчет суммы подписок
func (s *subService) GetSumSub(ctx context.Context, userID, serviceName string, from, to time.Time) (int, error) {

	sum, err := s.repos.GetTotalSum(ctx, userID, serviceName, from, to)
	if err != nil {
		return 0, fmt.Errorf("%w, %v", ErrSubSumFailed, err)
	}

	log.Printf("(service) Сумма подписок %s = %d", userID, sum)

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
