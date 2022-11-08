package core

import (
	"api-server/pkg/models/busi"
	"api-server/pkg/utils"
	"context"
	"net/http"
	"reflect"

	log "github.com/sirupsen/logrus"
)

func ListContracts(ctx context.Context, r *ListQuery) (interface{}, *utils.BuErrorResponse) {
	contracts := make([]*busi.EVMContract, 0)

	if err := busiSQLExecute(r, &contracts); err != nil {
		return nil, err
	}

	var err error
	contractsResp := make([]*Contract, 0, len(contracts))

	for _, contract := range contracts {
		var c Contract

		c.Txns, err = evmTransactionCount(contract.Address, contract.Version)
		if err != nil {
			log.Error(err)
			break
		}

		c.Address = contract.Address
		c.FilecoinAddress = contract.FilecoinAddress
		c.Balance = uint64(contract.Balance)
		c.Version = int64(contract.Version)

		contractsResp = append(contractsResp, &c)
	}

	return contractsResp, nil
}

func Getcontract(ctx context.Context, address string) (interface{}, *utils.BuErrorResponse) {
	var evmContract busi.EVMContract

	b, err := utils.EngineGroup[utils.DB].Where("address = ?", address).Get(&evmContract)
	if err != nil {
		log.Errorf("Execute sql error: %v", err)
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError, Response: utils.ErrBlockExplorerAPIServerInternal}
	}

	if !b {
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError, Response: utils.ErrBlockExplorerAPIServerNotFound}
	}

	return evmContract, nil
}

func ListTXNs(ctx context.Context, address string, r *ListQuery) (interface{}, *utils.BuErrorResponse) {
	transactions := make([]*busi.EVMTransaction, 0)
	return transactions, evmTransactionFind(address, r, &transactions)
}

func ListInternalTXNs(ctx context.Context, address string, r *ListQuery) (interface{}, *utils.BuErrorResponse) {
	internal_transactions := make([]*busi.EVMInternalTX, 0)
	return internal_transactions, evmTransactionFind(address, r, &internal_transactions)
}

func busiSQLExecute(r *ListQuery, rowsSlicePtr interface{}) *utils.BuErrorResponse {
	v := reflect.ValueOf(rowsSlicePtr)
	if v.Kind() != reflect.Ptr || reflect.Indirect(v).Kind() != reflect.Slice {
		log.Errorf("needs a pointer to a slice, v.Kind() = %v, reflect.Indirect(v).Kind() = %v", v.Kind(), reflect.Indirect(v).Kind())
		return nil
	}

	if err := utils.EngineGroup[utils.DB].Limit(r.Limit, r.Offset).Find(rowsSlicePtr); err != nil {
		log.Errorf("Execute sql error: %v", err)
		return &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError, Response: utils.ErrBlockExplorerAPIServerInternal}
	}
	return nil
}

func evmTransactionFind(address string, r *ListQuery, rowsSlicePtr interface{}) *utils.BuErrorResponse {
	v := reflect.ValueOf(rowsSlicePtr)
	if v.Kind() != reflect.Ptr || reflect.Indirect(v).Kind() != reflect.Slice {
		log.Errorf("needs a pointer to a slice, v.Kind() = %v, reflect.Indirect(v).Kind() = %v", v.Kind(), reflect.Indirect(v).Kind())
		return nil
	}

	if err := utils.EngineGroup[utils.DB].Where("\"from\" = ? or \"to\" = ?", address, address).Limit(r.Limit, r.Offset).Find(rowsSlicePtr); err != nil {
		log.Errorf("Execute sql error: %v", err)
		return &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError, Response: utils.ErrBlockExplorerAPIServerInternal}
	}

	return nil
}

func evmTransactionCount(address string, version int) (int64, error) {
	var (
		count int64
		err   error
		t     busi.EVMTransaction
	)

	if count, err = utils.EngineGroup[utils.DB].Where("(\"from\" = ? or \"to\" = ?) and version = ?", address, address, version).Count(&t); err != nil {
		return 0, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError, Response: utils.ErrBlockExplorerAPIServerInternal}
	}

	return count, nil
}
