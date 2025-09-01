package model

type MortgageRequest struct {
	ObjectCost     float64         `json:"object_cost" validate:"required,min=0"`
	InitialPayment float64         `json:"initial_payment" validate:"required,min=0"`
	Months         int             `json:"months" validate:"required,min=1,max=600"`
	Program        MortgageProgram `json:"program" validate:"required"`
}

type MortgageProgram struct {
	Salary   bool `json:"salary"`
	Military bool `json:"military"`
	Base     bool `json:"base"`
}
