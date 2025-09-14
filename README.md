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

## 🛠️ Makefile Commands Reference

Development Commands
| Command | Description |
|---------|-------------|
| `make all` | Build, test and lint the project |
| `make build` | Build the binary |
| `make run` | Run the application locally |
| `make stop` | Stop the local application |
| `make clean` | Clean build artifacts and test cache |

Testing Commands
| Command | Description |
|---------|-------------|
| `make test` | Run tests with coverage and race detection |
| `make test-verbose` | Run tests with verbose output |
| `make test-coverage` | Run tests and show coverage report |
| `make test-package pkg=./path` | Run tests for specific package |

Code Quality Commands
| Command | Description |
|---------|-------------|
| `make lint` | Run golangci-lint |
| `make lint-fix` | Run golangci-lint with auto-fix |
| `make fmt` | Format Go code |

Docker Commands
| Command | Description |
|---------|-------------|
| `make docker-build` | Build Docker image |
| `make docker-run` | Run Docker container |
| `make docker-stop` | Stop Docker container |
| `make docker-rm` | Remove Docker container |
| `make docker-clean` | Stop and remove Docker container |
| `make docker-logs` | Show Docker container logs |
| `make docker-push` | Build and push to Docker Hub |
| `make docker-pull` | Pull from Docker Hub |
| `make docker-shell` | Open shell in container |

Dependency Management
| Command | Description |
|---------|-------------|
| `make deps` | Download dependencies |
| `make deps-update` | Update dependencies |
| `make deps-vendor` | Vendor dependencies |

Utility Commands
| Command | Description |
|---------|-------------|
| `make help` | Show help message |
| `make version` | Show Go version |
 
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