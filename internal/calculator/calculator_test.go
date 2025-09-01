package calculator

import (
	"mortgage-calculator/internal/model"
	"testing"
)

func TestCalculator_Calculate(t *testing.T) {
	tests := []struct {
		name        string
		request     *model.MortgageRequest
		wantPayment float64
		wantError   error
	}{
		{
			name: "valid corporate program",
			request: &model.MortgageRequest{
				ObjectCost:     5_000_000,
				InitialPayment: 1_000_000,
				Months:         240,
				Program:        model.MortgageProgram{Salary: true},
			},
			wantPayment: 33458,
			wantError:   nil,
		},
		{
			name: "initial payment too low",
			request: &model.MortgageRequest{
				ObjectCost:     5_000_000,
				InitialPayment: 500_000, // 10%, should be 20% (1M)
				Months:         240,
				Program:        model.MortgageProgram{Salary: true},
			},
			wantPayment: 0,
			wantError:   ErrInitialPaymentTooLow,
		},
	}

	calc := NewCalculator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Calculate(tt.request)

			if tt.wantError != nil {
				if err == nil || err.Error() != tt.wantError.Error() {
					t.Errorf("Expected error %v, got %v", tt.wantError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result.Aggregates.MonthlyPayment != tt.wantPayment {
				t.Errorf("Expected monthly payment %f, got %f", tt.wantPayment, result.Aggregates.MonthlyPayment)
			}
		})
	}
}
