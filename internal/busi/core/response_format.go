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
	Height          int64
	Address         string
	FilecoinAddress string
	Name            string
	CompilerType    int
	CompilerVersion string
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

type EVMTransaction struct {
	busi.EVMTransaction `json:",inline"`
	ToIsContract        bool  `json:"to_is_contract"`
	ConfirmationBlocks  int64 `json:"confirmation_blocks"`
}

type InternalTxnsList struct {
	EVMInternalTX []*busi.EVMInternalTX `json:"evm_internal_txns"`
	Hits          int64                 `json:"hits"`
}

type CompileVersionList struct {
	Versions []*CompileVersion `json:"versions"`
}

type CompileVersion struct {
	Version     string `json:"version"`
	LongVersion string `json:"long_version"`
}

type SourceCode struct {
	FileName string `json:"file_name"`
	Code     string `json:"code"`
}

type ContractVerify struct {
	busi.EVMContractVerify
	SourceCodes []*SourceCode `json:"source_codes"`
	ABI         string        `json:"abi"`
	ErrMsg      string        `json:"err_msg"`
	Bytecode    string        `json:"bytecode"`
}

type ContractDetail struct {
	Address         string        `json:"address"`
	FilecoinAddress string        `json:"filecoin_address"`
	Balance         string        `json:"balance"`
	Nonce           uint64        `json:"nonce"`
	ByteCode        string        `json:"byte_code"`
	Creator         string        `json:"creator"`
	Txn             string        `json:"txn"`
	CompilerType    int           `json:"compiler_type"`
	CompilerVersion string        `json:"compiler_version"`
	LicenseType     string        `json:"license_type"`
	ContractName    string        `json:"contract_name"`
	Verified        time.Time     `json:"verified"`
	ABI             string        `json:"abi"`
	SourceCodes     []*SourceCode `json:"source_codes"`
}

type ContractIsVerify struct {
	IsVerify bool `json:"is_verify"`
}

type Address struct {
	Height          int64  `json:"height"`
	Address         string `json:"address"`
	FilecoinAddress string `json:"filecoin_address"`
	Balance         string `json:"balance"`
	Nonce           uint64 `json:"nonce"`
}

type EventList struct {
	Events []*Event `json:"events"`
}

type Event struct {
	Address      string                 `json:"address"`
	RawTopics    []string               `json:"raw_topics"`
	ParsedTopics map[string]interface{} `json:"parsed_topics"`
	RawData      string                 `json:"raw_data"`
	ParsedData   map[string]interface{} `json:"parsed_data"`
	BlockNumber  uint64                 `json:"block_number"`
	TxHash       string                 `json:"tx_hash"`
	TxIndex      uint                   `json:"tx_index"`
	BlockHash    string                 `json:"block_hash"`
	Index        uint                   `json:"index"`
	EventName    string                 `json:"event_name"`
}
