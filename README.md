# Запуск БД в контейнере

Для запуска БД в контейнере выполните:

```bash
make pg
```

Для остановки:

```bash
make stop-pg
```

Для удаления данных из БД:

```bash
make clean-data
```

# Сценарий

Создадим первую миграцию:

```bash
docker run --rm \
    -v $(realpath ./db/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        create \
        -dir /migrations \
        -ext .sql \
        -seq -digits 5 \
        init
```

Опишем первую версию БД. Добавим в `00001_init.up.sql` следующий код:

```sql
CREATE TABLE positions(
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title VARCHAR(200) UNIQUE NOT NULL
);

CREATE TABLE employees(
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    first_name VARCHAR(200) NOT NULL,
    last_name VARCHAR(200) NOT NULL,
    salary NUMERIC NOT NULL,
    position INT NOT NULL REFERENCES positions(id),
    CONSTRAINT employees_salary_positive_check CHECK (salary::numeric > 0)
);
```

Попробуем применить миграцию к БД. Контейнеры должны взаимодействовать друг с другом, для этого нужно узнать адрес контейнера с БД в сети docker'а:

```bash
docker inspect praktikum-webinar-db | grep IPAddress
```

После этого выполним:

```bash
docker run --rm \
    -v $(realpath ./db/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database postgres://gopher:gopher@172.17.0.2:5432/gopher_corp?sslmode=disable \
        up
```

Подключимся к БД:

```bash
psql -h localhost -p 5432 -U gopher -d gopher_corp
```

и посмотрим на результат:

```sql
\d employees; \d positions
```

Также обратим внимание на таблицу `schema_migrations`:

```sql
SELECT *
FROM schema_migrations;
```

Эта таблица автоматически созданная `go-migrate`, которая содержит информацию о текущем состоянии БД.

Попробуем вернуться к начальному состоянию:

```bash
docker run --rm \
    -v $(realpath ./db/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database postgres://gopher:gopher@172.17.0.2:5432/gopher_corp?sslmode=disable \
        down -all
```

Что наблюдаем?

Попробуем заново применить миграцию:

```bash
docker run --rm \
    -v $(realpath ./db/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database postgres://gopher:gopher@172.17.0.2:5432/gopher_corp?sslmode=disable \
        up
```

Что наблюдаем?

Исправим это. Установим версию в `1`:

```bash
docker run --rm \
    -v $(realpath ./db/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database postgres://gopher:gopher@172.17.0.2:5432/gopher_corp?sslmode=disable \
        force 1
```

Проверим сосотояние:

```sql
SELECT *
FROM schema_migrations;
```

Добавим код отменяющий миграцию в `00001_init.up.sql`:

```sql
DROP TABLE employees;
DROP TABLE positions;
```

"Откатим" миграцию:

```bash
docker run --rm \
    -v $(realpath ./db/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database postgres://gopher:gopher@172.17.0.2:5432/gopher_corp?sslmode=disable \
        down
```

Что наблюдаем?

Попробуем снова:

```bash
docker run --rm \
    -v $(realpath ./db/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database postgres://gopher:gopher@172.17.0.2:5432/gopher_corp?sslmode=disable \
        down -all
```

Миграции лучше оборачивать в транзакции, в случае ошибки легче откатиться на предыдущую версию. Добавим транзакции в файлы миграций:

```sql
BEGIN TRANSACTION;
-- остальной код
COMMIT;
```

Представим, что было решено добавить колонку `email` в таблицу `employees`. Давайте создадим для этого следующую миграцию:

```bash
docker run --rm \
    -v $(realpath ./db/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        create \
        -dir /migrations \
        -ext .sql \
        -seq -digits 5 \
        add_email_to_employees
```

И добавим следующий код в `00002_add_email_to_employees.up.sql`:

```sql
BEGIN TRANSACTION;

ALTER TABLE employees
ADD COLUMN email VARCHAR(200) NOT NULL;

COMMIT;
```

Добавим также описание комманд для отмены изменений в `00002_add_email_to_employees.down.sql`:

```sql
BEGIN TRANSACTION;

ALTER TABLE employees
DROP COLUMN email;

COMMIT;
```

Применим миграцию:

```bash
docker run --rm \
    -v $(realpath ./db/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database postgres://gopher:gopher@172.17.0.2:5432/gopher_corp?sslmode=disable \
        up
```

Добавим данные в БД:

```sql
INSERT INTO positions(title)
VALUES ('Go developer');

INSERT INTO employees (first_name, last_name, salary, position, email)
VALUES ('Alice', 'Liddell', '100000', (SELECT id FROM positions WHERE title='Go developer'), 'alice.liddell@gopher-corp.com');
```

Убедимся, что данные добавлены:

```sql
SELECT *
FROM employees;
```

Представим теперь, что мы захотели откатиться на одну версию назад. Что случится с таблицей после выполнения следующей команды?

```bash
docker run --rm \
    -v $(realpath ./db/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database postgres://gopher:gopher@172.17.0.2:5432/gopher_corp?sslmode=disable \
        down 1
```

## Как избежать потери данных

Зависит от подхода, принятого в комманде.

Один из вариантов - не удалять даные, а помечать как удаленные.

Применим миграцию `2` из `db/migrations-careful`:

```bash
docker run --rm \
    -v $(realpath ./db/migrations-careful):/migrations \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database postgres://gopher:gopher@172.17.0.2:5432/gopher_corp?sslmode=disable \
        up
```

Запишем значение электронной почты сотрудника:

```sql
UPDATE employees
SET email='alice.liddell@gopher-corp.com'
WHERE first_name='Alice' AND last_name='Liddell';
```

Теперь отменим последнюю миграцию:

```bash
docker run --rm \
    -v $(realpath ./db/migrations-careful):/migrations \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database postgres://gopher:gopher@172.17.0.2:5432/gopher_corp?sslmode=disable \
        down 1
```

Посмотрим на таблицу, а затем опять применим миграцию:

```bash
docker run --rm \
    -v $(realpath ./db/migrations-careful):/migrations \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database postgres://gopher:gopher@172.17.0.2:5432/gopher_corp?sslmode=disable \
        up
```

# Применение миграций из кода

Миграции можно автоматически применять при старте приложения. См. файл `app/internal/store/db.go`.

Откатим миграции:

```bash
docker run --rm \
    -v $(realpath ./db/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database postgres://gopher:gopher@172.17.0.2:5432/gopher_corp?sslmode=disable \
        down -all
```

Запустите приложение, для этого выполните эту комманду из папки `app`:

```bash
make && ./cmd/migrations/migrations -dsn "postgresql://gopher:gopher@localhost:5432/gopher_corp?sslmode=disable"
```

Создадим еще раз должность:

```sql
INSERT INTO positions(title)
VALUES ('Go developer');
```

Попробуем создать нового сотрудника:

```bash
curl --request POST http://127.0.0.1:8080/employee \
    --header "Content-Type: application/json" \
    --data '{"first_name": "Bob", "last_name": "Morane", "salary": 75000, "position": "Go developer", "email": "bob.morane@gopher-corp.com"}'
```

Попробуем добавить еще одного сотрудника (с отрицательной зарплатой):

```bash
curl --request POST http://127.0.0.1:8080/employee \
    --header "Content-Type: application/json" \
    --data '{"first_name": "Charlie", "last_name": "Bucket", "salary": -75000, "position": "Go developer", "email": "charlie.bucket@gopher-corp.com"}' \
    -v
```

Мы видим, что произошла ошибка. Это ожидаемо, т.к. БД проверяет положительность зарплаты. Было бы здорово тестировать это, чтобы убедиться, что этот функционал не будет утерян в результате регрессии.

# Интеграционное тестирование

См. `app/internal/store/db_integration_test.go`
