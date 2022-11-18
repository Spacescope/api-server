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
	total, err := utils.EngineGroup[utils.DB].Count(prt)
	if err != nil {
		log.Errorf("ListContracts execute sql error: %v\n", err)
		return 0, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}

	return total, nil
}

func busiSQLExecute(r *ListQuery, rowsSlicePtr interface{}) *utils.BuErrorResponse {
	v := reflect.ValueOf(rowsSlicePtr)
	if v.Kind() != reflect.Ptr || reflect.Indirect(v).Kind() != reflect.Slice {
		log.Errorf("needs a pointer to a slice, v.Kind() = %v, reflect.Indirect(v).Kind() = %v", v.Kind(),
			reflect.Indirect(v).Kind())
		return nil
	}

	if err := utils.EngineGroup[utils.DB].Limit(r.Limit, r.Offset).Find(rowsSlicePtr); err != nil {
		log.Errorf("Execute sql error: %v", err)
		return &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}
	return nil
}

func getContractVerify(address string, m *busi.EVMContractVerify) {

}

func evmTransactionFind(address string, r *ListQuery, rowsSlicePtr interface{}) *utils.BuErrorResponse {
	v := reflect.ValueOf(rowsSlicePtr)
	if v.Kind() != reflect.Ptr || reflect.Indirect(v).Kind() != reflect.Slice {
		log.Errorf("needs a pointer to a slice, v.Kind() = %v, reflect.Indirect(v).Kind() = %v", v.Kind(),
			reflect.Indirect(v).Kind())
		return nil
	}

	if err := utils.EngineGroup[utils.DB].Where("\"from\" = ? or \"to\" = ?", address, address).Limit(r.Limit,
		r.Offset).Find(rowsSlicePtr); err != nil {
		log.Errorf("Execute sql error: %v", err)
		return &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}

	return nil
}

func evmTransactionCountWitVersion(address string, version int) (int64, *utils.BuErrorResponse) {
	var (
		count int64
		err   error
		t     busi.EVMTransaction
	)

	if count, err = utils.EngineGroup[utils.DB].Where("(\"from\" = ? or \"to\" = ?) and version = ?", address, address,
		version).Count(&t); err != nil {
		return 0, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}

	return count, nil
}

func evmTransactionCount(address string /*, version int*/, t interface{}) (int64, *utils.BuErrorResponse) {
	var (
		count int64
		err   error
	)

	// if count, err = utils.EngineGroup[utils.DB].Where("(\"from\" = ? or \"to\" = ?) and version = ?", address, address, version).Count(&t); err != nil {
	if count, err = utils.EngineGroup[utils.DB].Where("\"from\" = ? or \"to\" = ?", address,
		address).Count(t); err != nil {
		return 0, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}

	return count, nil
}
