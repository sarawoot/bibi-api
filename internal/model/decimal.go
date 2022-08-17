package model

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

type Decimal struct {
	decimal.Decimal
}

func (m Decimal) MarshalJSON() ([]byte, error) {
	f, _ := m.Float64()
	return json.Marshal(f)
}

func (m *Decimal) UnmarshalJSON(data []byte) error {
	var a float64
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	m.Decimal = decimal.NewFromFloat(a)
	return nil
}

func (m Decimal) ToFloat64() float64 {
	v, _ := m.Float64()
	return v
}
