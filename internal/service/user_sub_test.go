package service

// Запуск тестов: go test ./internal/service/...
// С флагом -v для подробного вывода: go test -v ./internal/service/...

import (
	"context"
	"errors"
	"onlineusersub/internal/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Тесты для CreateSub
func TestCreateSub(t *testing.T) {
	validSub := domain.Subscriptions{
		UserID:      uuid.New(),
		ServiceName: "Yandex Plus",
		Price:       400,
		DateCreated: time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
	}

	// Таблица тест-кейсов
	tests := []struct {
		name      string               // название теста, выводится при -v
		sub       domain.Subscriptions // входные данные
		repoErr   error                // что вернёт репозиторий
		wantErr   bool                 // ожидаем ли ошибку вообще
		targetErr error                // конкретная ошибка, которую ожидаем
	}{
		{
			name:    "успешное создание",
			sub:     validSub,
			repoErr: nil,
			wantErr: false,
		},
		{
			name: "пустой UserID",
			sub: domain.Subscriptions{
				UserID:      uuid.Nil,
				ServiceName: "Netflix",
				Price:       500,
				DateCreated: time.Now(),
			},
			wantErr:   true,
			targetErr: ErrInvalidSub,
		},
		{
			name: "пустой ServiceName",
			sub: domain.Subscriptions{
				UserID:      uuid.New(),
				ServiceName: "",
				Price:       500,
				DateCreated: time.Now(),
			},
			wantErr:   true,
			targetErr: ErrInvalidSub,
		},
		{
			name: "нулевая цена",
			sub: domain.Subscriptions{
				UserID:      uuid.New(),
				ServiceName: "Spotify",
				Price:       0,
				DateCreated: time.Now(),
			},
			wantErr:   true,
			targetErr: ErrInvalidSub,
		},
		{
			name: "отрицательная цена",
			sub: domain.Subscriptions{
				UserID:      uuid.New(),
				ServiceName: "Spotify",
				Price:       -100,
				DateCreated: time.Now(),
			},
			wantErr:   true,
			targetErr: ErrInvalidSub,
		},
		{
			name: "нулевая дата",
			sub: domain.Subscriptions{
				UserID:      uuid.New(),
				ServiceName: "Spotify",
				Price:       300,
			},
			wantErr:   true,
			targetErr: ErrInvalidSub,
		},
		{
			name:      "ошибка репозитория",
			sub:       validSub,
			repoErr:   errors.New("connection refused"),
			wantErr:   true,
			targetErr: ErrSubSaveFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// mock-репозиторий.
			repo := &mockSubRepository{
				createFn: func(ctx context.Context, sub domain.Subscriptions) (domain.Subscriptions, error) {
					return sub, tt.repoErr
				},
			}

			svc := NewSubService(repo)
			got, err := svc.CreateSub(context.Background(), tt.sub)

			if tt.wantErr {

				require.Error(t, err)

				if tt.targetErr != nil {

					assert.True(t, errors.Is(err, tt.targetErr),
						"ожидали ошибку %v, получили %v", tt.targetErr, err)

					assert.Equal(t, domain.Subscriptions{}, got)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Тесты для GetSub
func TestGetSub(t *testing.T) {
	existingID := uuid.New()
	existingSub := domain.Subscriptions{
		ID:          existingID,
		UserID:      uuid.New(),
		ServiceName: "Yandex Plus",
		Price:       400,
		DateCreated: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	tests := []struct {
		name      string
		id        string
		repoSub   domain.Subscriptions
		repoErr   error
		wantErr   bool
		targetErr error
	}{
		{
			name:    "подписка найдена",
			id:      existingID.String(),
			repoSub: existingSub,
			repoErr: nil,
			wantErr: false,
		},
		{
			name:      "подписка не найдена",
			id:        uuid.New().String(),
			repoErr:   errors.New("no rows"),
			wantErr:   true,
			targetErr: ErrSubNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockSubRepository{
				getOneByIDFn: func(ctx context.Context, id string) (domain.Subscriptions, error) {
					return tt.repoSub, tt.repoErr
				},
			}

			svc := NewSubService(repo)
			got, err := svc.GetSub(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.targetErr))
				assert.Equal(t, domain.Subscriptions{}, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.repoSub, got)
			}
		})
	}
}

// Тесты для DeleteSub
func TestDeleteSub(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		repoErr   error
		wantErr   bool
		targetErr error
	}{
		{
			name:    "успешное удаление",
			id:      uuid.New().String(),
			repoErr: nil,
			wantErr: false,
		},
		{
			name:      "ошибка репозитория при удалении",
			id:        uuid.New().String(),
			repoErr:   errors.New("db error"),
			wantErr:   true,
			targetErr: ErrDeleteSub,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockSubRepository{
				deleteFn: func(ctx context.Context, id string) error {
					assert.Equal(t, tt.id, id)
					return tt.repoErr
				},
			}

			svc := NewSubService(repo)
			err := svc.DeleteSub(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.targetErr))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Тесты для GetSumSub
func TestGetSumSub(t *testing.T) {
	userID := uuid.New().String()
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		userID      string
		serviceName string
		repoSum     int
		repoErr     error
		wantSum     int
		wantErr     bool
	}{
		{
			name:        "сумма посчитана",
			userID:      userID,
			serviceName: "Yandex Plus",
			repoSum:     1200,
			wantSum:     1200,
			wantErr:     false,
		},
		{
			name:        "нет подписок — сумма 0",
			userID:      userID,
			serviceName: "Netflix",
			repoSum:     0,
			wantSum:     0,
			wantErr:     false,
		},
		{
			name:        "ошибка репозитория",
			userID:      userID,
			serviceName: "Spotify",
			repoErr:     errors.New("db timeout"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockSubRepository{
				getTotalSumFn: func(ctx context.Context, uid, svc string, f, t time.Time) (int, error) {
					return tt.repoSum, tt.repoErr
				},
			}

			svc := NewSubService(repo)
			sum, err := svc.GetSumSub(context.Background(), tt.userID, tt.serviceName, from, to)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantSum, sum)
			}
		})
	}
}

// Тесты для validateSub
func TestValidateSub(t *testing.T) {
	baseValid := domain.Subscriptions{
		UserID:      uuid.New(),
		ServiceName: "Test",
		Price:       100,
		DateCreated: time.Now(),
	}

	tests := []struct {
		name    string
		modify  func(*domain.Subscriptions) // что изменить в базовой валидной структуре
		wantErr bool
	}{
		{
			name:    "все поля валидны",
			modify:  func(s *domain.Subscriptions) {}, // ничего не меняем
			wantErr: false,
		},
		{
			name:    "UserID nil",
			modify:  func(s *domain.Subscriptions) { s.UserID = uuid.Nil },
			wantErr: true,
		},
		{
			name:    "ServiceName пустой",
			modify:  func(s *domain.Subscriptions) { s.ServiceName = "" },
			wantErr: true,
		},
		{
			name:    "Price = 0",
			modify:  func(s *domain.Subscriptions) { s.Price = 0 },
			wantErr: true,
		},
		{
			name:    "Price отрицательный",
			modify:  func(s *domain.Subscriptions) { s.Price = -1 },
			wantErr: true,
		},
		{
			name:    "DateCreated zero",
			modify:  func(s *domain.Subscriptions) { s.DateCreated = time.Time{} },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Копируем базовую структуру и применяем модификацию
			sub := baseValid
			tt.modify(&sub)

			err := validateSub(sub)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
