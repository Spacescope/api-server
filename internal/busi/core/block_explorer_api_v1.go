package core

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"api-server/pkg/models/busi"
	"api-server/pkg/utils"
	"github.com/goccy/go-json"
	"github.com/imxyb/solc-go"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
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

		contractVerify, err := GetSuccessContractVerifyByAddress(ctx, c.Address)
		if err != nil && err.HttpCode != http.StatusNotFound {
			return nil, err
		}
		if contractVerify != nil {
			c.Name = contractVerify.ContractName
			c.Compiler = contractVerify.CompilerVersion
			c.License = contractVerify.LicenseType
			c.Verified = contractVerify.CreateAt
		}

		contractsSlice = append(contractsSlice, &c)
	}

	contractsList.Contracts = contractsSlice

	return contractsList, nil
}

func Getcontract(ctx context.Context, address string) (interface{}, *utils.BuErrorResponse) {
	evmContracts := make([]*busi.EVMContract, 0)

	err := utils.EngineGroup[utils.DB].Where("address ilike ?", address).Find(&evmContracts)
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

func SubmitContractVerify(ctx context.Context, address string, r *SubmitContractVerifyRequest) (interface{},
	*utils.BuErrorResponse) {

	var contract busi.EVMContract
	exist, err := utils.EngineGroup[utils.DB].Where("address ilike ?", address).OrderBy("height desc").Get(&contract)
	if err != nil {
		log.Errorf("Execute sql error: %v", err)
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}
	if !exist {
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusNotFound,
			Response: utils.ErrBlockExplorerAPIServerNotFound}
	}

	count, err := utils.EngineGroup[utils.BusiDB].Table(new(busi.EVMContractVerify)).
		Where("address ilike ? and status=?", address, busi.EVMContractVerifyStatusSuccessfully).Count()
	if err != nil {
		log.Errorf("Execute sql error: %v", err)
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}
	if count > 0 {
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusBadRequest, Response: utils.ErrContractVerified}
	}

	input := &solc.Input{
		Language: "Solidity",
		Settings: solc.Settings{
			Optimizer: solc.Optimizer{
				Enabled: r.IsOptimization,
				Runs:    r.Runs,
			},
			EVMVersion: r.EVMVersion,
			OutputSelection: map[string]map[string][]string{
				"*": {
					"*": []string{"*"},
				},
			},
		},
	}
	if r.EVMVersion != "" && !strings.Contains(r.EVMVersion, "default") {
		input.Settings.EVMVersion = r.EVMVersion
	}

	// TODO need support more CompilerType
	if r.CompilerType == busi.CompilerTypeSingleFile {
		input.Sources = map[string]solc.SourceIn{
			"": {Content: r.SourceCode},
		}
	}

	ib, _ := json.Marshal(input)
	contractVerify := &busi.EVMContractVerify{
		Address:         address,
		CompilerType:    r.CompilerType,
		CompilerVersion: r.CompilerVersion,
		LicenseType:     r.LicenseType,
		Input:           string(ib),
		Status:          busi.EVMContractVerifyStatusDoing,
	}
	if _, err := utils.EngineGroup[utils.BusiDB].Insert(contractVerify); err != nil {
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}

	// async compile
	go asyncCompilerContract(input, contractVerify, &contract)

	return contractVerify, nil
}

var errByteCodeNotEqual = errors.New("bytecode not equal")

func asyncCompilerContract(input *solc.Input, cv *busi.EVMContractVerify, contract *busi.EVMContract) {
	var (
		output       *solc.Output
		err          error
		contractName string
	)

	defer func() {
		if err != nil {
			if err == errByteCodeNotEqual {
				cv.Status = busi.EVMContractVerifyStatusNotEqual
			} else {
				cv.Status = busi.EVMContractVerifyStatusUnknown
			}
			log.Errorf("compile error:%s", err)
		} else {
			cv.Status = busi.EVMContractVerifyStatusSuccessfully
		}

		if output != nil {
			o, _ := json.Marshal(output)
			cv.Output = string(o)
		}
		cv.ContractName = contractName
		if _, err = utils.EngineGroup[utils.BusiDB].ID(cv.ID).Update(cv); err != nil {
			log.Errorf("update contract verify failed, err:%s", err)
		}
	}()

	// like v0.3.6+commit.3fc68da5.js
	tmp := strings.Split(cv.CompilerVersion, "+")
	version := strings.Replace(tmp[0], "v", "", -1)
	compiler, err := solc.GetCompiler(version)
	if err != nil {
		log.Errorf("GetCompiler err：%s", err)
		return
	}
	output, err = compiler.Compile(input)
	if err != nil {
		log.Errorf("Compile err：%s", err)
		return
	}

	var (
		compliedByteCode string
		bytecodeHash     string
	)

	if cv.CompilerType == busi.CompilerTypeSingleFile {
		// get first key because there only one contract
		for cn, c := range output.Contracts[""] {
			compliedByteCode = c.EVM.DeployedBytecode.Object
			result := gjson.Get(c.Metadata, "settings.metadata.bytecodeHash")
			bytecodeHash = result.String()
			contractName = cn
			break
		}
	}

	equal, err := solc.Verify(compliedByteCode, contract.ByteCode, bytecodeHash)
	if err != nil {
		log.Errorf("solc.Verify failed, err:%s", err)
		return
	}

	if !equal {
		log.Errorf("bytecode not equal")
		err = errByteCodeNotEqual
		return
	}
}

func GetContractVerifyByID(ctx context.Context, id int) (interface{}, *utils.BuErrorResponse) {
	cv, buErr := getContractVerifyByQuery(ctx, "id=?", id)
	if buErr != nil {
		return nil, buErr
	}
	contractVerify := new(ContractVerify)
	contractVerify.ID = cv.ID
	contractVerify.ContractName = cv.ContractName
	contractVerify.Status = cv.Status
	contractVerify.CreateAt = cv.CreateAt
	contractVerify.UpdatedAt = cv.UpdatedAt
	contractVerify.LicenseType = cv.LicenseType
	contractVerify.CompilerVersion = cv.CompilerVersion
	contractVerify.CompilerType = cv.CompilerType
	contractVerify.Address = cv.Address

	var output solc.Output
	if err := json.Unmarshal([]byte(cv.Output), &output); err != nil {
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}
	if contractVerify.Status != busi.EVMContractVerifyStatusSuccessfully {
		if len(output.Errors) == 0 {
			contractVerify.ErrMsg = "unknown"
		} else {
			contractVerify.ErrMsg = output.Errors[0].Message
		}
		return contractVerify, nil
	}

	var input solc.Input
	if err := json.Unmarshal([]byte(cv.Input), &input); err != nil {
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}

	for fileName, code := range input.Sources {
		// use single file no filename
		if contractVerify.CompilerType == busi.CompilerTypeSingleFile {
			fileName = fmt.Sprintf("%s.sol", contractVerify.ContractName)
		}
		contractVerify.SourceCodes = append(contractVerify.SourceCodes, &SourceCode{
			FileName: fileName,
			Code:     code.Content,
		})
	}

	if contractVerify.CompilerType == busi.CompilerTypeSingleFile {
		for _, c := range output.Contracts[""] {
			contractVerify.ABI = gjson.Get(c.Metadata, "output.abi").String()
			break
		}
	}
	return contractVerify, nil
}

func GetSuccessContractVerifyByAddress(ctx context.Context, address string) (*busi.EVMContractVerify,
	*utils.BuErrorResponse) {
	return getContractVerifyByQuery(ctx, "address ilike ? and status=?", address,
		busi.EVMContractVerifyStatusSuccessfully)
}

func getContractVerifyByQuery(ctx context.Context, query interface{}, args ...interface{}) (*busi.EVMContractVerify,
	*utils.BuErrorResponse) {
	var contractVerify busi.EVMContractVerify
	exist, err := utils.EngineGroup[utils.BusiDB].Where(query, args...).Get(&contractVerify)
	if err != nil {
		log.Errorf("Execute sql error: %v", err)
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}
	if !exist {
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusNotFound,
			Response: utils.ErrBlockExplorerAPIServerNotFound}
	}
	return &contractVerify, nil
}

func ListCompileVersion(ctx context.Context) (interface{}, *utils.BuErrorResponse) {
	buildList, err := solc.GetBuildList()
	if err != nil {
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}
	var versions CompileVersionList
	for _, build := range buildList {
		versions.Versions = append(versions.Versions, &CompileVersion{
			Version:     build.Version,
			LongVersion: build.LongVersion,
		})
	}
	return versions, nil
}
