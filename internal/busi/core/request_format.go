package core

import (
	"errors"
	"time"
)

const Day = 24 * time.Hour
const Month = 31 * Day
const Year = 365 * Day

type APIRequest struct {
	StartDate time.Time `form:"start_date" json:"start_date" time_format:"2006-01-02" example:"2022-06-01"`
	EndDate   time.Time `form:"end_date" json:"end_date" time_format:"2006-01-02" example:"2022-07-01"`
}

func (r *APIRequest) baseValidate() error {
	if r.StartDate.IsZero() && r.EndDate.IsZero() {
		// r.StartDate = utils.GetStartOfTheDay()
		// r.EndDate = utils.GetEndOfTheDay()
		return nil
	}

	if r.StartDate.After(r.EndDate) {
		return errors.New("end_date should be greater than start_date")
	}

	if r.StartDate.After(time.Now()) {
		return errors.New("not a valid time period")
	}

	return nil
}

func (r *APIRequest) Validate() error {
	if err := r.baseValidate(); err != nil {
		return err
	}

	if r.EndDate.Sub(r.StartDate).Hours() > Month.Hours() {
		return errors.New("the time window should be within 31 days")
	}

	return nil
}
