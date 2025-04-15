# Medods Auth API

Сервис аутентификации на Go с использованием JWT (SHA512) и PostgreSQL.
Реализует выдачу и обновление пары Access/Refresh токенов с защитой от повторного использования и контролем IP-адреса.

---

## Технологии

- Go
- PostgreSQL
- JWT (SHA512)
- bcrypt
- sqlx
- chi
- zap (логирование)
- docker-compose
- golang-migrate
- testify (unit-тесты)

---

## Инструкция для запуска

### 1. Клонируйте репозиторий

```bash
git clone https://github.com/Egorpalan/medods-api.git
cd medods-api
```


### 2. Настройте переменные окружения

Создайте `.env` на основе `.env.example`.

Проверьте настройки docker-compose.yml и настройте корректные порты в случае необходимости.

```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=auth_db
JWT_SECRET=supersecretkey
ACCESS_TOKEN_TTL=15m
REFRESH_TOKEN_TTL=168h
APP_PORT=8080
```


### 3. Запустите сервис и базу данных

```bash
make compose-up
```


### 4. Примените миграции

В новом терминале:

```bash
make migrate-up
```


### 5. Запустите тесты

```bash
make test
```

---

## API

### 1. Получить пару токенов

**POST** `/token?user_id={GUID};`

- **Параметры:**
  `user_id` — идентификатор пользователя (GUID)
- **Ответ:**

```json
{
  "accessToken": "&lt;JWT&gt;",
  "refreshToken": "&lt;refresh_token&gt;"
}
```


### 2. Обновить пару токенов

**POST** `/refresh`

- **Body (JSON):**

```json
{
  "access_token": "&lt;accessToken&gt;",
  "refresh_token": "&lt;refreshToken&gt;"
}
```

- **Заголовок:**
  `X-Forwarded-For` — IP клиента (опционально, для теста смены IP.)
- Как протестировать смену IP
  * Получите новую пару токенов через /token.
  * Сразу (первый раз) отправьте запрос на /refresh с другим IP (например, через заголовок X-Forwarded-For).
  * В логах появится предупреждение о смене IP и будет выдана новая пара токенов.
  * Если повторить refresh с тем же токеном — будет ошибка "refresh token already used"
- **Ответ:**

```json
{
  "accessToken": "&lt;новый JWT&gt;",
  "refreshToken": "&lt;новый refresh&gt;"
}
```

---

## Особенности

- **Access токен** — JWT (SHA512), не хранится в базе.
- **Refresh токен** — base64, хранится только в виде bcrypt-хеша, защищён от повторного использования.
- **Связь токенов** — refresh можно использовать только с тем access, с которым он был выдан.
- **Контроль IP** — при смене IP при refresh отправляется warning (мокируется логом).
- **Логирование** — zap, все события и ошибки логируются.
- **Graceful shutdown** — сервис корректно завершает работу по SIGINT/SIGTERM.
- **Тесты** — покрытие бизнес-логики unit-тестами (testify).

---

## Примеры запросов

### Получить токены

```bash
curl -X POST "http://localhost:8080/token?user_id=123e4567-e89b-12d3-a456-426614174000"
```


### Обновить токены

```bash
curl -X POST "http://localhost:8080/refresh" \
  -H "Content-Type: application/json" \
  -H "X-Forwarded-For: 129.0.0.2" \
  -d '{"access_token":"&lt;accessToken&gt;","refresh_token":"&lt;refreshToken&gt;"}'
```

---

## Makefile команды

- `make run` — локальный запуск
- `make build` — сборка бинарника
- `make compose-up` — запуск docker-compose
- `make compose-down` — остановка docker-compose
- `make migrate-up` — применить миграции
- `make migrate-down` — откатить миграции
- `make test` — запустить unit-тесты

