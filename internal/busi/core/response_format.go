package core

import (
	"time"

	"api-server/pkg/models/busi"
)

type ContractsList struct {
	Contracts []*Contract `json:"contracts"`
	Hits      int64       `json:"hits"`
}

type Contract struct {
	Address         string
	FilecoinAddress string
	Name            string
	Compiler        string
	Version         int64
	Balance         string
	Txns            int64
	Verified        time.Time
	License         string
}

type TxnsList struct {
	EVMTransaction []*busi.EVMTransaction `json:"evm_txns"`
	Hits           int64                  `json:"hits"`
}

type InternalTxnsList struct {
	EVMInternalTX []*busi.EVMInternalTX `json:"evm_internal_txns"`
	Hits          int64                 `json:"hits"`
}
