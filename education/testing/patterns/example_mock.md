Interfce testing (like mocks)

- Т.к моков в Go нет, надо прикалываться с мокированием через интерфейсы, например:

Сценарий для тестирования: делаю сервис который должен сохранять данные в БД

Шаг 1: Определяю интерфейс на стороне потребителя

// repository.go
package main

// Storage - это интерфейс, который ОПИСЫВАЕТ ПОВЕДЕНИЕ, необходимое нашему сервису.
type Storage interface {
    Save(data string) error
    Get(id int) (string, error)
}

Шаг 2: Создаю сервис который зависит от интерфейса

// service.go
package main

type Service struct {
    // Сервис зависит от ИНТЕРФЕЙСА, а не от конкретной БД, файла или API.
    storage Storage
}

func (s *Service) ProcessData(data string) error {
    // Сервису абсолютно всё равно, КАК работает storage.
    // Его волнует только ЧТО он может делать.
    return s.storage.Save(data)
}

Шаг 3: Пром реализация

// postgres_storage.go
package main

import "database/sql"

// PostgresStorage - реализация Storage для реальной PostgreSQL БД.
type PostgresStorage struct {
    db *sql.DB
}

// Реализуем метод Save для Postgres. Теперь PostgresStorage неявно удовлетворяет интерфейсу Storage.
func (p *PostgresStorage) Save(data string) error {
    _, err := p.db.Exec("INSERT INTO data (value) VALUES ($1)", data)
    return err
}
func (p *PostgresStorage) Get(id int) (string, error) { /* ... */ }

Шаг 4: Тестовая реализация (по сути, аналог mock)

// service_test.go
package main

import (
    "errors"
    "testing"
)

// MockStorage - наша ТЕСТОВАЯ реализация интерфейса Storage.
// Это не "мок" в классическом понимании, а просто другая реализация.
type MockStorage struct {
    // Мы можем добавить поля для управления поведением в тестах.
    shouldFail bool
    savedData  []string // "In-memory база" для проверки, что сохранилось
}

// Реализуем интерфейс Storage для MockStorage.
func (m *MockStorage) Save(data string) error {
    if m.shouldFail {
        return errors.New("fake storage error")
    }
    m.savedData = append(m.savedData, data) // Сохраняем данные в память для последующей проверки
    return nil
}
func (m *MockStorage) Get(id int) (string, error) { return "", nil }

// ТЕСТ
func TestService_ProcessData_Success(t *testing.T) {
    // 1. Arrange (Подготовка)
    // Создаем наш "мок" - экземпляр MockStorage
    mockStorage := &MockStorage{shouldFail: false}
    // Внедряем зависимость в сервис. Service принимает Storage, а MockStorage реализует Storage.
    service := &Service{storage: mockStorage}

    // 2. Act (Действие)
    testData := "test data"
    err := service.ProcessData(testData)

    // 3. Assert (Проверка)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    // Проверяем, что данные действительно "сохранились" в наш mock
    if len(mockStorage.savedData) != 1 || mockStorage.savedData[0] != testData {
        t.Errorf("Expected saved data to be [%s], got %v", testData, mockStorage.savedData)
    }
}

Пояснение:

Мы абстрагировались от конкретной реализации Storage. Для тестов мы подменили тяжелую PostgresStorage на легковесную MockStorage, которая реализует тот же самый интерфейс. Сервису всё равно, он работает с обоими одинаково.

Для сложных случаев можно использовать библиотеки для генерации моков (например, github.com/golang/mock или github.com/vektra/mockery/v2), но ручное создание моков для интерфейсов — очень распространённая практика