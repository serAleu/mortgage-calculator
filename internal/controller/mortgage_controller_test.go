package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"mortgage-calculator/internal/calculator"
	"mortgage-calculator/internal/model"
)

// ============================================================================
// MOCK СТРУКТУРЫ ДЛЯ ТЕСТИРОВАНИЯ
// ============================================================================

// MockCalculator имитирует работу калькулятора для изоляции тестов
// В реальном приложении calculator будет делать сложные расчеты,
// а здесь мы просто возвращаем заранее заданные значения
type MockCalculator struct {
	// result хранит результат, который должен вернуть мок
	result *model.MortgageCalculation
	// err хранит ошибку, которую должен вернуть мок
	err error
}

// Calculate - метод мока, который имитирует работу реального калькулятора
// Принимает указатель на запрос (*model.MortgageRequest) чтобы избежать копирования большой структуры
// Возвращает указатель на результат и ошибку - стандартная практика в Go
func (m *MockCalculator) Calculate(req *model.MortgageRequest) (*model.MortgageCalculation, error) {
	return m.result, m.err
}

// MockCache имитирует работу кэша для тестирования
// Теперь он полностью имплементирует интерфейс Cache
type MockCache struct {
	storedValue *model.MortgageCalculation
	storedItems []*model.MortgageCalculation // Для реализации GetAll
}

// Store - метод мока кэша, который сохраняет значение и возвращает ID
func (m *MockCache) Store(calculation *model.MortgageCalculation) int {
	m.storedValue = calculation
	m.storedItems = append(m.storedItems, calculation)
	return 1 // Возвращаем фиктивный ID для тестов
}

// GetAll - метод мока кэша, который возвращает все сохраненные расчеты
func (m *MockCache) GetAll() []*model.MortgageCalculation {
	return m.storedItems
}

// ============================================================================
// ОСНОВНОЙ ТЕСТОВЫЙ МЕТОД
// ============================================================================

// TestHandleCalculate тестирует метод handleCalculate с различными сценариями
// Table-driven tests - стандартный подход в Go для тестирования множества случаев
func TestHandleCalculate(t *testing.T) {
	// Фиксированное время для тестов чтобы избежать flaky tests
	// (тестов, которые могут проходить или не проходить в зависимости от времени)
	testTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// Тестовые случаи: каждый элемент массива - отдельный тестовый сценарий
	tests := []struct {
		name           string                     // Название теста для читаемости
		requestBody    string                     // Тело HTTP запроса в формате JSON
		mockResult     *model.MortgageCalculation // Что должен вернуть мок калькулятора
		mockError      error                      // Какую ошибку должен вернуть мок калькулятора
		expectedStatus int                        // Ожидаемый HTTP статус код ответа
		expectedBody   string                     // Ожидаемое тело ответа в формате JSON
	}{
		{
			name:        "successful calculation with base program",
			requestBody: `{"object_cost": 10000000, "initial_payment": 2000000, "months": 240, "program": {"base": true}}`,
			mockResult: &model.MortgageCalculation{
				ID: 1,
				Params: model.MortgageParams{
					ObjectCost:     10000000,
					InitialPayment: 2000000,
					Months:         240,
				},
				Program: model.MortgageProgram{
					Base: true,
				},
				Aggregates: model.MortgageAggregates{
					Rate:            7.5,
					LoanSum:         8000000,
					MonthlyPayment:  65608.42,
					Overpayment:     7740208,
					LastPaymentDate: testTime.AddDate(20, 0, 0),
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"result":{"id":1,"params":{"object_cost":10000000,"initial_payment":2000000,"months":240},"program":{"salary":false,"military":false,"base":true},"aggregates":{"rate":7.5,"loan_sum":8000000,"monthly_payment":65608.42,"overpayment":7740208,"last_payment_date":"2044-01-01T00:00:00Z"}}}`,
		},
		{
			name:        "successful calculation with military program",
			requestBody: `{"object_cost": 10000000, "initial_payment": 2000000, "months": 240, "program": {"military": true}}`,
			mockResult: &model.MortgageCalculation{
				ID: 1,
				Params: model.MortgageParams{
					ObjectCost:     10000000,
					InitialPayment: 2000000,
					Months:         240,
				},
				Program: model.MortgageProgram{
					Military: true,
				},
				Aggregates: model.MortgageAggregates{
					Rate:            5.5, // Военная ипотека обычно имеет льготную ставку
					LoanSum:         8000000,
					MonthlyPayment:  54789.31,
					Overpayment:     5149434,
					LastPaymentDate: testTime.AddDate(20, 0, 0),
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"result":{"id":1,"params":{"object_cost":10000000,"initial_payment":2000000,"months":240},"program":{"salary":false,"military":true,"base":false},"aggregates":{"rate":5.5,"loan_sum":8000000,"monthly_payment":54789.31,"overpayment":5149434,"last_payment_date":"2044-01-01T00:00:00Z"}}}`,
		},
		{
			name:           "invalid JSON syntax",
			requestBody:    `{invalid json syntax`, // Невалидный JSON
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid json"}`,
		},
		{
			name:           "missing required field - object_cost",
			requestBody:    `{"initial_payment": 2000000, "months": 240, "program": {"base": true}}`, // Нет object_cost
			expectedStatus: http.StatusBadRequest,
			// Validator вернет ошибку о пропущенном обязательном поле
		},
		{
			name:           "negative object cost",
			requestBody:    `{"object_cost": -10000000, "initial_payment": 2000000, "months": 240, "program": {"base": true}}`, // Отрицательная стоимость
			expectedStatus: http.StatusBadRequest,
			// Validator вернет ошибку о минимальном значении
		},
		{
			name:           "too many months (more than 600)",
			requestBody:    `{"object_cost": 10000000, "initial_payment": 2000000, "months": 601, "program": {"base": true}}`, // Слишком большой срок
			expectedStatus: http.StatusBadRequest,
			// Validator вернет ошибку о максимальном значении
		},
		{
			name:           "initial payment too low error from calculator",
			requestBody:    `{"object_cost": 10000000, "initial_payment": 500000, "months": 240, "program": {"base": true}}`, // Маленький первоначальный взнос
			mockError:      calculator.ErrInitialPaymentTooLow,                                                               // Калькулятор возвращает специфическую ошибку
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"initial payment too low"}`,
		},
		{
			name:           "internal calculator error",
			requestBody:    `{"object_cost": 10000000, "initial_payment": 2000000, "months": 240, "program": {"base": true}}`,
			mockError:      errors.New("database connection failed"), // Любая другая ошибка
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
		},
	}

	// Итерируемся по всем тестовым случаям
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ==================================================================
			// ПОДГОТОВКА ТЕСТОВОГО ОКРУЖЕНИЯ
			// ==================================================================

			// Создаем мок калькулятора с заранее заданным поведением
			mockCalc := &MockCalculator{
				result: tt.mockResult, // Результат, который должен вернуть мок
				err:    tt.mockError,  // Ошибка, которую должен вернуть мок
			}

			// Создаем мок кэша для проверки вызова метода Store
			mockCache := &MockCache{}

			// Создаем контроллер с моками вместо реальных зависимостей
			// Это принцип Dependency Injection - зависимости внедряются извне
			controller := &MortgageController{
				calc:  mockCalc,  // Мок калькулятора
				cache: mockCache, // Мок кэша
			}

			// ==================================================================
			// СОЗДАНИЕ HTTP ЗАПРОСА
			// ==================================================================

			// Создаем HTTP POST запрос с телом из requestBody
			// bytes.NewBufferString преобразует строку в io.Reader (интерфейс для чтения)
			req, err := http.NewRequest(
				"POST",                                // HTTP метод
				"/calculate",                          // URL пути
				bytes.NewBufferString(tt.requestBody), // Тело запроса
			)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			// Устанавливаем Content-Type заголовок для указания формата данных
			// Это важно для правильного парсинга на стороне сервера
			req.Header.Set("Content-Type", "application/json")

			// ==================================================================
			// СОЗДАНИЕ HTTP РЕСПОНСА ДЛЯ ЗАХВАТА ВЫВОДА
			// ==================================================================

			// httptest.NewRecorder создает реализацию http.ResponseWriter
			// которая записывает все в память вместо отправки по сети
			// Это позволяет проверять что именно вернул наш обработчик
			rr := httptest.NewRecorder()

			// ==================================================================
			// ВЫЗОВ ТЕСТИРУЕМОГО МЕТОДА
			// ==================================================================

			// Вызываем приватный метод handleCalculate контроллера
			// Передаем ему ResponseWriter и Request
			controller.handleCalculate(rr, req)

			// ==================================================================
			// ПРОВЕРКА РЕЗУЛЬТАТОВ ТЕСТА
			// ==================================================================

			// Проверяем HTTP статус код ответа
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf(
					"handler returned wrong status code: got %v want %v. Response body: %s",
					status, tt.expectedStatus, rr.Body.String(),
				)
			}

			// Проверяем тело ответа только если в тест кейсе указано ожидаемое тело
			if tt.expectedBody != "" {
				// Нормализуем JSON для сравнения (убирает пробелы, приводит к единому формату)
				expectedNormalized, err := normalizeJSON(tt.expectedBody)
				if err != nil {
					t.Fatalf("Failed to normalize expected JSON: %v", err)
				}

				actualNormalized, err := normalizeJSON(rr.Body.String())
				if err != nil {
					t.Fatalf("Failed to normalize actual JSON: %v", err)
				}

				// Сравниваем нормализованные JSON строки
				if actualNormalized != expectedNormalized {
					t.Errorf(
						"handler returned unexpected body:\ngot:  %v\nwant: %v",
						actualNormalized, expectedNormalized,
					)
				}
			}

			// Для успешных запросов проверяем, что результат сохранился в кэш
			if tt.mockResult != nil && tt.mockError == nil {
				// Проверяем через GetAll() вместо прямого доступа к storedValue
				allItems := mockCache.GetAll()
				if len(allItems) == 0 || allItems[len(allItems)-1] != tt.mockResult {
					t.Errorf(
						"result was not stored in cache: got %v, want %v",
						allItems, tt.mockResult,
					)
				}

				// Дополнительно проверяем, что Store вернул ID (если нужно)
				// if storedID := mockCache.Store(tt.mockResult); storedID != expectedID {
				//     t.Errorf("expected store to return ID %d, got %d", expectedID, storedID)
				// }
			}
		})
	}
}

// ============================================================================
// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ И ТЕСТЫ
// ============================================================================

// normalizeJSON нормализует JSON строку для сравнения
// Убирает пробелы, приводит к единому порядку полей и формату
func normalizeJSON(s string) (string, error) {
	var data interface{}
	// Парсим JSON в interface{} для денормализации
	if err := json.Unmarshal([]byte(s), &data); err != nil {
		return "", err
	}

	// Снова маршалим в JSON чтобы получить каноническую форму
	normalized, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(normalized), nil
}

// TestValidateProgram тестирует функцию валидации ипотечных программ
func TestValidateProgram(t *testing.T) {
	tests := []struct {
		name    string
		program model.MortgageProgram
		wantErr bool
	}{
		{
			name:    "valid base program",
			program: model.MortgageProgram{Base: true},
			wantErr: false,
		},
		{
			name:    "valid military program",
			program: model.MortgageProgram{Military: true},
			wantErr: false,
		},
		{
			name:    "valid salary program",
			program: model.MortgageProgram{Salary: true},
			wantErr: false,
		},
		{
			name:    "invalid program - no flags set",
			program: model.MortgageProgram{}, // Все флаги false
			wantErr: true,
		},
		{
			name:    "invalid program - multiple flags set",
			program: model.MortgageProgram{Base: true, Military: true}, // Несколько программ
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Вызываем функцию валидации программы
			err := validateProgram(tt.program)

			// Проверяем соответствие ожидаемого и фактического результата
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"validateProgram() error = %v, wantErr %v. Program: %+v",
					err, tt.wantErr, tt.program,
				)
			}
		})
	}
}

// TestErrorHandlingFunctions тестирует вспомогательные функции отправки ошибок
func TestErrorHandlingFunctions(t *testing.T) {
	t.Run("sendError with custom message", func(t *testing.T) {
		rr := httptest.NewRecorder()
		sendError(rr, "custom error message", http.StatusBadRequest)

		// Проверяем статус код
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}

		// Проверяем тело ответа
		expected := `{"error":"custom error message"}`
		if rr.Body.String() != expected {
			t.Errorf("Expected body %s, got %s", expected, rr.Body.String())
		}

		// Проверяем Content-Type заголовок
		if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", ct)
		}
	})

	t.Run("sendValidationError with validation error", func(t *testing.T) {
		rr := httptest.NewRecorder()

		// Создаем mock ошибку валидации
		validationErr := errors.New("field 'object_cost' is required")
		sendValidationError(rr, validationErr)

		// Проверяем что статус код 400
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}

		// Проверяем что тело ответа содержит ошибку
		var response struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Error == "" {
			t.Error("Expected error message in response")
		}
	})
}
