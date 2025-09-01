package cache

import (
	"mortgage-calculator/internal/model"
	"sync"
	"sync/atomic"
)

type Cache interface {
	Store(calculation *model.MortgageCalculation) int
	GetAll() []*model.MortgageCalculation
}

type inMemoryCache struct {
	mu      sync.RWMutex
	store   map[int]*model.MortgageCalculation
	counter atomic.Int32
}

func NewInMemoryCache() Cache {
	return &inMemoryCache{
		store: make(map[int]*model.MortgageCalculation),
	}
}

func (c *inMemoryCache) Store(calculation *model.MortgageCalculation) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	id := int(c.counter.Add(1))
	calculation.ID = id
	c.store[id] = calculation

	return id
}

func (c *inMemoryCache) GetAll() []*model.MortgageCalculation {
	c.mu.RLock()
	defer c.mu.RUnlock()

	calculations := make([]*model.MortgageCalculation, 0, len(c.store))
	for _, calc := range c.store {
		calculations = append(calculations, calc)
	}

	return calculations
}
