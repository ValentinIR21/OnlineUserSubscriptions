package repository

import (
	"context"
	"fmt"
	"log"
	"onlineusersub/internal/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRep(ctx context.Context, connString string) (*PostgresRepository, error) {

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("Не удалось создать пул соединений: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("БД не пингуется: %w", err)
	}

	return &PostgresRepository{pool: pool}, nil
}

func (p *PostgresRepository) Close() {
	p.pool.Close()
}

func (p *PostgresRepository) Create(ctx context.Context, subscriptions domain.Subscriptions) error {

	const query = `
		INSERT INTO subscriptions (
			user_id, service_name, price, date_created, date_conclusion)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := p.pool.Exec(ctx, query,
		subscriptions.UserID,
		subscriptions.ServiceName,
		subscriptions.Price,
		subscriptions.DateCreated,
		subscriptions.DateConclusion,
	)
	if err != nil {
		return fmt.Errorf("Ошибка сохранения в БД %w", err)
	}

	return nil
}

func (p *PostgresRepository) GetOneByID(ctx context.Context, id string) (domain.Subscriptions, error) {

	const query = `
		SELECT 
			id, user_id, service_name, price, date_created, date_conclusion
		FROM subscriptions
		WHERE id = $1
	`
	var subscription domain.Subscriptions

	if err := p.pool.QueryRow(ctx, query, id).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.ServiceName,
		&subscription.Price,
		&subscription.DateCreated,
		&subscription.DateConclusion,
	); err != nil {
		return domain.Subscriptions{}, fmt.Errorf("Подписка с id %s не найдена: %w", id, err)
	}

	return subscription, nil
}

func (p *PostgresRepository) Update(ctx context.Context, sub domain.Subscriptions) error {

	const query = `
		UPDATE subscriptions
		SET service_name = $1,
			price = $2,
			date_created = $3, 
			date_conclusion = $4
		WHERE id = $5
	`

	result, err := p.pool.Exec(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.DateCreated,
		sub.DateConclusion,
		sub.ID,
	)
	if err != nil {
		return fmt.Errorf("Ошибка обновления записи: %w", err)
	}

	// Проверяем, была ли такая запись в БД
	if result.RowsAffected() == 0 {
		return fmt.Errorf("Запись с id %s не найдена", sub.ID)
	}

	return nil
}

func (p *PostgresRepository) Delete(ctx context.Context, id string) error {

	const query = `DELETE from subscriptions WHERE id = $1`

	_, err := p.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("Ошибка удаления: %w", err)
	}

	return nil
}

func (p *PostgresRepository) ListBySubs(ctx context.Context) ([]domain.Subscriptions, error) {

	const query = `
		SELECT 
			id, user_id, service_name, price, date_created, date_conclusion
		FROM subscriptions
		ORDER BY date_created DESC
	`

	rows, err := p.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Ошибка запроса: %w", err)
	}
	defer rows.Close()

	var subscriptions []domain.Subscriptions

	for rows.Next() {
		var subscription domain.Subscriptions

		if err := rows.Scan(
			&subscription.ID, &subscription.UserID, &subscription.ServiceName, &subscription.Price,
			&subscription.DateCreated, &subscription.DateConclusion,
		); err != nil {
			log.Printf("Пропуск строки при сканировании: %v", err)
			continue
		}

		subscriptions = append(subscriptions, subscription)

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Ошибка итерации строк: %w", err)
	}

	return subscriptions, nil
}

func (p *PostgresRepository) GetTotalSum(ctx context.Context, userID, serviceName string, from, to time.Time) (int, error) {

	const query = `
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions
		WHERE user_id = $1
			AND service_name = $2
			AND date_created >= $3
			AND date_created <= $4
	`

	var total int

	err := p.pool.QueryRow(ctx, query, userID, serviceName, from, to).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("Ошибка подсчета суммы: %w", err)
	}

	return total, nil
}
