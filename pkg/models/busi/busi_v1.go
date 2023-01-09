package busi

import (
	"time"
)

type EVMBlockHeader struct {
	Height           int64  `json:"height"`
	Version          int    `json:"version"`
	Number           int64  `json:"number"`
	Hash             string `json:"hash"`
	ParentHash       string `json:"parentHash"`
	Sha3Uncles       string `json:"sha3_uncles"`
	Miner            string `json:"miner"`
	StateRoot        string `json:"state_root"`
	TransactionsRoot string `json:"transactions_root"`
	ReceiptsRoot     string `json:"receipts_root"`
	Difficulty       int64  `json:"difficulty"`
	GasLimit         int64  `json:"gas_limit"`
	GasUsed          int64  `json:"gas_used"`
	Timestamp        int64  `json:"timestamp"`
	ExtraData        string `json:"extra_data"`
	MixHash          string `json:"mix_hash"`
	Nonce            string `json:"nonce"`
	BaseFeePerGas    string `json:"base_fee_per_gas"`
	Size             uint64 `json:"size"`
}

func (b *EVMBlockHeader) TableName() string {
	return "evm_block_header"
}

// Contract evm smart contract
type EVMContract struct {
	Height          int64  `json:"height"`
	Version         int    `json:"version"`
	Address         string `json:"address"`
	FilecoinAddress string `json:"filecoin_address"`
	Balance         string `json:"balance"`
	Nonce           uint64 `json:"nonce"`
	ByteCode        string `json:"byte_code"`
}

func (c *EVMContract) TableName() string {
	return "evm_contract"
}

// InternalTX contract internal transaction
type EVMInternalTX struct {
	Height     int64  `json:"height"`
	Version    int    `json:"version"`
	Hash       string `json:"hash"`
	ParentHash string `json:"parent_hash"`
	Type       uint64 `json:"type"`
	From       string `json:"from"`
	To         string `json:"to"`
	Value      string `json:"value"`
}

func (i *EVMInternalTX) TableName() string {
	return "evm_internal_tx"
}

// Receipt evm transaction receipt
type EVMReceipt struct {
	Height            int64  `json:"height"`
	Version           int    `json:"version"`
	TransactionHash   string `json:"transaction_hash"`
	TransactionIndex  int64  `json:"transaction_index"`
	BlockHash         string `json:"block_hash"`
	BlockNumber       int64  `json:"block_number"`
	From              string `json:"from"`
	To                string `json:"to"`
	StateRoot         string `json:"state_root"`
	Status            int64  `json:"status"`
	ContractAddress   string `json:"contract_address"`
	CumulativeGasUsed int64  `json:"cumulative_gas_used"`
	GasUsed           int64  `json:"gas_used"`
	EffectiveGasPrice int64  `json:"effective_gas_price"`
	LogsBloom         string `json:"logs_bloom"`
	Logs              string `json:"logs"`
}

func (r *EVMReceipt) TableName() string {
	return "evm_receipt"
}

// Transaction evm transaction
type EVMTransaction struct {
	Height               int64  `json:"height"`
	Version              int    `json:"version"`
	Hash                 string `json:"hash"`
	ChainID              uint64 `json:"chain_id"`
	Nonce                uint64 `json:"nonce"`
	BlockHash            string `json:"block_hash"`
	BlockNumber          uint64 `json:"block_number"`
	TransactionIndex     uint64 `json:"transaction_index"`
	From                 string `json:"from"`
	To                   string `json:"to"`
	Value                string `json:"value"`
	Type                 uint64 `json:"type"`
	Input                string `json:"input"`
	GasLimit             uint64 `json:"gas_limit"`
	MaxFeePerGas         string `json:"max_fee_per_gas"`
	MaxPriorityFeePerGas string `json:"max_priority_fee_per_gas"`
	V                    string `json:"v"`
	R                    string `json:"r"`
	S                    string `json:"s"`

	MethodName string                 `xorm:"-" json:"method_name"`
	MethodSig  string                 `xorm:"-" json:"method_sig"`
	Params     map[string]interface{} `xorm:"-" json:"params"`
}

func (m *EVMTransaction) TableName() string {
	return "evm_transaction"
}

const (
	CompilerTypeSingleFile   = 1
	CompilerTypeMultiPart    = 2
	CompilerTypeStdJsonInput = 3

	EVMContractVerifyStatusDoing        = 0
	EVMContractVerifyStatusSuccessfully = 1
	EVMContractVerifyStatusNotEqual     = 2
	EVMContractVerifyStatusUnknown      = 3
)

type EVMContractVerify struct {
	ID              int64     `xorm:"pk autoincr" json:"id"`
	Address         string    `xorm:"varchar(255) notnull default '' index" json:"address"`
	CompilerType    int       `xorm:"int notnull default 1" json:"compiler_type"`
	CompilerVersion string    `xorm:"varchar(100) notnull default ''" json:"compiler_version"`
	LicenseType     string    `xorm:"varchar(255) notnull default ''" json:"license_type"`
	ContractName    string    `xorm:"varchar(100) notnull default ''" json:"contract_name"`
	Input           string    `xorm:"text notnull default ''" json:"-"`
	Output          string    `xorm:"text notnull default ''" json:"-"`
	Status          int       `xorm:"int notnull default 0" json:"status"`
	CreateAt        time.Time `xorm:"created" json:"create_at"`
	UpdatedAt       time.Time `xorm:"updated" json:"updated_at"`
}

func (c *EVMContractVerify) TableName() string {
	return "evm_contract_verify"
}

// EVMAddress evm address
type EVMAddress struct {
	Height          int64  `xorm:"bigint notnull pk" json:"height"`
	Version         int    `xorm:"integer notnull pk" json:"version"`
	Address         string `xorm:"varchar(255) notnull pk" json:"address"`
	FilecoinAddress string `xorm:"varchar(255) notnull default ''" json:"filecoin_address"`
	Balance         string `xorm:"varchar(100) notnull default '0'" json:"balance"`
	Nonce           uint64 `xorm:"bigint notnull default 0" json:"nonce"`
}

func (a *EVMAddress) TableName() string {
	return "evm_address"
}

type FVMSummaryDaily struct {
	StatDate                  time.Time
	TotalContractCount        int64
	TotalInternalTxnCount     int64
	TotalExternalTxnCount     int64
	TotalTxnCount             int64
	TotalDeployerAddressCount int64
	TotalNetworkAddressCount  int64
	IsLatest                  bool
	CreateAt                  time.Time
}

func (f *FVMSummaryDaily) TableName() string {
	return "fvm_summary_daily"
}

type FVMContractSummaryDaily struct {
	StatDate        time.Time
	ContractAddress string
	FilecoinAddress string
	TxnCount        int64
	UserCount       int64
	IsLatest        bool
	CreateAt        time.Time
}

func (f *FVMContractSummaryDaily) TableName() string {
	return "fvm_contract_summary_daily"
}

type FVMContractCountDaily struct {
	StatDate                   time.Time
	ContractAddress            string
	NewContractCountDaily      int
	NewContractCountWeekly     int
	NewContractCountMonthly    int
	ActiveContractCountDaily   int
	ActiveContractCountWeekly  int
	ActiveContractCountMonthly int
	IsLatest                   bool
	CreateAt                   time.Time
}

func (f *FVMContractCountDaily) TableName() string {
	return "fvm_contract_count_daily"
}

type FVMUserCountDaily struct {
	StatDate               time.Time
	ActiveUserCountDaily   int64
	ActiveUserCountWeekly  int64
	ActiveUserCountMonthly int64
	IsLatest               bool
	CreateAt               time.Time
}

func (f *FVMUserCountDaily) TableName() string {
	return "fvm_user_count_daily"
}

type FVMTotalValueLockedDaily struct {
	StatDate              time.Time
	Fil2usd               float64 `xorm:"'fil2usd'"`
	TotalValueLocked      float64
	TotalValueLockedUsd   float64
	TotalValueReceived    float64
	TotalValueReceivedUsd float64
	TotalValueSent        float64
	IsLatest              bool
	CreateAt              time.Time
}

func (f *FVMTotalValueLockedDaily) TableName() string {
	return "fvm_total_value_locked_daily"
}

type FVMContractCallCountDaily struct {
	StatDate        time.Time
	ContractAddress string
	CallDirection   string
	CallCount       int64
	IsLatest        bool
	CreateAt        time.Time
}

func (f *FVMContractCallCountDaily) TableName() string {
	return "fvm_contract_call_count_daily"
}

type FVMStatDataIsReady struct {
	Tablename      string `xorm:"'table_name'"`
	LatestStatDate time.Time
}

func (f *FVMStatDataIsReady) TableName() string {
	return "fvmstat_data_is_ready"
}
