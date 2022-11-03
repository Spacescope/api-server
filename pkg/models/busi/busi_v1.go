package busi

import (
	"time"
)

type GasFeeOverviewH struct {
	StatDate              time.Time `json:"stat_date"`
	HourDate              time.Time `json:"hour_date"`
	PrecommitCostSector   float64   `json:"precommit_cost_sector"`
	ProvecommitCostSector float64   `json:"provecommit_cost_sector"`
	PreBatchCostSector    float64   `json:"pre_batch_cost_sector"`
	ProveBatchCostSector  float64   `json:"prove_batch_cost_sector"`
	BaseFee               float64   `json:"base_fee"`
	AvgPreAggCount        float64   `json:"avg_pre_agg_count"`
	AvgProveAggCount      float64   `json:"avg_prove_agg_count"`
	IsLatest              bool      `json:"-"`
	CreateAt              time.Time `json:"-"`
}

type CirculatingSupply struct {
	StatDate                    time.Time `json:"stat_date"`
	Height                      uint64    `json:"-"`
	ValueHeight                 uint64    `json:"-"`
	CirculatingFil              float64   `json:"circulating_fil"`
	CirculatingFilIncrease      float64   `json:"circulating_fil_increase"`
	MinedFil                    float64   `json:"mined_fil"`
	MinedFilIncrease            float64   `json:"mined_fil_increase"`
	VestedFil                   float64   `json:"vested_fil"`
	VestedFilIncrease           float64   `json:"vested_fil_increase"`
	ReserveDisbursedFil         float64   `json:"reserve_disbursed_fil"`
	ReserveDisbursedFilIncrease float64   `json:"reserve_disbursed_fil_increase"`
	LockedFil                   float64   `json:"locked_fil"`
	LockedFilIncrease           float64   `json:"locked_fil_increase"`
	BurntFil                    float64   `json:"burnt_fil"`
	BurntFilIncrease            float64   `json:"burnt_fil_increase"`
	IsLatest                    bool      `json:"-"`
	CreateAt                    time.Time `json:"-"`
}
