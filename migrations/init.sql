CREATE TABLE IF NOT EXISTS subscriptions(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    date_created DATE NOT NULL,
    date_conclusion DATE
);

-- Индес на ускорения запроса суммы подписок
CREATE INDEX IF NOT EXISTS idx_subscriptions_sum_price
    ON subscriptions (user_id, service_name, date_created, price);