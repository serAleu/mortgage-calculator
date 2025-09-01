package calculator

import (
	"math"
	"mortgage-calculator/internal/model"
	"time"
)

type Calculator interface {
	Calculate(request *model.MortgageRequest) (*model.MortgageCalculation, error)
}

type calculatorImpl struct{}

func NewCalculator() Calculator {
	return &calculatorImpl{}
}

func (c *calculatorImpl) Calculate(req *model.MortgageRequest) (*model.MortgageCalculation, error) {
	// Validate initial payment (20%)
	minInitialPayment := req.ObjectCost * 0.2
	if req.InitialPayment < minInitialPayment {
		return nil, ErrInitialPaymentTooLow
	}

	// Determine interest rate based on program
	annualRate := c.getAnnualRate(req.Program)
	monthlyRate := annualRate / 12 / 100

	// Calculate loan sum
	loanSum := req.ObjectCost - req.InitialPayment

	// Calculate annuity payment
	annuityCoeff := (monthlyRate * math.Pow(1+monthlyRate, float64(req.Months))) /
		(math.Pow(1+monthlyRate, float64(req.Months)) - 1)
	monthlyPayment := loanSum * annuityCoeff

	// Calculate overpayment
	totalPayment := monthlyPayment * float64(req.Months)
	overpayment := totalPayment - loanSum

	// Calculate last payment date
	lastPaymentDate := time.Now().AddDate(0, req.Months, 0)

	return &model.MortgageCalculation{
		Params: model.MortgageParams{
			ObjectCost:     req.ObjectCost,
			InitialPayment: req.InitialPayment,
			Months:         req.Months,
		},
		Program: req.Program,
		Aggregates: model.MortgageAggregates{
			Rate:            annualRate,
			LoanSum:         loanSum,
			MonthlyPayment:  math.Round(monthlyPayment),
			Overpayment:     math.Round(overpayment),
			LastPaymentDate: lastPaymentDate,
		},
	}, nil
}

func (c *calculatorImpl) getAnnualRate(program model.MortgageProgram) float64 {
	switch {
	case program.Salary:
		return 8
	case program.Military:
		return 9
	case program.Base:
		return 10
	default:
		return 0
	}
}

var (
	ErrInitialPaymentTooLow = &BusinessError{"the initial payment should be more"}
)

type BusinessError struct {
	Message string
}

func (e *BusinessError) Error() string {
	return e.Message
}
