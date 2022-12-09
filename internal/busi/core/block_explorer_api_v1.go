package core

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"api-server/pkg/models/busi"
	"api-server/pkg/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/imxyb/solc-go"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	ethcommon "github.com/ethereum/go-ethereum/common"
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
		var c Contract

		c.Txns, err = evmTransactionCount(contract.Address)
		if err != nil {
			log.Error(err)
			break
		}

		c.Height = contract.Height
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
			c.CompilerType = contractVerify.CompilerType
			c.CompilerVersion = contractVerify.CompilerVersion
			c.License = contractVerify.LicenseType
			c.Verified = contractVerify.CreateAt
		}

		contractsSlice = append(contractsSlice, &c)
	}

	contractsList.Contracts = contractsSlice

	return contractsList, nil
}

func GetContract(ctx context.Context, address string) (interface{}, *utils.BuErrorResponse) {
	evmContract := new(busi.EVMContract)

	b, err := utils.EngineGroup[utils.TaskDB].
		Where("address=? or filecoin_address=?", address, address).
		OrderBy("height desc").Get(evmContract)
	if err != nil {
		log.Errorf("Execute sql error: %v", err)
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}

	var ethAddress string
	// if query address is filecoin address
	if address[0] == 'f' || address[0] == 't' {
		ethAddress = evmContract.Address
	} else {
		ethAddress = address
	}

	if !b {
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusOK, Response: utils.ErrBlockExplorerAPIServerNotFound}
	}

	// get the heightest and new version address
	// t := newContractsArr(evmContract)
	// sort.Sort(t)

	// contract := t[len(t)-1]

	contractDetail := ContractDetail{
		Address:         ethcommon.HexToAddress(evmContract.Address).Hex(),
		FilecoinAddress: evmContract.FilecoinAddress,
		Balance:         evmContract.Balance,
		Nonce:           evmContract.Nonce,
		ByteCode:        evmContract.ByteCode,
	}

	transaction, resErr := findCreatorTransaction(ethAddress)
	if err != nil {
		return nil, resErr
	}
	if transaction != nil {
		contractDetail.Creator = transaction.From
		contractDetail.Txn = transaction.Hash
	}

	var contractVerify busi.EVMContractVerify
	exist, err := utils.EngineGroup[utils.APIDB].Where("address=? and status=?", ethAddress,
		busi.EVMContractVerifyStatusSuccessfully).Get(&contractVerify)
	if err != nil {
		log.Errorf("Execute sql error: %v", err)
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}
	if exist {
		contractDetail.Verified = contractVerify.CreateAt
		contractDetail.ContractName = contractVerify.ContractName
		contractDetail.LicenseType = contractVerify.LicenseType
		contractDetail.CompilerVersion = contractVerify.CompilerVersion

		var output solc.Output
		if err := json.Unmarshal([]byte(contractVerify.Output), &output); err != nil {
			return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
				Response: utils.ErrBlockExplorerAPIServerInternal}
		}
		var input solc.Input
		if err := json.Unmarshal([]byte(contractVerify.Input), &input); err != nil {
			return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
				Response: utils.ErrBlockExplorerAPIServerInternal}
		}

		for fileName, code := range input.Sources {
			// use single file no filename
			if contractVerify.CompilerType == busi.CompilerTypeSingleFile {
				fileName = fmt.Sprintf("%s.sol", contractVerify.ContractName)
			}
			contractDetail.SourceCodes = append(contractDetail.SourceCodes, &SourceCode{
				FileName: fileName,
				Code:     code.Content,
			})
		}
		// main contract must at first
		sort.Slice(contractDetail.SourceCodes, func(i, j int) bool {
			if strings.Contains(contractDetail.SourceCodes[i].FileName, contractDetail.ContractName) {
				return true
			}
			return false
		})

		if contractVerify.CompilerType == busi.CompilerTypeSingleFile {
			for _, c := range output.Contracts[""] {
				contractDetail.ABI = gjson.Get(c.Metadata, "output.abi").String()
				break
			}
		} else if contractVerify.CompilerType == busi.CompilerTypeMultiPart {
			key := fmt.Sprintf("%s.sol", contractDetail.ContractName)
			metadata := output.Contracts[key][contractDetail.ContractName].Metadata
			contractDetail.ABI = gjson.Get(metadata, "output.abi").String()
		}
	}

	return contractDetail, nil
}

func ListContractTXNs(ctx context.Context, address string, r *ListQuery) (interface{}, *utils.BuErrorResponse) {
	var (
		txnsList TxnsList
	)

	// if address is contract, must add creator hash
	creatorTx, err := findCreatorTransaction(address)
	if err != nil {
		return nil, err
	}

	// get the numbers of transactions
	total, err := evmTransactionCount(address)
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
	// creator tx must at last
	if creatorTx != nil && len(transactions) < r.Limit {
		transactions = append(transactions, creatorTx)
		txnsList.Hits += 1
	}

	for _, transaction := range transactions {
		if transaction.To == "" {
			transaction.MethodName = "create"
		} else {
			transaction.MethodName, transaction.MethodSig, transaction.Params =
				parseMethodAndParamsFromContract(transaction.Input, address)
			if err != nil {
				return nil, err
			}
		}
	}
	txnsList.EVMTransaction = transactions

	return txnsList, nil
}

func ListContractEvents(ctx context.Context, address string, r *ListQuery) (interface{}, *utils.BuErrorResponse) {
	var (
		evmReceipts []busi.EVMReceipt
		events      []*Event
	)
	err := utils.EngineGroup[utils.TaskDB].Where("`to`=? and logs!='[]'", address).
		OrderBy("height desc").Limit(r.Limit).Find(&evmReceipts)
	if err != nil {
		log.Errorf("Execute sql error: %v", err)
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}
	tokenABI, err := getContractABI(address)
	if err != nil {
		log.Errorf("getContractABI error: %v", err)
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}
	if tokenABI == nil {
		return events, nil
	}

	for _, receipt := range evmReceipts {
		var ethLogs []types.Log
		if err = json.Unmarshal([]byte(receipt.Logs), &ethLogs); err != nil {
			log.Errorf("json.Unmarshal error: %v", err)
			return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
				Response: utils.ErrBlockExplorerAPIServerInternal}
		}
		for _, ethLog := range ethLogs {
			abiEvent, err := tokenABI.EventByID(ethLog.Topics[0])
			if err != nil {
				log.Errorf("tokenABI.EventByID error: %v", err)
				return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
					Response: utils.ErrBlockExplorerAPIServerInternal}
			}
			event := &Event{
				Address:     ethLog.Address.String(),
				RawData:     hex.EncodeToString(ethLog.Data),
				BlockNumber: ethLog.BlockNumber,
				TxHash:      ethLog.TxHash.String(),
				TxIndex:     ethLog.TxIndex,
				BlockHash:   ethLog.BlockHash.String(),
				Index:       ethLog.Index,
				EventName:   abiEvent.String(),
			}
			for _, topic := range ethLog.Topics {
				event.RawTopics = append(event.RawTopics, topic.String())
			}
			var indexedArgs []abi.Argument
			for _, input := range abiEvent.Inputs {
				if input.Indexed {
					indexedArgs = append(indexedArgs, input)
				}
			}
			event.ParsedTopics = make(map[string]interface{})
			if err = abi.ParseTopicsIntoMap(event.ParsedTopics, indexedArgs, ethLog.Topics[1:]); err != nil {
				log.Errorf("abi.ParseTopicsIntoMap error: %v", err)
				return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
					Response: utils.ErrBlockExplorerAPIServerInternal}
			}
			event.ParsedData = make(map[string]interface{})
			if err = abiEvent.Inputs.UnpackIntoMap(event.ParsedData, ethLog.Data); err != nil {
				log.Errorf("abiEvent.Inputs.UnpackIntoMap error: %v", err)
				return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
					Response: utils.ErrBlockExplorerAPIServerInternal}
			}
			events = append(events, event)
		}
	}

	// 因为一个receipt可能有多个event, 只取r.limit个，参考etherscan
	if len(events) > r.Limit {
		events = events[:r.Limit]
	}

	return EventList{Events: events}, nil
}

func ListInternalTXNs(ctx context.Context, address string, r *ListQuery) (interface{}, *utils.BuErrorResponse) {
	var (
		internalTXNsList InternalTxnsList
	)

	count, err := utils.EngineGroup[utils.TaskDB].Table(new(busi.EVMInternalTX)).
		Join("inner", "evm_transaction", "evm_internal_tx.parent_hash=evm_transaction.hash").
		Where("evm_transaction.from=? or evm_transaction.to=?", address, address).Count()
	if err != nil {
		log.Errorf("Execute sql error: %v", err)
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}

	internalTXNsList.Hits = count
	if internalTXNsList.Hits <= 0 {
		return internalTXNsList, nil
	}

	internalTxs := make([]*busi.EVMInternalTX, 0)
	err = utils.EngineGroup[utils.TaskDB].Select("evm_internal_tx.*").
		Join("inner", "evm_transaction", "evm_internal_tx.parent_hash=evm_transaction.hash").
		Where("evm_transaction.from=? or evm_transaction.to=?", address, address).
		Limit(r.Limit, r.Offset).OrderBy("height desc").Find(&internalTxs)
	if err != nil {
		log.Errorf("Execute sql error: %v", err)
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}

	internalTXNsList.EVMInternalTX = internalTxs

	return internalTXNsList, nil
}

func SubmitContractVerify(ctx context.Context, address string, r *SubmitContractVerifyRequest) (interface{},
	*utils.BuErrorResponse) {

	var contract busi.EVMContract
	exist, err := utils.EngineGroup[utils.TaskDB].Where("address=?", address).OrderBy("height desc").Get(&contract)
	if err != nil {
		log.Errorf("Execute sql error: %v", err)
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}
	if !exist {
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusNotFound,
			Response: utils.ErrBlockExplorerAPIServerNotFound}
	}

	count, err := utils.EngineGroup[utils.APIDB].Table(new(busi.EVMContractVerify)).
		Where("address=? and status=?", address, busi.EVMContractVerifyStatusSuccessfully).Count()
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
		Sources: map[string]solc.SourceIn{},
	}
	if r.EVMVersion != "" && !strings.Contains(r.EVMVersion, "default") {
		input.Settings.EVMVersion = r.EVMVersion
	}

	var mainContractFileName string
	input.Sources, mainContractFileName, err = buildSource(r)
	if err != nil {
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
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
	if _, err := utils.EngineGroup[utils.APIDB].Insert(contractVerify); err != nil {
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}

	// async compile
	go asyncCompilerContract(input, mainContractFileName, contractVerify, &contract)

	return contractVerify, nil
}

func buildSource(r *SubmitContractVerifyRequest) (map[string]solc.SourceIn, string, error) {
	var mainContractFileName string
	sources := make(map[string]solc.SourceIn)
	if r.CompilerType == busi.CompilerTypeSingleFile {
		sources[""] = solc.SourceIn{Content: r.SourceCode}
	} else if r.CompilerType == busi.CompilerTypeMultiPart {
		for _, part := range r.SourceCodeParts {
			baseName := filepath.Base(part.Filename)

			// main contract at first
			if mainContractFileName == "" {
				mainContractFileName = baseName
			}

			resp, err := http.Get(part.SourceCodeUrl)
			if err != nil {
				log.Errorf("get url %s faild, err:%s", part.SourceCodeUrl, err)
				return nil, "", &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
					Response: utils.ErrBlockExplorerAPIServerInternal}
			}
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				resp.Body.Close()
				log.Errorf("get url %s ReadAll failed, err:%s", part.SourceCodeUrl, err)
				return nil, "", err
			}
			resp.Body.Close()
			sources[baseName] = solc.SourceIn{Content: string(b)}
		}
	}
	return sources, mainContractFileName, nil
}

var errByteCodeNotEqual = errors.New("bytecode not equal")

func asyncCompilerContract(input *solc.Input, mainContractFileName string,
	cv *busi.EVMContractVerify, contract *busi.EVMContract) {
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
		if _, err = utils.EngineGroup[utils.APIDB].ID(cv.ID).Update(cv); err != nil {
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
	} else if cv.CompilerType == busi.CompilerTypeMultiPart {
		mainContract := strings.Replace(mainContractFileName, ".sol", "", -1)
		c := output.Contracts[mainContractFileName][mainContract]
		compliedByteCode = c.EVM.DeployedBytecode.Object
		result := gjson.Get(c.Metadata, "settings.metadata.bytecodeHash")
		bytecodeHash = result.String()
		contractName = mainContract
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
	json.Unmarshal([]byte(cv.Output), &output)
	if contractVerify.Status != busi.EVMContractVerifyStatusSuccessfully {
		if len(output.Errors) == 0 {
			contractVerify.ErrMsg = "bytecode not equal"
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
			contractVerify.Bytecode = c.EVM.DeployedBytecode.Object
			break
		}
	}
	return contractVerify, nil
}

func GetSuccessContractVerifyByAddress(ctx context.Context, address string) (*busi.EVMContractVerify,
	*utils.BuErrorResponse) {
	return getContractVerifyByQuery(ctx, "address=? and status=?", address, busi.EVMContractVerifyStatusSuccessfully)
}

func getContractVerifyByQuery(ctx context.Context, query interface{}, args ...interface{}) (*busi.EVMContractVerify,
	*utils.BuErrorResponse) {
	var contractVerify busi.EVMContractVerify
	exist, err := utils.EngineGroup[utils.APIDB].Where(query, args...).Get(&contractVerify)
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
		if strings.Contains(build.LongVersion, "nightly") {
			continue
		}
		versions.Versions = append(versions.Versions, &CompileVersion{
			Version:     build.Version,
			LongVersion: build.LongVersion,
		})
	}
	return versions, nil
}

func GetContractIsVerify(ctx context.Context, address string) (interface{}, *utils.BuErrorResponse) {
	count, err := utils.EngineGroup[utils.APIDB].Where("address=? and status=?",
		address, busi.EVMContractVerifyStatusSuccessfully).Table(new(busi.EVMContractVerify)).Count()
	if err != nil {
		log.Errorf("Execute sql error: %v", err)
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}
	var result ContractIsVerify
	if count > 0 {
		result.IsVerify = true
	}
	return result, nil
}

func ListTXNs(ctx context.Context, r *ListQuery) (interface{}, *utils.BuErrorResponse) {
	var (
		c        busi.EVMContract
		txnsList TxnsList
	)

	// get the numbers of txns
	total, err := busiTableRecordsCount(&c)
	if err != nil {
		return nil, err
	}

	txnsList.Hits = total
	if txnsList.Hits <= 0 {
		return txnsList, nil
	}

	evmTransaction := make([]*busi.EVMTransaction, 0)
	if err := busiSQLExecute(r, &evmTransaction); err != nil {
		return nil, err
	}
	for _, transaction := range evmTransaction {
		if transaction.To == "" {
			transaction.MethodName = "create"
		} else {
			transaction.MethodName, transaction.MethodSig, transaction.Params =
				parseMethodAndParamsFromContract(transaction.Input, transaction.To)
			if err != nil {
				return nil, err
			}
		}
	}
	txnsList.EVMTransaction = evmTransaction
	return txnsList, nil
}

func GetTXN(ctx context.Context, hash string) (interface{}, *utils.BuErrorResponse) {
	evmTransactions := make([]*busi.EVMTransaction, 0)

	// get txn
	err := utils.EngineGroup[utils.TaskDB].Where("hash = ?", hash).Find(&evmTransactions)
	if err != nil {
		log.Errorf("Execute sql error: %v", err)
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}

	if len(evmTransactions) == 0 {
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusOK, Response: utils.ErrBlockExplorerAPIServerNotFound}
	}

	var (
		maxHeight      int64
		evmTransaction busi.EVMTransaction
	)
	for _, t := range evmTransactions {
		if t.Height > maxHeight {
			maxHeight = t.Height
			evmTransaction = *t
		}
	}

	if evmTransaction.To == "" {
		evmTransaction.MethodName = "create"
	} else {
		evmTransaction.MethodName, evmTransaction.MethodSig, evmTransaction.Params =
			parseMethodAndParamsFromContract(evmTransaction.Input, evmTransaction.To)
		if err != nil {
			log.Errorf("parseMethodAndParamsFromContract error: %v", err)
			return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
				Response: utils.ErrBlockExplorerAPIServerInternal}
		}
	}

	// check if "to" is contract
	var resp EVMTransaction
	resp.EVMTransaction = evmTransaction

	evmContract := new(busi.EVMContract)
	if evmTransaction.To != "" {
		resp.ToIsContract, err = utils.EngineGroup[utils.TaskDB].Where("address = ?",
			evmTransaction.To).Get(evmContract)
		if err != nil {
			return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
				Response: utils.ErrBlockExplorerAPIServerInternal}
		}
	}

	// confirmation blocks count
	result, err := utils.EngineGroup[utils.TaskDB].QueryString("select max(height) from evm_block_header;")
	if err != nil {
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}
	max_height, err := strconv.ParseInt(result[0]["max"], 10, 64)
	if err != nil {
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}
	if max_height > 0 && (max_height > resp.Height) {
		resp.ConfirmationBlocks = int64(math.Max(50, float64(max_height-resp.Height)))
	}

	return resp, nil
}

func GetBlock(ctx context.Context, heightStr string) (interface{}, *utils.BuErrorResponse) {

	height, err := strconv.ParseInt(heightStr, 10, 64)
	if err != nil {
		log.Errorf("GetBlock ParseInt err: %v", err)
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusOK, Response: utils.ErrBlockExplorerAPIServerParams}
	}

	evmBlockHeader := new(busi.EVMBlockHeader)

	b, err := utils.EngineGroup[utils.TaskDB].Where("height = ?", height).Get(evmBlockHeader)
	if err != nil {
		log.Errorf("Execute sql error: %v", err)
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}

	if !b {
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusOK, Response: utils.ErrBlockExplorerAPIServerNotFound}
	}

	return evmBlockHeader, nil
}

func GetAddress(ctx context.Context, address string) (interface{}, *utils.BuErrorResponse) {
	evmAddress := new(busi.EVMAddress)

	b, err := utils.EngineGroup[utils.TaskDB].
		Where("address=? or filecoin_address=?", address, address).
		OrderBy("height desc").Get(evmAddress)
	if err != nil {
		log.Errorf("Execute sql error: %v", err)
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusInternalServerError,
			Response: utils.ErrBlockExplorerAPIServerInternal}
	}
	if !b {
		return nil, &utils.BuErrorResponse{HttpCode: http.StatusOK, Response: utils.ErrBlockExplorerAPIServerNotFound}
	}
	return evmAddress, nil
}

func ListAddressTXNs(ctx context.Context, address string, r *ListQuery) (interface{}, *utils.BuErrorResponse) {
	var (
		txnsList TxnsList
	)

	// get the numbers of transactions
	total, err := evmTransactionCount(address)
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
	for _, transaction := range transactions {
		if transaction.To == "" {
			transaction.MethodName = "create"
		} else {
			transaction.MethodName, transaction.MethodSig, transaction.Params =
				parseMethodAndParamsFromContract(transaction.Input, transaction.To)
			if err != nil {
				return nil, err
			}
		}
	}
	txnsList.EVMTransaction = transactions

	return txnsList, nil
}
