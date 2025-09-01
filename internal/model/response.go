package model

import "time"

type MortgageResponse struct {
	Result *MortgageCalculation `json:"result,omitempty"`
	Error  string               `json:"error,omitempty"`
}

type MortgageCalculation struct {
	ID         int                `json:"id,omitempty"`
	Params     MortgageParams     `json:"params"`
	Program    MortgageProgram    `json:"program"`
	Aggregates MortgageAggregates `json:"aggregates"`
}

type MortgageParams struct {
	ObjectCost     float64 `json:"object_cost"`
	InitialPayment float64 `json:"initial_payment"`
	Months         int     `json:"months"`
}

type MortgageAggregates struct {
	Rate            float64   `json:"rate"`
	LoanSum         float64   `json:"loan_sum"`
	MonthlyPayment  float64   `json:"monthly_payment"`
	Overpayment     float64   `json:"overpayment"`
	LastPaymentDate time.Time `json:"last_payment_date"`
}
