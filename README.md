# 🏠 Mortgage Calculator Service

Cервис для расчета параметров ипотеки с кэшированием

## 📖 Обзор

RESTful API сервис на Go для расчета ключевых параметров ипотеки:
- **Процентная ставка** (в зависимости от программы кредитования)
- **Сумма кредита**
- **Аннуитетный ежемесячный платеж**
- **Общая переплата**
- **Дата последнего платежа**

Все расчеты сохраняются в in-memory кэше и могут быть получены по запросу

## 🚀 Возможности

- ✅ REST API эндпоинты
- ✅ In-memory кэширование расчетов
- ✅ Структурированное логирование через middleware
- ✅ Комплексная валидация входных данных
- ✅ Докеризация
- ✅ Юнит-тесты
- ✅ Проверка качества кода golangci-lint
- ✅ Конфигурация через YAML файл

## 📋 API Эндпоинты

🟢 POST /execute
Расчет параметров ипотеки.

**Запрос:**
```bash
curl -X POST http://localhost:8080/execute \
  -H "Content-Type: application/json" \
  -d '{
    "object_cost": 5000000,
    "initial_payment": 1000000,
    "months": 240,
    "program": {
      "salary": true
    }
  }'
```

**Ответ:**
```
{
  "result": {
    "params": {
      "object_cost": 5000000,
      "initial_payment": 1000000,
      "months": 240
    },
    "program": {
      "salary": true
    },
    "aggregates": {
      "rate": 8,
      "loan_sum": 4000000,
      "monthly_payment": 33458,
      "overpayment": 4029920,
      "last_payment_date": "2044-02-18"
    }
  }
}
```

🟢 GET /cache
Получить все рассчитанные ипотеки из кэша.

**Запрос:**
```bash
curl -X GET http://localhost:8080/cache
```

**Ответ:**
```
[
  {
    "id": 0,
    "params": {
      "object_cost": 5000000,
      "initial_payment": 1000000,
      "months": 240
    },
    "program": {
      "salary": true
    },
    "aggregates": {
      "rate": 8,
      "loan_sum": 4000000,
      "monthly_payment": 33458,
      "overpayment": 4029920,
      "last_payment_date": "2044-02-18"
    }
  }
]
```

## 🔮 HTTP statuses

Сервис возвращает соответствующие HTTP статусы:
* 400 Bad Request - Неверные входные данные, неправильный выбор программы или недостаточный первоначальный взнос
* 200 OK - Успешный расчет
* 500 Internal Server Error - Проблемы на стороне сервера

## 🏗️ Программы кредитования

Программа	            Ставка	Мин. первоначальный взнос
🏢 Корпоративная	    8%	    20%
🪖 Военная ипотека	    9%	    20%
🏠 Базовая программа	10%	    20%

## ⚙️ Конфигурация

Сервис настраивается через config.yml:
```
port: 8080
```

## 🐳 Docker

Pull docker image:
```
docker pull seraleu/mortgage-calculator:latest
```
Сборка образа:
```
docker build -t seraleu/mortgage-calculator:latest .
```
Запуск контейнера: 
```
docker run -d -p 8080:8080 --name mortgage-calc seraleu/mortgage-calculator:latest
```
Пуш в Docker Hub:
```
docker push seraleu/mortgage-calculator:latest
```

## 🛠️ Разработка

Требования
1. Go 1.25+
2. Docker

Локальная разработка:

1. git clone https://github.com/serAleu/mortgage-calculator.git
2. cd mortgage-calculator
3. go mod download
4. go test ./... -v
5. golangci-lint run
6. go run ./cmd/server

Makefile:
```
make build        # Собрать бинарник
make test         # Запустить тесты
make lint         # Запустить линтер
make docker-build # Собрать Docker образ
make run          # Запустить Docker контейнер
make stop         # Остановить Docker контейнер
```
 
## 📦 Структура проекта
```
mortgage-calculator/
├── cmd/
│   └── server/
│       └── main.go          # Точка входа
├── internal/
│   ├── app/                 # Инициализация приложения
│   ├── config/              # Конфигурация
│   ├── controller/          # HTTP обработчики
│   ├── calculator/          # Бизнес-логика
│   ├── cache/               # In-memory кэш
│   ├── model/               # Структуры данных
│   └── middleware/          # HTTP middleware
├── config.yml               # Файл конфигурации
├── Dockerfile               # Docker конфигурация
├── Makefile                 # Автоматизация сборки
└── .golangci.yml            # Конфигурация линтера
```

## 🧪 Тестирование

Запуск всех тестов

```
go test ./... -cover
```

Отчет о покрытии кода:

```
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 📊 Качество кода

Проект использует golangci-lint со строгой конфигурацией:

```
golangci-lint run # Проверить качество кода
```