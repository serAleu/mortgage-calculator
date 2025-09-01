package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"mortgage-calculator/internal/cache"
	"mortgage-calculator/internal/calculator"
	"mortgage-calculator/internal/model"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type MortgageController struct {
	calc  calculator.Calculator
	cache cache.Cache
}

func NewMortgageController(calc calculator.Calculator, cache cache.Cache) *MortgageController {
	return &MortgageController{calc: calc, cache: cache}
}

func (c *MortgageController) RegisterRoutes(r *chi.Mux) {
	r.Post("/execute", c.handleCalculate)
	r.Get("/cache", c.handleGetCache)
}

func (c *MortgageController) handleCalculate(w http.ResponseWriter, r *http.Request) {
	var req model.MortgageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	// Validate program selection
	if err := validateProgram(req.Program); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate other fields
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		sendValidationError(w, err)
		return
	}

	// Business logic calculation
	result, err := c.calc.Calculate(&req)
	if err != nil {
		if errors.Is(err, calculator.ErrInitialPaymentTooLow) {
			sendError(w, err.Error(), http.StatusBadRequest)
			return
		}
		sendError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Store in cache
	c.cache.Store(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.MortgageResponse{Result: result})
}

func (c *MortgageController) handleGetCache(w http.ResponseWriter, r *http.Request) {
	calculations := c.cache.GetAll()
	if len(calculations) == 0 {
		sendError(w, "empty cache", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(calculations)
}

func validateProgram(program model.MortgageProgram) error {
	count := 0
	if program.Salary {
		count++
	}
	if program.Military {
		count++
	}
	if program.Base {
		count++
	}

	switch count {
	case 0:
		return errors.New("choose program")
	case 1:
		return nil
	default:
		return errors.New("choose only 1 program")
	}
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model.MortgageResponse{Error: message})
}

func sendValidationError(w http.ResponseWriter, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		sendError(w, fmt.Sprintf("validation error: %s", validationErrors[0].Error()), http.StatusBadRequest)
		return
	}
	sendError(w, "invalid input", http.StatusBadRequest)
}
