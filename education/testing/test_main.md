TestMain

TestMain — это механизм для выполнения одноразовых операций до и после запуска всех тестов в конкретном пакете (*_test.go).
Аналог из Java: @BeforeAll / @AfterAll из JUnit 5. Но в Go это реализовано на уровне всего пакета, а не класса.

1. Типичные use-cases:
    * поднятие и остановка тестовой базы данных
    * запуск и остановка HTTP-сервера для интеграционных тестов
    * установка соединения с внешним сервисом (например, Redis)
    * глобальная инициализация тестовых данных 

2. Как это работает:
    * создаю одну функцию TestMain в тестовом файле
    * go test обнаруживает её и передаёт ей управление вместо того, чтобы запускать тесты напрямую
    * сам вызываю m.Run(), который и запускает все тесты в пакете
    * контролирую что происходит до (setup()) и после (shutdown()) запуска тестов.
    * завершает работу, когда передается код возврата тестов операционной системе (os.Exit(code))

Пример:

package user

import (
    "database/sql"
    "fmt"
    "os"
    "testing"

    _ "github.com/lib/pq" // Драйвер PostgreSQL
)

// Глобальные переменные для общих ресурсов
var testDB *sql.DB

// setup создает тестовую БД и применяет миграции
func setup() error {
    var err error
    // Подключаемся к тестовой БД (например, к специально созданной testdb)
    connStr := "user=postgres dbname=testdb sslmode=disable"
    testDB, err = sql.Open("postgres", connStr)
    if err != nil {
        return err
    }
    return testDB.Ping() // Проверяем соединение
}

// shutdown очищает тестовую БД и закрывает соединение
func shutdown() {
    if testDB != nil {
        // Выполняем очистку: DROP TABLE, etc.
        testDB.Exec("DROP TABLE IF EXISTS users")
        testDB.Close()
    }
}

// TestMain - точка входа для всех тестов в пакете 'user'
func TestMain(m *testing.M) {
    // 1. Setup
    err := setup()
    if err != nil {
        fmt.Printf("Failed to setup: %v\n", err)
        os.Exit(1) // Выходим с ошибкой, если не удалось настроить
    }

    // 2. Запускаем все тесты в этом пакете (TestXxx функции)
    exitCode := m.Run()

    // 3. Shutdown (выполнится ВНЕЗАВИСИМО от того, упали тесты или нет)
    shutdown()

    // 4. Передаем код возврата тестов ОС
    os.Exit(exitCode)
}

// Теперь все тесты могут использовать testDB
func TestUserRepository_Create(t *testing.T) {
    repo := NewRepository(testDB) // Передаем глобальное соединение с БД
    // ... логика теста, использующая repo
}

Важные нюансы:
1. TestMain выполняется один раз на пакет
2. Порядок выполнения: setup() -> ALL TestXxx() -> shutdown()
3. shutdown() выполняется всегда, даже если тесты упали (благодаря os.Exit(exitCode))
4. не использовать TestMain для юнит-тестов! Он нужен для тяжёлой интеграции. Для изолированных юнитов он избыточен