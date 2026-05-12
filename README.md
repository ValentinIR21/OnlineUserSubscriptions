# Subscriptions Service

Cервис для агрегации данных об онлайн-подписках пользователей

## Запуск
 
```bash
# Клонировать репозиторий
git clone https://github.com/ValentinIR21/OnlineUserSubscriptions.git
 
# Запустить все 
docker-compose up --build
```
Сервер запустится на `http://localhost:8081`


### Остановка
 
```bash
docker compose down
```

Для удаления данных (volume с БД):
 
```bash
docker compose down -v
```


## Swagger UI
 
После запуска документация доступна по адресу:
 
```
http://localhost:8081/swagger/index.html
```

Там можно посмотреть все эндпоинты и отправить запросы прямо из браузера.


## API
 
Формат дат: `MM-YYYY` (например `07-2025`).
 
### Создать подписку
 
```
POST /subscription/publish
```
 
Тело запроса:
 
```json
{
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "service_name": "Yandex Plus",
  "price": 400,
  "start_date": "07-2025",
  "conclusion_date": "12-2025"
}
```
 
`conclusion_date` — опциональное поле.
 
---
 
### Получить подписку по ID
 
```
GET /subscription/{id}
```
 
---
 
### Обновить подписку
 
```
PATCH /subscription/{id}
```
 
Тело запроса — аналогично созданию.
 
---
 
### Удалить подписку
 
```
DELETE /subscription/{id}
```
 
---
 
### Список всех подписок
 
```
GET /subscriptions
```
 
---
 
### Сумма подписок за период
 
```
GET /subscriptions/sum?user_id=...&service=...&from=MM-YYYY&to=MM-YYYY
```


## Тесты
 
Тесты покрывают сервисный слой: бизнес-логику и валидацию. В тестах реальная БД не используется — вместо неё подставляется mock-репозиторий.
 
```bash
# Запустить тесты
go test ./internal/service/...
 
# С подробным выводом
go test -v ./internal/service/...
```
