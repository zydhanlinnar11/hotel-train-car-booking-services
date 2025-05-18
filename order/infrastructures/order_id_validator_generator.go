package infrastructures

import (
	"github.com/google/uuid"
)

type OrderIdValidatorGenerator struct {
}

func NewOrderIdValidatorGenerator() *OrderIdValidatorGenerator {
	return &OrderIdValidatorGenerator{}
}

func (o *OrderIdValidatorGenerator) Validate(id string) error {
	_, err := uuid.Parse(id)
	return err
}

func (o *OrderIdValidatorGenerator) Generate() string {
	return uuid.New().String()
}
