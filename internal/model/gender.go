package model

import (
	"encoding/json"
	"errors"
)

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
	LGBTQ  Gender = "lgbtq+"
)

func (g Gender) MarshalJSON() ([]byte, error) {
	return []byte(g), nil
}

func (g *Gender) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	gender := Gender(s)
	switch gender {
	case Male, Female, LGBTQ:
		*g = gender
		return nil
	}
	return errors.New("Invalid Gender")
}
