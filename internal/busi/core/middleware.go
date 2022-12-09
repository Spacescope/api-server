package core

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"api-server/pkg/models/busi"
	"api-server/pkg/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/imxyb/solc-go"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
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

	if err := utils.EngineGroup[utils.TaskDB].Limit(r.Limit, r.Offset).Desc("height").Find(rowsSlicePtr); err != nil {
		log.Errorf("Execute sql error: %v", err)
		return &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}
	return nil
}

func findCreatorTransaction(address string) (*busi.EVMTransaction, *utils.BuErrorResponse) {
	var receipt busi.EVMReceipt
	exist, err := utils.EngineGroup[utils.TaskDB].Where("`to`='' and contract_address=?", address).Get(&receipt)
	if err != nil {
		log.Errorf("Execute sql error: %v", err)
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}
	var tx busi.EVMTransaction
	if exist {
		exist, err = utils.EngineGroup[utils.TaskDB].Where("hash=?", receipt.TransactionHash).Get(&tx)
		if err != nil {
			log.Errorf("Execute sql error: %v", err)
			return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
				Response: utils.ErrBlockExplorerAPIServerInternal}
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
		log.Errorf("needs a pointer to a slice, v.Kind() = %v, reflect.Indirect(v).Kind() = %v", v.Kind(),
			reflect.Indirect(v).Kind())
		return nil
	}

	if err := utils.EngineGroup[utils.TaskDB].Where("\"from\" = ? or \"to\" = ?", address, address).
		Limit(r.Limit, r.Offset).OrderBy("height desc").Find(rowsSlicePtr); err != nil {
		log.Errorf("Execute sql error: %v", err)
		return &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}

	return nil
}

func evmTransactionCount(address string) (int64, *utils.BuErrorResponse) {
	var (
		count int64
		err   error

		t busi.EVMTransaction
	)

	if count, err = utils.EngineGroup[utils.TaskDB].Where("\"from\" = ? or \"to\" = ?", address,
		address).Count(&t); err != nil {
		return 0, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}

	return count, nil
}

func parseMethodAndParamsFromContract(input, contractAddress string) (string, string, map[string]interface{}) {
	if input == "" {
		return "unknown", "", nil
	}
	inputData, err := hex.DecodeString(input)
	if len(inputData) < 4 {
		return "unknown", "", nil
	}
	tokenABI, err := getContractABI(contractAddress)
	if err != nil {
		log.Errorf("getContractABI failed, err:%s", err)
		return fmt.Sprintf("0x%s", hex.EncodeToString(inputData[:4])), "", nil
	}
	if tokenABI == nil {
		return fmt.Sprintf("0x%s", hex.EncodeToString(inputData[:4])), "", nil
	}
	abiMethod, err := tokenABI.MethodById(inputData[:4])
	if err != nil {
		log.Errorf("tokenABI.MethodById failed, err:%s", err)
		return "unknown", "", nil
	}
	params := make(map[string]interface{})
	if err = abiMethod.Inputs.UnpackIntoMap(params, inputData[4:]); err != nil {
		return "", "", nil
	}
	return abiMethod.RawName, abiMethod.String(), params
}

var (
	cacheABI sync.Map
)

func getContractABI(contractAddress string) (*abi.ABI, error) {
	contractAddress = strings.ToLower(contractAddress)
	v, ok := cacheABI.Load(contractAddress)
	if ok {
		return v.(*abi.ABI), nil
	}

	var contractVerify busi.EVMContractVerify
	exist, err := utils.EngineGroup[utils.APIDB].Where("address=? and status=?",
		contractAddress, busi.EVMContractVerifyStatusSuccessfully).Get(&contractVerify)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}

	var (
		output    solc.Output
		abiString string
	)
	if err = json.Unmarshal([]byte(contractVerify.Output), &output); err != nil {
		return nil, err
	}
	if contractVerify.CompilerType == busi.CompilerTypeSingleFile {
		for _, c := range output.Contracts[""] {
			abiString = gjson.Get(c.Metadata, "output.abi").String()
			break
		}
	} else if contractVerify.CompilerType == busi.CompilerTypeMultiPart {
		key := fmt.Sprintf("%s.sol", contractVerify.ContractName)
		metadata := output.Contracts[key][contractVerify.ContractName].Metadata
		abiString = gjson.Get(metadata, "output.abi").String()
	}
	tokenABI, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return nil, err
	}
	cacheABI.LoadOrStore(contractAddress, &tokenABI)
	return &tokenABI, nil
}
