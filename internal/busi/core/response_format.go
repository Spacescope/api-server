package core

import (
	"time"
)

type Contract struct {
	Address         string
	FilecoinAddress string
	Name            string
	Compiler        string
	Version         int64
	Balance         uint64
	Txns            int64
	Verified        time.Time
	License         string
}
