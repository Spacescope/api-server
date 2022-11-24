package core

import (
	"net/http"
	"reflect"

	"api-server/pkg/models/busi"
	"api-server/pkg/utils"

	log "github.com/sirupsen/logrus"
)

type contractsArr []*busi.EVMContract

func newContractsArr(t []*busi.EVMContract) contractsArr {
	return contractsArr(t)
}

func (x contractsArr) Len() int {
	return len(x)
}

func (x contractsArr) Less(i, j int) bool {
	if x[i].Height == x[j].Height {
		return x[i].Version < x[j].Version
	} else {
		return x[i].Height < x[j].Height
	}
}

func (x contractsArr) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func busiTableRecordsCount(prt interface{}) (int64, *utils.BuErrorResponse) {
	total, err := utils.EngineGroup[utils.TaskDB].Count(prt)
	if err != nil {
		log.Errorf("ListContracts execute sql error: %v", err)
		return 0, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError, Response: utils.ErrBlockExplorerAPIServerInternal}
	}

	return total, nil
}

func busiSQLExecute(r *ListQuery, rowsSlicePtr interface{}) *utils.BuErrorResponse {
	v := reflect.ValueOf(rowsSlicePtr)
	if v.Kind() != reflect.Ptr || reflect.Indirect(v).Kind() != reflect.Slice {
		log.Errorf("needs a pointer to a slice, v.Kind() = %v, reflect.Indirect(v).Kind() = %v", v.Kind(), reflect.Indirect(v).Kind())
		return nil
	}

	if err := utils.EngineGroup[utils.TaskDB].Limit(r.Limit, r.Offset).Find(rowsSlicePtr); err != nil {
		log.Errorf("Execute sql error: %v", err)
		return &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError, Response: utils.ErrBlockExplorerAPIServerInternal}
	}
	return nil
}

func findCreatorTransaction(address string) (*busi.EVMTransaction, *utils.BuErrorResponse) {
	var receipt busi.EVMReceipt
	exist, err := utils.EngineGroup[utils.TaskDB].Where("`to`='' and contract_address=?", address).Get(&receipt)
	if err != nil {
		log.Errorf("Execute sql error: %v", err)
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError, Response: utils.ErrBlockExplorerAPIServerInternal}
	}
	var tx busi.EVMTransaction
	if exist {
		exist, err = utils.EngineGroup[utils.TaskDB].Where("hash=?", receipt.TransactionHash).Get(&tx)
		if err != nil {
			log.Errorf("Execute sql error: %v", err)
			return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError, Response: utils.ErrBlockExplorerAPIServerInternal}
		}
		if exist {
			return &tx, nil
		}
	}
	return nil, nil
}

func evmTransactionFind(address string, r *ListQuery, rowsSlicePtr interface{}) *utils.BuErrorResponse {
	v := reflect.ValueOf(rowsSlicePtr)
	if v.Kind() != reflect.Ptr || reflect.Indirect(v).Kind() != reflect.Slice {
		log.Errorf("needs a pointer to a slice, v.Kind() = %v, reflect.Indirect(v).Kind() = %v", v.Kind(), reflect.Indirect(v).Kind())
		return nil
	}

	if err := utils.EngineGroup[utils.TaskDB].Where("\"from\" = ? or \"to\" = ?", address, address).Limit(r.Limit, r.Offset).Find(rowsSlicePtr); err != nil {
		log.Errorf("Execute sql error: %v", err)
		return &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError, Response: utils.ErrBlockExplorerAPIServerInternal}
	}

	return nil
}

func evmTransactionCount(address string) (int64, *utils.BuErrorResponse) {
	var (
		count int64
		err   error

		t busi.EVMTransaction
	)

	if count, err = utils.EngineGroup[utils.TaskDB].Where("\"from\" = ? or \"to\" = ?", address, address).Count(&t); err != nil {
		return 0, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError, Response: utils.ErrBlockExplorerAPIServerInternal}
	}

	return count, nil
}
