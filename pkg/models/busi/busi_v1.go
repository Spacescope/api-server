package busi

import (
	"time"
)

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

type BlockHeader struct {
	Height          int64  `xorm:"integer notnull pk"`
	Version         int    `xorm:"integer notnull pk"`
	Cid             string `xorm:"varchar(255) pk notnull index default ''"`
	Miner           string `xorm:"varchar(255) notnull default ''"`
	ParentWeight    string `xorm:"varchar(100) notnull default ''"`
	ParentBaseFee   string `xorm:"varchar(100) notnull default ''"`
	ParentStateRoot string `xorm:"varchar(100) notnull default 0"`
	WinCount        int64  `xorm:"integer notnull default 0"`
	Timestamp       uint64 `xorm:"integer notnull default 0"`
	ForkSignaling   uint64 `xorm:"integer notnull default 0"`
	CreatedAt       int64  `xorm:"created"`
}

func (b *BlockHeader) TableName() string {
	return "block_header"
}

type EVMBlockHeader struct {
	Height           int64  `xorm:"integer notnull pk"`
	Version          int    `xorm:"integer notnull pk"`
	Number           int64  `xorm:"integer pk notnull default 0"`
	ParentHash       string `xorm:"varchar(255) notnull default ''"`
	Sha3Uncles       string `xorm:"varchar(255) notnull default ''"`
	Miner            string `xorm:"varchar(255) notnull default ''"`
	StateRoot        string `xorm:"varchar(255) notnull default ''"`
	TransactionsRoot string `xorm:"varchar(255) notnull default ''"`
	ReceiptsRoot     string `xorm:"varchar(255) notnull default ''"`
	Difficulty       int64  `xorm:"integer notnull default 0"`
	GasLimit         int64  `xorm:"integer notnull default 0"`
	GasUsed          int64  `xorm:"integer notnull default 0"`
	Timestamp        int64  `xorm:"integer notnull default 0"`
	ExtraData        string `xorm:"text notnull default ''"`
	MixHash          string `xorm:"varchar(255) notnull default ''"`
	Nonce            string `xorm:"varchar(255) notnull default ''"`
	BaseFeePerGas    int64  `xorm:"integer notnull default 0"`
	Size             uint64 `xorm:"integer notnull default 0"`
	CreatedAt        int64  `xorm:"created"`
}

func (b *EVMBlockHeader) TableName() string {
	return "evm_block_header"
}

type BlockHeader struct {
	Height          int64  `xorm:"integer notnull pk"`
	Version         int    `xorm:"integer notnull pk"`
	Cid             string `xorm:"varchar(255) pk notnull index default ''"`
	Miner           string `xorm:"varchar(255) notnull default ''"`
	ParentWeight    string `xorm:"varchar(100) notnull default ''"`
	ParentBaseFee   string `xorm:"varchar(100) notnull default ''"`
	ParentStateRoot string `xorm:"varchar(100) notnull default 0"`
	WinCount        int64  `xorm:"integer notnull default 0"`
	Timestamp       uint64 `xorm:"integer notnull default 0"`
	ForkSignaling   uint64 `xorm:"integer notnull default 0"`
	CreatedAt       int64  `xorm:"created"`
}

func (b *BlockHeader) TableName() string {
	return "block_header"
}

type EVMBlockHeader struct {
	Height           int64  `xorm:"integer notnull pk"`
	Version          int    `xorm:"integer notnull pk"`
	Number           int64  `xorm:"integer pk notnull default 0"`
	ParentHash       string `xorm:"varchar(255) notnull default ''"`
	Sha3Uncles       string `xorm:"varchar(255) notnull default ''"`
	Miner            string `xorm:"varchar(255) notnull default ''"`
	StateRoot        string `xorm:"varchar(255) notnull default ''"`
	TransactionsRoot string `xorm:"varchar(255) notnull default ''"`
	ReceiptsRoot     string `xorm:"varchar(255) notnull default ''"`
	Difficulty       int64  `xorm:"integer notnull default 0"`
	GasLimit         int64  `xorm:"integer notnull default 0"`
	GasUsed          int64  `xorm:"integer notnull default 0"`
	Timestamp        int64  `xorm:"integer notnull default 0"`
	ExtraData        string `xorm:"text notnull default ''"`
	MixHash          string `xorm:"varchar(255) notnull default ''"`
	Nonce            string `xorm:"varchar(255) notnull default ''"`
	BaseFeePerGas    int64  `xorm:"integer notnull default 0"`
	Size             uint64 `xorm:"integer notnull default 0"`
	CreatedAt        int64  `xorm:"created"`
}

func (b *EVMBlockHeader) TableName() string {
	return "evm_block_header"
}
