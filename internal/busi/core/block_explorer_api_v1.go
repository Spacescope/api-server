package core

import (
	"context"
	"net/http"
	"sort"

	"api-server/pkg/models/busi"
	"api-server/pkg/utils"

	log "github.com/sirupsen/logrus"
)

func ListContracts(ctx context.Context, r *ListQuery) (interface{}, *utils.BuErrorResponse) {
	var (
		c             busi.EVMContract
		contractsList ContractsList
	)

	// get the numbers of contracts
	total, err := busiTableRecordsCount(&c)
	if err != nil {
		return nil, err
	}

	contractsList.Hits = total
	if contractsList.Hits <= 0 {
		return contractsList, nil
	}

	// get contracts list
	contracts := make([]*busi.EVMContract, 0)
	if err := busiSQLExecute(r, &contracts); err != nil {
		return nil, err
	}

	contractsSlice := make([]*Contract, 0, len(contracts))

	for _, contract := range contracts {
		var (
			c Contract
		)

		c.Txns, err = evmTransactionCountWitVersion(contract.Address, contract.Version)
		if err != nil {
			log.Error(err)
			break
		}

		c.Address = contract.Address
		c.FilecoinAddress = contract.FilecoinAddress
		c.Balance = contract.Balance
		c.Version = int64(contract.Version)

		contractsSlice = append(contractsSlice, &c)
	}

	contractsList.Contracts = contractsSlice

	return contractsList, nil
}

func Getcontract(ctx context.Context, address string) (interface{}, *utils.BuErrorResponse) {
	evmContracts := make([]*busi.EVMContract, 0)

	err := utils.EngineGroup[utils.DB].Where("address = ?", address).Find(&evmContracts)
	if err != nil {
		log.Errorf("Execute sql error: %v", err)
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}

	if len(evmContracts) == 0 {
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusOK, Response: utils.ErrBlockExplorerAPIServerNotFound}
	}

	// get the heightest and new version address
	t := newContractsArr(evmContracts)
	sort.Sort(t)
	return t[len(t)-1], nil
}

func ListTXNs(ctx context.Context, address string, r *ListQuery) (interface{}, *utils.BuErrorResponse) {
	var (
		txnsList TxnsList
	)

	// get the numbers of transactions
	var t busi.EVMTransaction
	total, err := evmTransactionCount(address, &t)
	if err != nil {
		return nil, err
	}

	txnsList.Hits = total
	if txnsList.Hits <= 0 {
		return txnsList, nil
	}

	// get transactions list
	transactions := make([]*busi.EVMTransaction, 0)
	if err := evmTransactionFind(address, r, &transactions); err != nil {
		return nil, err
	}
	txnsList.EVMTransaction = transactions

	return txnsList, nil
}

func ListInternalTXNs(ctx context.Context, address string, r *ListQuery) (interface{}, *utils.BuErrorResponse) {
	var (
		internalTXNsList InternalTxnsList
	)

	// get the numbers of internal transactions

	var t busi.EVMInternalTX
	total, err := evmTransactionCount(address, &t)
	if err != nil {
		return nil, err
	}

	internalTXNsList.Hits = total
	if internalTXNsList.Hits <= 0 {
		return internalTXNsList, nil
	}

	// get internal transactions list
	internal_transactions := make([]*busi.EVMInternalTX, 0)
	if err := evmTransactionFind(address, r, &internal_transactions); err != nil {
		return nil, err
	}
	internalTXNsList.EVMInternalTX = internal_transactions

	return internalTXNsList, nil
}
