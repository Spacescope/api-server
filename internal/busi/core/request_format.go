package core

import (
	"errors"
)

type ListQuery struct {
	Offset int `form:"o" json:"o"`
	Limit  int `form:"l" json:"l"`
}

func (r *ListQuery) ListValidate() error {
	if r.Offset < 0 {
		return errors.New("the o(ffset) should be greater than or equal 0")
	}

	switch r.Limit {
	case 25, 50, 100:
	default:
		return errors.New("the l(imit) should be one of 25/50/100")
	}

	return nil
}
